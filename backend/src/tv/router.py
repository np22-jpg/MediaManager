import json
import pprint
from uuid import UUID

import psycopg.errors
from fastapi import APIRouter, Depends, status
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from sqlmodel import select

import dowloadClients
import indexer
import metadataProvider
from auth.users import current_active_user
from database import DbSessionDependency
from database.torrents import Torrent
from database.tv import Season, Show
from indexer import IndexerQueryResult
from tv import log

router = APIRouter(
    prefix="/tv",
)


class ShowDetails(BaseModel):
    show: Show
    seasons: list[Season]


@router.post("/show", status_code=status.HTTP_201_CREATED, dependencies=[Depends(current_active_user)],
             responses={
                 status.HTTP_201_CREATED: {"model": Show, "description": "Successfully created show"},
                 status.HTTP_409_CONFLICT: {"model": str, "description": "Show already exists"},
             })
def add_show(db: DbSessionDependency, show_id: int, metadata_provider: str = "tmdb", version: str = ""):
    res = db.exec(select(Show).
                  where(Show.external_id == show_id).
                  where(Show.metadata_provider == metadata_provider).
                  where(Show.version == version)).first()

    if res is not None:
        return JSONResponse(status_code=status.HTTP_409_CONFLICT, content={"message": "Show already exists"})

    show = metadataProvider.get_show_metadata(id=show_id, provider=metadata_provider)
    show.version = version
    log.info("Adding show: " + json.dumps(show.model_dump(), default=str))
    db.add(show)
    try:
        db.commit()
        db.refresh(show)
    except psycopg.errors.UniqueViolation as e:
        log.debug(e)
        log.info("Show already exists " + show.__str__())
        return JSONResponse(status_code=status.HTTP_409_CONFLICT, content={"message": "Show already exists"})

    return show


@router.delete("/{show_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def delete_show(db: DbSessionDependency, show_id: UUID):
    db.delete(db.get(Show, show_id))
    db.commit()


@router.patch("/{show_id}/{season_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
              response_model=Season)
def add_season(db: DbSessionDependency, season_id: UUID):
    """
    adds requested flag to a season
    """
    season = db.get(Season, season_id)
    season.requested = True
    db.add(season)
    db.commit()
    db.refresh(season)

    return season


@router.delete("/{show_id}/{season_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
               response_model=Show)
def delete_season(db: DbSessionDependency, show_id: UUID, season: int):
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
    current_active_user)],
            response_model=list[IndexerQueryResult])
def get_season_torrents(db: DbSessionDependency, show_id: UUID, season_id: UUID):
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
    current_active_user)], response_model=list[Season])
def download_seasons_torrent(db: DbSessionDependency, show_id: UUID, torrent_id: UUID):
    """
    downloads torrents for a show season, links the torrent for all seasons the torrent contains

    """
    torrent = db.get(Torrent, torrent_id)

    if torrent is None:
        return JSONResponse(status_code=status.HTTP_404_NOT_FOUND, content={"message": "Torrent not found"})

    seasons = []
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


@router.post("/{show_id}/{season_id}/torrent", status_code=status.HTTP_200_OK, dependencies=[Depends(
    current_active_user)], response_model=list[Season])
def delete_seasons_torrent(db: DbSessionDependency, show_id: UUID, season_id: UUID, torrent_id: UUID):
    """
    downloads torrents for a season, links the torrent only to the specified season
    this means that multiple torrents can contain a season but you can choose from one which the content should be
    imported

    """
    torrent = db.get(Torrent, torrent_id)

    if torrent is None:
        return JSONResponse(status_code=status.HTTP_404_NOT_FOUND, content={"message": "Torrent not found"})

    seasons = []
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


@router.get("/", dependencies=[Depends(current_active_user)], response_model=list[Show])
def get_shows(db: DbSessionDependency):
    """"""
    return db.exec(select(Show)).unique().fetchall()


@router.get("/{show_id}", dependencies=[Depends(current_active_user)], response_model=ShowDetails)
def get_show(db: DbSessionDependency, show_id: UUID):
    """

    :param show_id:
    :type show_id:
    :return:
    :rtype:
    """
    shows = db.execute(select(Show, Season).where(Show.id == show_id).join(Season).order_by(Season.number)).fetchall()
    seasons = []
    for show in shows:
        seasons.append(show[1])

    shows = db.execute(select(Show, Season).where(Show.id == show_id).join(Season).order_by(Season.number))

    return ShowDetails(show=shows.first()[0], seasons=seasons)


@router.get("/search", dependencies=[Depends(current_active_user)])
def search_show(query: str, metadata_provider: str = "tmdb"):
    return metadataProvider.search_show(query, metadata_provider)
