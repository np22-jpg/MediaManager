from sqlalchemy.orm import Session

import media_manager.indexer.service
import media_manager.metadataProvider
import media_manager.torrent.repository
import media_manager.tv.repository
from media_manager.database import SessionLocal
from media_manager.indexer import IndexerQueryResult
from media_manager.indexer.schemas import IndexerQueryResultId
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.torrent.repository import get_seasons_files_of_torrent
from media_manager.torrent.schemas import Torrent
from media_manager.torrent.service import TorrentService
from media_manager.tv import log
from media_manager.tv.exceptions import MediaAlreadyExists
from media_manager.tv.repository import add_season_file, get_season_files_by_season_id
from media_manager.tv.schemas import (
    Show,
    ShowId,
    SeasonRequest,
    SeasonFile,
    SeasonId,
    Season,
    RichShowTorrent,
    RichSeasonTorrent,
    PublicSeason,
    PublicShow,
    PublicSeasonFile,
    SeasonNumber,
    SeasonRequestId,
    RichSeasonRequest,
)
from media_manager.torrent.schemas import QualityStrings


def add_show(db: Session, external_id: int, metadata_provider: str) -> Show | None:
    if check_if_show_exists(
            db=db, external_id=external_id, metadata_provider=metadata_provider
    ):
        raise MediaAlreadyExists(
            f"Show with external ID {external_id} and"
            + f" metadata provider {metadata_provider} already exists"
        )
    show_with_metadata = media_manager.metadataProvider.get_show_metadata(
        id=external_id, provider=metadata_provider
    )
    saved_show = media_manager.tv.repository.save_show(db=db, show=show_with_metadata)
    return saved_show


def add_season_request(db: Session, season_request: SeasonRequest) -> None:
    media_manager.tv.repository.add_season_request(db=db, season_request=season_request)


def get_season_request_by_id(
        db: Session, season_request_id: SeasonRequestId
) -> SeasonRequest | None:
    return media_manager.tv.repository.get_season_request(
        db=db, season_request_id=season_request_id
    )


def update_season_request(db: Session, season_request: SeasonRequest) -> None:
    media_manager.tv.repository.delete_season_request(
        db=db, season_request_id=season_request.id
    )
    media_manager.tv.repository.add_season_request(db=db, season_request=season_request)


def delete_season_request(db: Session, season_request_id: SeasonRequestId) -> None:
    media_manager.tv.repository.delete_season_request(
        db=db, season_request_id=season_request_id
    )


def get_public_season_files_by_season_id(
        db: Session, season_id: SeasonId
) -> list[PublicSeasonFile]:
    season_files = get_season_files_by_season_id(db=db, season_id=season_id)
    public_season_files = [PublicSeasonFile.model_validate(x) for x in season_files]
    result = []
    for season_file in public_season_files:
        if season_file_exists_on_file(db=db, season_file=season_file):
            season_file.downloaded = True
        result.append(season_file)
    return result


def get_public_season_files_by_season_number(
        db: Session, season_number: SeasonNumber, show_id: ShowId
) -> list[PublicSeasonFile]:
    season = media_manager.tv.repository.get_season_by_number(
        db=db, season_number=season_number, show_id=show_id
    )
    return get_public_season_files_by_season_id(db=db, season_id=season.id)


def check_if_show_exists(
        db: Session,
        external_id: int = None,
        metadata_provider: str = None,
        show_id: ShowId = None,
) -> bool:
    if external_id and metadata_provider:
        if media_manager.tv.repository.get_show_by_external_id(
                external_id=external_id, metadata_provider=metadata_provider, db=db
        ):
            return True
        else:
            return False
    elif show_id:
        if media_manager.tv.repository.get_show(show_id=show_id, db=db):
            return True
        else:
            return False
    else:
        raise ValueError(
            "External ID and metadata provider or Show ID must be provided"
        )


def get_all_available_torrents_for_a_season(
        db: Session, season_number: int, show_id: ShowId, search_query_override: str = None
) -> list[IndexerQueryResult]:
    log.debug(
        f"getting all available torrents for season {season_number} for show {show_id}"
    )
    show = media_manager.tv.repository.get_show(show_id=show_id, db=db)
    if search_query_override:
        search_query = search_query_override
    else:
        # TODO: add more Search query strings and combine all the results, like "season 3", "s03", "s3"
        search_query = show.name + " s" + str(season_number).zfill(2)
    torrents: list[IndexerQueryResult] = media_manager.indexer.service.search(
        query=search_query, db=db
    )
    if search_query_override:
        log.debug(
            f"Found with search query override {torrents.__len__()} torrents: {torrents}"
        )
        return torrents
    result: list[IndexerQueryResult] = []
    for torrent in torrents:
        if season_number in torrent.season:
            result.append(torrent)
    result.sort()
    return result


def get_all_shows(db: Session) -> list[Show]:
    return media_manager.tv.repository.get_shows(db=db)


def search_for_show(
        query: str, metadata_provider: str, db: Session
) -> list[MetaDataProviderShowSearchResult]:
    results = media_manager.metadataProvider.search_show(query, metadata_provider)
    for result in results:
        if check_if_show_exists(
                db=db, external_id=result.external_id, metadata_provider=metadata_provider
        ):
            result.added = True
    return results


def get_popular_shows(metadata_provider: str, db: Session):
    results: list[MetaDataProviderShowSearchResult] = (
        media_manager.metadataProvider.search_show(provider=metadata_provider)
    )

    for result in results:
        if check_if_show_exists(
                db=db, external_id=result.external_id, metadata_provider=metadata_provider
        ):
            results.pop(results.index(result))
    return results


def get_public_show_by_id(db: Session, show_id: ShowId) -> PublicShow:
    show = media_manager.tv.repository.get_show(show_id=show_id, db=db)
    seasons = [PublicSeason.model_validate(season) for season in show.seasons]
    for season in seasons:
        season.downloaded = is_season_downloaded(db=db, season_id=season.id)
    public_show = PublicShow.model_validate(show)
    public_show.seasons = seasons
    return public_show


def get_show_by_id(db: Session, show_id: ShowId) -> Show:
    return media_manager.tv.repository.get_show(show_id=show_id, db=db)


def is_season_downloaded(db: Session, season_id: SeasonId) -> bool:
    season_files = get_season_files_by_season_id(db=db, season_id=season_id)
    for season_file in season_files:
        if season_file_exists_on_file(db=db, season_file=season_file):
            return True

    return False


def season_file_exists_on_file(db: Session, season_file: SeasonFile) -> bool:
    if season_file.torrent_id is None:
        return True
    else:
        torrent_file = media_manager.torrent.repository.get_torrent_by_id(
            db=db, torrent_id=season_file.torrent_id
        )
        if torrent_file.imported:
            return True

    return False


def get_show_by_external_id(
        db: Session, external_id: int, metadata_provider: str
) -> Show | None:
    return media_manager.tv.repository.get_show_by_external_id(
        external_id=external_id, metadata_provider=metadata_provider, db=db
    )


def get_season(db: Session, season_id: SeasonId) -> Season:
    return media_manager.tv.repository.get_season(season_id=season_id, db=db)


def get_all_season_requests(db: Session) -> list[RichSeasonRequest]:
    return media_manager.tv.repository.get_season_requests(db=db)


def get_torrents_for_show(db: Session, show: Show) -> RichShowTorrent:
    show_torrents = media_manager.tv.repository.get_torrents_by_show_id(
        db=db, show_id=show.id
    )
    rich_season_torrents = []
    for show_torrent in show_torrents:
        seasons = media_manager.tv.repository.get_seasons_by_torrent_id(
            db=db, torrent_id=show_torrent.id
        )
        season_files = get_seasons_files_of_torrent(db=db, torrent_id=show_torrent.id)
        file_path_suffix = season_files[0].file_path_suffix
        season_torrent = RichSeasonTorrent(
            torrent_id=show_torrent.id,
            torrent_title=show_torrent.title,
            status=show_torrent.status,
            quality=show_torrent.quality,
            imported=show_torrent.imported,
            seasons=seasons,
            file_path_suffix=file_path_suffix,
        )
        rich_season_torrents.append(season_torrent)
    return RichShowTorrent(
        show_id=show.id,
        name=show.name,
        year=show.year,
        metadata_provider=show.metadata_provider,
        torrents=rich_season_torrents,
    )


def get_all_shows_with_torrents(db: Session) -> list[RichShowTorrent]:
    shows = media_manager.tv.repository.get_all_shows_with_torrents(db=db)
    return [get_torrents_for_show(show=show, db=db) for show in shows]


def download_torrent(
        db: Session,
        public_indexer_result_id: IndexerQueryResultId,
        show_id: ShowId,
        override_show_file_path_suffix: str = "",
) -> Torrent:
    indexer_result = media_manager.indexer.service.get_indexer_query_result(
        db=db, result_id=public_indexer_result_id
    )
    show_torrent = TorrentService(db=db).download(indexer_result=indexer_result)

    for season_number in indexer_result.season:
        season = media_manager.tv.repository.get_season_by_number(
            db=db, season_number=season_number, show_id=show_id
        )
        season_file = SeasonFile(
            season_id=season.id,
            quality=indexer_result.quality,
            torrent_id=show_torrent.id,
            file_path_suffix=override_show_file_path_suffix,
        )
        add_season_file(db=db, season_file=season_file)
    return show_torrent


def download_approved_season_request(
        db: Session,
        season_request: SeasonRequest,
        show_id: ShowId,
) -> bool:
    if not season_request.authorized:
        log.error(f"Season request {season_request.id} is not authorized for download")
        raise ValueError(
            f"Season request {season_request.id} is not authorized for download"
        )
    log.info(f"Downloading approved season request {season_request.id}")

    season = get_season(db=db, season_id=season_request.season_id)
    torrents = get_all_available_torrents_for_a_season(
        db=db, season_number=season.number, show_id=show_id
    )
    available_torrents: list[IndexerQueryResult] = []
    for torrent in torrents:
        if (
                torrent.quality > season_request.wanted_quality
                or torrent.quality < season_request.min_quality
                or torrent.seeders < 3
        ):
            log.info(
                f"Skipping torrent {torrent.title} with quality {torrent.quality} for season {season.id}, because it does not match the requested quality {season_request.wanted_quality}"
            )
        elif torrent.season == [
            season.number,
        ]:
            log.info(
                f"Skipping torrent {torrent.title} with quality {torrent.quality} for season {season.id}, because it contains to many/wrong seasons {torrent.season} (wanted: {season.number})"
            )
        else:
            available_torrents.append(torrent)
            log.info(
                f"Taking torrent  {torrent.title} with quality {torrent.quality} for season {season.id} into consideration"
            )
    if len(available_torrents) == 0:
        log.warning(
            f"No torrents matching criteria were found (wanted quality: {season_request.wanted_quality}, min_quality: {season_request.min_quality} for season {season.id})"
        )
        return False

    available_torrents.sort()

    torrent = TorrentService(db=db).download(indexer_result=available_torrents[0])
    season_file = SeasonFile(
        season_id=season.id,
        quality=torrent.quality,
        torrent_id=torrent.id,
        file_path_suffix=QualityStrings[torrent.quality.name].value.upper(),
    )
    add_season_file(db=db, season_file=season_file)
    return True


def auto_download_all_approved_season_requests() -> None:
    db: Session = SessionLocal()
    log.info("Auto downloading all approved season requests")
    season_requests = media_manager.tv.repository.get_season_requests(db=db)
    log.info(f"Found {len(season_requests)} season requests to process")
    count = 0
    for season_request in season_requests:
        if season_request.authorized:
            log.info(f"Processing season request {season_request.id} for download")
            show = media_manager.tv.repository.get_show_by_season_id(
                db=db, season_id=season_request.season_id
            )
            if download_approved_season_request(
                    db=db, season_request=season_request, show_id=show.id
            ):
                count += 1
            else:
                log.warning(
                    f"Failed to download season request {season_request.id} for show {show.name}"
                )
    log.info(f"Auto downloaded {count} approved season requests")
