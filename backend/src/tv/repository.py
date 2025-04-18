from sqlalchemy import select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session, joinedload

from tv.models import Season, Show, Episode, SeasonRequest, SeasonFile
from tv.schemas import Season as SeasonSchema, SeasonId, Show as ShowSchema, ShowId, \
    SeasonRequest as SeasonRequestSchema, SeasonFile as SeasonFileSchema


def get_show(show_id: ShowId, db: Session) -> ShowSchema | None:
    """
    Retrieve a show by its ID, including seasons and episodes.

    :param show_id: The ID of the show to retrieve.
    :param db: The database session.
    :return: A ShowSchema object if found, otherwise None.
    """
    stmt = (
        select(Show)
        .where(Show.id == show_id)
        .options(
            joinedload(Show.seasons).joinedload(Season.episodes)  # Load relationships
        )
    )

    result = db.execute(stmt).unique().scalar_one_or_none()
    if not result:
        return None

    return ShowSchema.model_validate(result)


def get_show_by_external_id(external_id: int, db: Session, metadata_provider: str) -> ShowSchema | None:
    """
    Retrieve a show by its external ID, including nested seasons and episodes.

    :param external_id: The ID of the show to retrieve.
    :param metadata_provider: The metadata provider associated with the ID.
    :param db: The database session.
    :return: A ShowSchema object if found, otherwise None.
    """
    stmt = (
        select(Show)
        .where(Show.external_id == external_id)
        .where(Show.metadata_provider == metadata_provider)
        .options(
            joinedload(Show.seasons).joinedload(Season.episodes)  # Load relationships
        )
    )

    result = db.execute(stmt).unique().scalar_one_or_none()
    if not result:
        return None

    return ShowSchema(
        **result.__dict__
    )


def get_shows(db: Session) -> list[ShowSchema]:
    """
    Retrieve all shows from the database, including nested seasons and episodes.

    :param db: The database session.
    :return: A list of ShowSchema objects.
    """
    stmt = select(Show)

    results = db.execute(stmt).scalars().all()

    return [ShowSchema.model_validate(show) for show in results]


def save_show(show: ShowSchema, db: Session) -> ShowSchema:
    """
    Save a new show to the database, including its seasons and episodes.

    :param show: The ShowSchema object to save.
    :param db: The database session.
    :return: The saved ShowSchema object.
    :raises ValueError: If a show with the same primary key already exists.
    """
    db_show = Show(
        id=show.id,
        external_id=show.external_id,
        metadata_provider=show.metadata_provider,
        name=show.name,
        overview=show.overview,
        year=show.year,
        seasons=[
            Season(
                id=season.id,
                show_id=show.id,  # Correctly linking to the parent show
                number=season.number,
                external_id=season.external_id,
                name=season.name,
                overview=season.overview,
                episodes=[
                    Episode(
                        id=episode.id,
                        season_id=season.id,  # Correctly linking to the parent season
                        number=episode.number,
                        external_id=episode.external_id,
                        title=episode.title
                    ) for episode in season.episodes  # Convert episodes properly
                ]
            ) for season in show.seasons  # Convert seasons properly
        ]
    )

    db.add(db_show)
    try:
        db.commit()
        db.refresh(db_show)
        return ShowSchema.model_validate(db_show)
    except IntegrityError:
        db.rollback()
        raise ValueError("Show with this primary key already exists.")


def delete_show(show_id: ShowId, db: Session) -> None:
    """
    Delete a show by its ID.

    :param show_id: The ID of the show to delete.
    :param db: The database session.
    :return: The deleted ShowSchema object if found, otherwise None.
    """
    show = db.get(Show, show_id)
    db.delete(show)
    db.commit()


def get_season(season_id: SeasonId, db: Session) -> SeasonSchema:
    """

    :param season_id: The ID of the season to get.
    :param db: The database session.
    :return: a Season object.
    """
    return SeasonSchema.model_validate(db.get(Season(), season_id))


def add_season_to_requested_list(season_request: SeasonRequestSchema, db: Session) -> None:
    """
    Adds a Season to the SeasonRequest table, which marks it as requested.

    """
    db.add(SeasonRequest(**season_request.model_dump()))
    db.commit()


def remove_season_from_requested_list(season_request: SeasonRequestSchema, db: Session) -> None:
    """
    Removes a Season from the SeasonRequest table, which removes it from the 'requested' list.

    """
    db.delete(SeasonRequest(**season_request.model_dump()))
    db.commit()


def get_season_by_number(db: Session, season_number: int, show_id: ShowId) -> SeasonSchema:
    stmt = (
        select(Season).
        where(Season.show_id == show_id).
        where(Season.number == season_number).
        options(
            joinedload(Season.episodes).joinedload(Season.show)
        )
    )
    result = db.execute(stmt).unique().scalar_one_or_none()
    return SeasonSchema.model_validate(result)

def get_season_requests(db: Session) -> list[SeasonRequestSchema]:
    stmt = select(SeasonRequest)
    result = db.execute(stmt).scalars().all()
    return [SeasonRequestSchema.model_validate(season) for season in result]


def add_season_file(db: Session, season_file: SeasonFileSchema) -> SeasonFileSchema:
    db.add(SeasonFile(**season_file.model_dump()))
    db.commit()
    return season_file
