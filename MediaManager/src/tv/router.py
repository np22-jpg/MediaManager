import json
from typing import List
from uuid import UUID

import psycopg.errors
from fastapi import APIRouter, Depends, status
from fastapi.responses import JSONResponse
from sqlmodel import select
from tmdbsimple import TV, TV_Seasons

import auth
from database import SessionDependency
from routers.users import Message
from tv import Show, tmdb, log, Season, Episode

router = APIRouter(
    prefix="/tv",
)


@router.post("/show", status_code=status.HTTP_201_CREATED, dependencies=[Depends(auth.get_current_user)],
             responses={
                 status.HTTP_201_CREATED: {"model": Show, "description": "Successfully created show"},
                 status.HTTP_409_CONFLICT: {"model": Message, "description": "Show already exists"},
             })
def add_show(db: SessionDependency, show_id: int, metadata_provider: str = "tmdb"):
    show_metadata = TV(show_id).info()
    show = Show(
        external_id=show_id,
        metadata_provider=metadata_provider,
        name=show_metadata["name"],
        overview=show_metadata["overview"]
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


@router.patch("/{show_id}/{season}", status_code=status.HTTP_200_OK, dependencies=[Depends(auth.get_current_user)],
              response_model=Show)
def add_season(db: SessionDependency, show_id: UUID, season: int):
    """
    adds requested flag to a season
    """
    season = db.get(Season, (show_id, season))
    season.requested = True
    db.add(season)
    db.commit()
    db.refresh(season)
    return season


@router.delete("/{show_id}/{season}", status_code=status.HTTP_200_OK, dependencies=[Depends(auth.get_current_user)],
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


@router.get("/show", dependencies=[Depends(auth.get_current_user)], response_model=List[Show])
def get_shows(db: SessionDependency):
    return db.exec(select(Show)).unique().fetchall()


@router.get("/search")
def search_show(query: str):
    search = tmdb.Search()
    return search.tv(query=query)
