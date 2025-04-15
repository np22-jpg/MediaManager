from sqlalchemy.orm import Session

import database
import indexer
# import indexer
import metadataProvider
import tv.repository
from indexer import IndexerQueryResult
from tv.exceptions import MediaAlreadyExists
from tv.repository import get_show_by_external_id
from tv.schemas import Show, ShowId, SeasonRequest


def add_show(db: Session, external_id: int, metadata_provider: str) -> Show | None:
    if check_if_show_exists(db=db, external_id=external_id, metadata_provider=metadata_provider):
        raise MediaAlreadyExists(f"Show with external ID {external_id} and" +
                                 f" metadata provider {metadata_provider} already exists")
    show_with_metadata = metadataProvider.get_show_metadata(id=external_id, provider=metadata_provider)
    saved_show = tv.repository.save_show(db=db, show=show_with_metadata)
    return saved_show


def request_season(db: Session, season_request: SeasonRequest) -> None:
    tv.repository.add_season_to_requested_list(db=db, season_request=season_request)


def unrequest_season(db: Session, season_request: SeasonRequest) -> None:
    tv.repository.remove_season_from_requested_list(db=db, season_request=season_request)


def check_if_show_exists(db: Session,
                         external_id: int = None,
                         metadata_provider: str = None,
                         show_id: ShowId = None) -> bool:
    if external_id and metadata_provider:
        if tv.repository.get_show_by_external_id(external_id=external_id, metadata_provider=metadata_provider, db=db):
            return True
        else:
            return False
    elif show_id:
        if tv.repository.get_show(show_id=show_id, db=db):
            return True
        else:
            return False
    else:
        raise ValueError("External ID and metadata provider or Show ID must be provided")


def get_all_available_torrents_for_a_season(db: Session, season_number: int, show_id: ShowId) -> list[
    IndexerQueryResult]:
    show = tv.repository.get_show(show_id=show_id, db=db)
    torrents: list[IndexerQueryResult] = indexer.search(show.name + " S" + str(season_number))
    result = []
    for torrent in torrents:
        if season.number in torrent.season:
            result.append(torrent)
    return result


def get_all_shows(db: Session) -> list[Show]:
    return tv.repository.get_shows(db=db)


def get_show_by_id(db: Session, show_id: ShowId) -> Show | None:
    return tv.repository.get_show(show_id=show_id, db=db)


def get_all_requested_seasons(db: Session) -> list[SeasonRequest]:
    return tv.repository.get_season_requests(db=db)


if __name__ == "__main__":
    session = database.SessionLocal()

    try:
        show = add_show(db=session, external_id=1418, metadata_provider="tmdb")
    except MediaAlreadyExists as e:
        print(e)
        show = get_show_by_external_id(db=session, external_id=1418, metadata_provider="tmdb")

    print(show)
    print(show.name)
    for season in show.seasons:
        print(season)
        # print(season.number)
        for episode in season.episodes:
            print(episode)
            # print(episode.title)
