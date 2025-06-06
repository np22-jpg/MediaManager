from sqlalchemy.orm import Session

import media_manager.indexer.service
import media_manager.metadataProvider
import media_manager.torrent.repository
from media_manager.database import SessionLocal
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.indexer.schemas import IndexerQueryResultId
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.torrent.repository import get_seasons_files_of_torrent
from media_manager.torrent.schemas import Torrent
from media_manager.torrent.service import TorrentService
from media_manager.tv import log
from media_manager.tv.exceptions import MediaAlreadyExists
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
from media_manager.tv.repository import TvRepository


class TvService:
    def __init__(self, tv_repository: TvRepository):
        self.tv_repository = tv_repository

    def add_show(self, external_id: int, metadata_provider: str) -> Show | None:
        """
        Add a new show to the database.

        :param external_id: The ID of the show in the metadata provider's system.
        :param metadata_provider: The name of the metadata provider.
        :return: The saved show object or None if it failed.
        """
        show_with_metadata = media_manager.metadataProvider.get_show_metadata(
            id=external_id, provider=metadata_provider
        )
        saved_show = self.tv_repository.save_show(show=show_with_metadata)
        return saved_show

    def add_season_request(self, season_request: SeasonRequest) -> SeasonRequest:
        """
        Add a new season request.

        :param season_request: The season request to add.
        :return: The added season request.
        """
        return self.tv_repository.add_season_request(season_request=season_request)

    def get_season_request_by_id(self, season_request_id: SeasonRequestId) -> SeasonRequest | None:
        """
        Get a season request by its ID.

        :param season_request_id: The ID of the season request.
        :return: The season request or None if not found.
        """
        return self.tv_repository.get_season_request(season_request_id=season_request_id)

    def update_season_request(self, season_request: SeasonRequest) -> SeasonRequest:
        """
        Update an existing season request.

        :param season_request: The season request to update.
        :return: The updated season request.
        """
        self.tv_repository.delete_season_request(season_request_id=season_request.id)
        return self.tv_repository.add_season_request(season_request=season_request)

    def delete_season_request(self, season_request_id: SeasonRequestId) -> None:
        """
        Delete a season request by its ID.

        :param season_request_id: The ID of the season request to delete.
        """
        self.tv_repository.delete_season_request(season_request_id=season_request_id)

    def get_public_season_files_by_season_id(self, season_id: SeasonId) -> list[PublicSeasonFile]:
        """
        Get all public season files for a given season ID.

        :param season_id: The ID of the season.
        :return: A list of public season files.
        """
        season_files = self.tv_repository.get_season_files_by_season_id(season_id=season_id)
        public_season_files = [PublicSeasonFile.model_validate(x) for x in season_files]
        result = []
        for season_file in public_season_files:
            if self.season_file_exists_on_file(season_file=season_file):
                season_file.downloaded = True
            result.append(season_file)
        return result

    def get_public_season_files_by_season_number(self, season_number: SeasonNumber, show_id: ShowId) -> list[
        PublicSeasonFile]:
        """
        Get all public season files for a given season number and show ID.

        :param season_number: The number of the season.
        :param show_id: The ID of the show.
        :return: A list of public season files.
        """
        season = self.tv_repository.get_season_by_number(season_number=season_number, show_id=show_id)
        return self.get_public_season_files_by_season_id(season_id=season.id)

    def check_if_show_exists(self, external_id: int = None, metadata_provider: str = None,
                             show_id: ShowId = None) -> bool:
        """
        Check if a show exists in the database.

        :param external_id: The external ID of the show.
        :param metadata_provider: The metadata provider.
        :param show_id: The ID of the show.
        :return: True if the show exists, False otherwise.
        :raises ValueError: If neither external ID and metadata provider nor show ID are provided.
        """
        if external_id and metadata_provider:
            try:
                self.tv_repository.get_show_by_external_id(external_id=external_id, metadata_provider=metadata_provider)
                return True
            except:
                return False
        elif show_id:
            try:
                self.tv_repository.get_show_by_id(show_id=show_id)
                return True
            except:
                return False
        else:
            raise ValueError("External ID and metadata provider or Show ID must be provided")

    def get_all_available_torrents_for_a_season(self, season_number: int, show_id: ShowId,
                                                search_query_override: str = None) -> list[IndexerQueryResult]:
        """
        Get all available torrents for a given season.

        :param season_number: The number of the season.
        :param show_id: The ID of the show.
        :param search_query_override: Optional override for the search query.
        :return: A list of indexer query results.
        """
        log.debug(f"getting all available torrents for season {season_number} for show {show_id}")
        show = self.tv_repository.get_show_by_id(show_id=show_id)
        if search_query_override:
            search_query = search_query_override
        else:
            # TODO: add more Search query strings and combine all the results, like "season 3", "s03", "s3"
            search_query = show.name + " s" + str(season_number).zfill(2)

        torrents: list[IndexerQueryResult] = media_manager.indexer.service.search(query=search_query,
                                                                                  db=self.tv_repository.db)

        if search_query_override:
            log.debug(f"Found with search query override {torrents.__len__()} torrents: {torrents}")
            return torrents

        result: list[IndexerQueryResult] = []
        for torrent in torrents:
            if season_number in torrent.season:
                result.append(torrent)
        result.sort()
        return result

    def get_all_shows(self) -> list[Show]:
        """
        Get all shows.

        :return: A list of all shows.
        """
        return self.tv_repository.get_shows()

    def search_for_show(self, query: str, metadata_provider: str) -> list[MetaDataProviderShowSearchResult]:
        """
        Search for shows using a given query.

        :param query: The search query.
        :param metadata_provider: The metadata provider to search.
        :return: A list of metadata provider show search results.
        """
        results = media_manager.metadataProvider.search_show(query, metadata_provider)
        for result in results:
            if self.check_if_show_exists(external_id=result.external_id, metadata_provider=metadata_provider):
                result.added = True
        return results

    def get_popular_shows(self, metadata_provider: str) -> list[MetaDataProviderShowSearchResult]:
        """
        Get popular shows from a given metadata provider.

        :param metadata_provider: The metadata provider to use.
        :return: A list of metadata provider show search results.
        """
        results: list[MetaDataProviderShowSearchResult] = (
            media_manager.metadataProvider.search_show(provider=metadata_provider)
        )

        filtered_results = []
        for result in results:
            if not self.check_if_show_exists(external_id=result.external_id, metadata_provider=metadata_provider):
                filtered_results.append(result)

        return filtered_results

    def get_public_show_by_id(self, show_id: ShowId) -> PublicShow:
        """
        Get a public show by its ID.

        :param show_id: The ID of the show.
        :return: A public show.
        """
        show = self.tv_repository.get_show_by_id(show_id=show_id)
        seasons = [PublicSeason.model_validate(season) for season in show.seasons]
        for season in seasons:
            season.downloaded = self.is_season_downloaded(season_id=season.id)
        public_show = PublicShow.model_validate(show)
        public_show.seasons = seasons
        return public_show

    def get_show_by_id(self, show_id: ShowId) -> Show:
        """
        Get a show by its ID.

        :param show_id: The ID of the show.
        :return: The show.
        """
        return self.tv_repository.get_show_by_id(show_id=show_id)

    def is_season_downloaded(self, season_id: SeasonId) -> bool:
        """
        Check if a season is downloaded.

        :param season_id: The ID of the season.
        :return: True if the season is downloaded, False otherwise.
        """
        season_files = self.tv_repository.get_season_files_by_season_id(season_id=season_id)
        for season_file in season_files:
            if self.season_file_exists_on_file(season_file=season_file):
                return True
        return False

    def season_file_exists_on_file(self, season_file: SeasonFile) -> bool:
        """
        Check if a season file exists on the filesystem.

        :param season_file: The season file to check.
        :return: True if the file exists, False otherwise.
        """
        if season_file.torrent_id is None:
            return True
        else:
            torrent_file = media_manager.torrent.repository.get_torrent_by_id(
                db=self.tv_repository.db, torrent_id=season_file.torrent_id
            )
            if torrent_file.imported:
                return True
        return False

    def get_show_by_external_id(self, external_id: int, metadata_provider: str) -> Show | None:
        """
        Get a show by its external ID and metadata provider.

        :param external_id: The external ID of the show.
        :param metadata_provider: The metadata provider.
        :return: The show or None if not found.
        """
        return self.tv_repository.get_show_by_external_id(
            external_id=external_id, metadata_provider=metadata_provider
        )

    def get_season(self, season_id: SeasonId) -> Season:
        """
        Get a season by its ID.

        :param season_id: The ID of the season.
        :return: The season.
        """
        return self.tv_repository.get_season(season_id=season_id)

    def get_all_season_requests(self) -> list[RichSeasonRequest]:
        """
        Get all season requests.

        :return: A list of rich season requests.
        """
        return self.tv_repository.get_season_requests()

    def get_torrents_for_show(self, show: Show) -> RichShowTorrent:
        """
        Get torrents for a given show.

        :param show: The show.
        :return: A rich show torrent.
        """
        show_torrents = self.tv_repository.get_torrents_by_show_id(show_id=show.id)
        rich_season_torrents = []
        for show_torrent in show_torrents:
            seasons = self.tv_repository.get_seasons_by_torrent_id(torrent_id=show_torrent.id)
            season_files = get_seasons_files_of_torrent(db=self.tv_repository.db, torrent_id=show_torrent.id)
            file_path_suffix = season_files[0].file_path_suffix if season_files else ""
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

    def get_all_shows_with_torrents(self) -> list[RichShowTorrent]:
        """
        Get all shows with torrents.

        :return: A list of rich show torrents.
        """
        shows = self.tv_repository.get_all_shows_with_torrents()
        return [self.get_torrents_for_show(show=show) for show in shows]

    def download_torrent(self, public_indexer_result_id: IndexerQueryResultId, show_id: ShowId,
                         override_show_file_path_suffix: str = "") -> Torrent:
        """
        Download a torrent for a given indexer result and show.

        :param public_indexer_result_id: The ID of the indexer result.
        :param show_id: The ID of the show.
        :param override_show_file_path_suffix: Optional override for the file path suffix.
        :return: The downloaded torrent.
        """
        indexer_result = media_manager.indexer.service.get_indexer_query_result(
            db=self.tv_repository.db, result_id=public_indexer_result_id
        )
        show_torrent = TorrentService(db=self.tv_repository.db).download(indexer_result=indexer_result)

        for season_number in indexer_result.season:
            season = self.tv_repository.get_season_by_number(season_number=season_number, show_id=show_id)
            season_file = SeasonFile(
                season_id=season.id,
                quality=indexer_result.quality,
                torrent_id=show_torrent.id,
                file_path_suffix=override_show_file_path_suffix,
            )
            self.tv_repository.add_season_file(season_file=season_file)
        return show_torrent

    def download_approved_season_request(self, season_request: SeasonRequest, show_id: ShowId) -> bool:
        """
        Download an approved season request.

        :param season_request: The season request to download.
        :param show_id: The ID of the show.
        :return: True if the download was successful, False otherwise.
        :raises ValueError: If the season request is not authorized.
        """
        if not season_request.authorized:
            log.error(f"Season request {season_request.id} is not authorized for download")
            raise ValueError(f"Season request {season_request.id} is not authorized for download")

        log.info(f"Downloading approved season request {season_request.id}")

        season = self.get_season(season_id=season_request.season_id)
        torrents = self.get_all_available_torrents_for_a_season(season_number=season.number, show_id=show_id)
        available_torrents: list[IndexerQueryResult] = []

        for torrent in torrents:
            if (torrent.quality > season_request.wanted_quality or
                    torrent.quality < season_request.min_quality or
                    torrent.seeders < 3):
                log.info(
                    f"Skipping torrent {torrent.title} with quality {torrent.quality} for season {season.id}, because it does not match the requested quality {season_request.wanted_quality}"
                )
            elif torrent.season == [season.number]:
                log.info(
                    f"Skipping torrent {torrent.title} with quality {torrent.quality} for season {season.id}, because it contains to many/wrong seasons {torrent.season} (wanted: {season.number})"
                )
            else:
                available_torrents.append(torrent)
                log.info(
                    f"Taking torrent {torrent.title} with quality {torrent.quality} for season {season.id} into consideration"
                )

        if len(available_torrents) == 0:
            log.warning(
                f"No torrents matching criteria were found (wanted quality: {season_request.wanted_quality}, min_quality: {season_request.min_quality} for season {season.id})"
            )
            return False

        available_torrents.sort()

        torrent = TorrentService(db=self.tv_repository.db).download(indexer_result=available_torrents[0])
        season_file = SeasonFile(
            season_id=season.id,
            quality=torrent.quality,
            torrent_id=torrent.id,
            file_path_suffix=QualityStrings[torrent.quality.name].value.upper(),
        )
        self.tv_repository.add_season_file(season_file=season_file)
        return True


def auto_download_all_approved_season_requests() -> None:
    """
    Auto download all approved season requests.
    This is a standalone function as it creates its own DB session.
    """
    db: Session = SessionLocal()
    tv_repository = TvRepository(db)
    tv_service = TvService(tv_repository)

    log.info("Auto downloading all approved season requests")
    season_requests = tv_repository.get_season_requests()
    log.info(f"Found {len(season_requests)} season requests to process")
    count = 0

    for season_request in season_requests:
        if season_request.authorized:
            log.info(f"Processing season request {season_request.id} for download")
            show = tv_repository.get_show_by_season_id(season_id=season_request.season_id)
            if tv_service.download_approved_season_request(season_request=season_request, show_id=show.id):
                count += 1
            else:
                log.warning(f"Failed to download season request {season_request.id} for show {show.name}")

    log.info(f"Auto downloaded {count} approved season requests")
    db.close()
