import json
import pprint
from typing import List
from uuid import UUID

import psycopg.errors
from fastapi import APIRouter, Depends, status
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from sqlmodel import select
from tmdbsimple import TV, TV_Seasons

import auth
import dowloadClients
import indexer
from database import SessionDependency
from database.torrents import Torrent
from database.tv import Episode, Season, Show
from indexer import IndexerQueryResult
from routers.users import Message
from tv import log, tmdb

router = APIRouter(
    prefix="/tv",
)


class ShowDetails(BaseModel):
    show: Show
    seasons: list[Season]


@router.post("/show", status_code=status.HTTP_201_CREATED, dependencies=[Depends(auth.get_current_user)],
             responses={
                 status.HTTP_201_CREATED: {"model": Show, "description": "Successfully created show"},
                 status.HTTP_409_CONFLICT: {"model": Message, "description": "Show already exists"},
             })
def add_show(db: SessionDependency, show_id: int, metadata_provider: str = "tmdb"):
    show_metadata = TV(show_id).info()

    # For some shows the first_air_date isn't known, therefore it needs to be nullable
    year: str | None = show_metadata["first_air_date"]
    if year:
        year: int = int(year.split('-')[0])
    else:
        year = None

    show = Show(
        external_id=show_id,
        metadata_provider=metadata_provider,
        name=show_metadata["name"],
        overview=show_metadata["overview"],
        year=year,
    )

    log.info("Adding show: " + json.dumps(show.model_dump(), default=str))
    db.add(show)
    db.commit()

    for season in show_metadata["seasons"]:
        season_metadata = TV_Seasons(tv_id=show_metadata["id"], season_number=season["season_number"]).info()
        db.add(Season(
            show_id=show.id,
            number=int(season_metadata["season_number"]),
            name=season_metadata["name"],
            overview=season_metadata["overview"],
            external_id=int(season_metadata["id"]))
        )
        db.commit()

        for episode in season_metadata["episodes"]:
            db.add(Episode(
                show_id=show.id,
                season_number=int(season_metadata["season_number"]),
                title=episode["name"],
                number=int(episode["episode_number"]),
                external_id=int(episode["id"]),
            ))

    try:
        db.commit()
        db.refresh(show)
    except psycopg.errors.UniqueViolation as e:
        log.debug(e)
        log.info("Show already exists " + show.__str__())
        return JSONResponse(status_code=status.HTTP_409_CONFLICT, content={"message": "Show already exists"})

    return show


@router.delete("/{show_id}", status_code=status.HTTP_200_OK)
def delete_show(db: SessionDependency, show_id: UUID):
    db.delete(db.get(Show, show_id))
    db.commit()


@router.patch("/{show_id}/{season_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(auth.get_current_user)],
              response_model=Season)
def add_season(db: SessionDependency, season_id: UUID):
    """
    adds requested flag to a season
    """
    season = db.get(Season, season_id)
    season.requested = True
    db.add(season)
    db.commit()
    db.refresh(season)

    return season


@router.delete("/{show_id}/{season_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(auth.get_current_user)],
               response_model=Show)
def delete_season(db: SessionDependency, show_id: UUID, season: int):
    """
    removes requested flag from a season
    """
    season = db.get(Season, (show_id, season))
    season.requested = False
    db.add(season)
    db.commit()
    db.refresh(season)
    return season


@router.get("/{show_id}/{season_id}/torrent", status_code=status.HTTP_200_OK, dependencies=[Depends(
    auth.get_current_user)],
            response_model=list[IndexerQueryResult])
def get_season_torrents(db: SessionDependency, show_id: UUID, season_id: UUID):
    season = db.get(Season, season_id)

    if season is None:
        return JSONResponse(status_code=status.HTTP_404_NOT_FOUND, content={"message": "Season not found"})

    torrents: list[IndexerQueryResult] = indexer.search(season)
    result = []
    for torrent in torrents:
        if season.number in torrent.season:
            result.append(torrent)

    db.commit()
    if len(result) == 0:
        return result
    result.sort()

    log.info(f"Found {torrents.__len__()} torrents for show {season.show.name} season {season.number}, of which "
             f"{result.__len__()} torrents fit the query")
    log.debug(f"unfiltered torrents: \n{pprint.pformat(torrents)}\nfiltered torrents: \n{pprint.pformat(result)}")
    return result


@router.post("/{show_id}/torrent", status_code=status.HTTP_200_OK, dependencies=[Depends(
    auth.get_current_user)], response_model=list[Season])
def download_seasons_torrent(db: SessionDependency, show_id: UUID, torrent: IndexerQueryResult, ):
    seasons: list[Season] = []
    for season_number in torrent.season:
        seasons.append(
            db.exec(select(Season)
                    .where(Season.show_id == show_id)
                    .where(Season.number == season_number)
                    ).first()
        )

    torrent = torrent.download()

    dowloadClients.client.download(Torrent)

    for season in seasons:
        season.requested = True
        season.torrent_id = torrent.id

    return seasons


@router.get("/", dependencies=[Depends(auth.get_current_user)], response_model=List[Show])
def get_shows(db: SessionDependency):
    return db.exec(select(Show)).unique().fetchall()


@router.get("/{show_id}", dependencies=[Depends(auth.get_current_user)], response_model=ShowDetails)
def get_show(db: SessionDependency, show_id: UUID):
    shows = db.execute(select(Show, Season).where(Show.id == show_id).join(Season).order_by(Season.number)).fetchall()
    seasons = []
    for show in shows:
        seasons.append(show[1])

    shows = db.execute(select(Show, Season).where(Show.id == show_id).join(Season).order_by(Season.number))

    return ShowDetails(show=shows.first()[0], seasons=seasons)


@router.get("/search")
def search_show(query: str):
    search = tmdb.Search()
    return search.tv(query=query)
