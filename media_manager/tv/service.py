import re

from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session

import media_manager.indexer.service
import media_manager.metadataProvider
import media_manager.torrent.repository
from media_manager.indexer.repository import IndexerRepository
from media_manager.database import SessionLocal
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.indexer.schemas import IndexerQueryResultId
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.torrent.schemas import Torrent, TorrentStatus
from media_manager.torrent.service import TorrentService
from media_manager.tv import log
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
    SeasonRequestId,
    RichSeasonRequest,
)
from media_manager.torrent.schemas import QualityStrings
from media_manager.tv.repository import TvRepository
from media_manager.exceptions import NotFoundError
import pprint
from pathlib import Path
from media_manager.config import BasicConfig
from media_manager.torrent.repository import TorrentRepository
from media_manager.torrent.utils import import_file, import_torrent
from media_manager.indexer.service import IndexerService


class TvService:
    def __init__(
        self,
        tv_repository: TvRepository,
        torrent_service: TorrentService,
        indexer_service: IndexerService,
    ):
        self.tv_repository = tv_repository
        self.torrent_service = torrent_service
        self.indexer_service = indexer_service

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

    def get_season_request_by_id(
        self, season_request_id: SeasonRequestId
    ) -> SeasonRequest | None:
        """
        Get a season request by its ID.

        :param season_request_id: The ID of the season request.
        :return: The season request or None if not found.
        """
        return self.tv_repository.get_season_request(
            season_request_id=season_request_id
        )

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

    def get_public_season_files_by_season_id(
        self, season_id: SeasonId
    ) -> list[PublicSeasonFile]:
        """
        Get all public season files for a given season ID.

        :param season_id: The ID of the season.
        :return: A list of public season files.
        """
        season_files = self.tv_repository.get_season_files_by_season_id(
            season_id=season_id
        )
        public_season_files = [PublicSeasonFile.model_validate(x) for x in season_files]
        result = []
        for season_file in public_season_files:
            if self.season_file_exists_on_file(season_file=season_file):
                season_file.downloaded = True
            result.append(season_file)
        return result

    def check_if_show_exists(
        self,
        external_id: int = None,
        metadata_provider: str = None,
        show_id: ShowId = None,
    ) -> bool:
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
                self.tv_repository.get_show_by_external_id(
                    external_id=external_id, metadata_provider=metadata_provider
                )
                return True
            except NotFoundError:
                return False
        elif show_id:
            try:
                self.tv_repository.get_show_by_id(show_id=show_id)
                return True
            except NotFoundError:
                return False
        else:
            raise ValueError(
                "External ID and metadata provider or Show ID must be provided"
            )

    def get_all_available_torrents_for_a_season(
        self, season_number: int, show_id: ShowId, search_query_override: str = None
    ) -> list[IndexerQueryResult]:
        """
        Get all available torrents for a given season.

        :param season_number: The number of the season.
        :param show_id: The ID of the show.
        :param search_query_override: Optional override for the search query.
        :return: A list of indexer query results.
        """
        log.debug(
            f"getting all available torrents for season {season_number} for show {show_id}"
        )
        show = self.tv_repository.get_show_by_id(show_id=show_id)
        if search_query_override:
            search_query = search_query_override
        else:
            # TODO: add more Search query strings and combine all the results, like "season 3", "s03", "s3"
            search_query = show.name + " s" + str(season_number).zfill(2)

        torrents: list[IndexerQueryResult] = self.indexer_service.search(
            query=search_query
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

    def get_all_shows(self) -> list[Show]:
        """
        Get all shows.

        :return: A list of all shows.
        """
        return self.tv_repository.get_shows()

    def search_for_show(
        self, query: str, metadata_provider: str
    ) -> list[MetaDataProviderShowSearchResult]:
        """
        Search for shows using a given query.

        :param query: The search query.
        :param metadata_provider: The metadata provider to search.
        :return: A list of metadata provider show search results.
        """
        results = media_manager.metadataProvider.search_show(query, metadata_provider)
        for result in results:
            if self.check_if_show_exists(
                external_id=result.external_id, metadata_provider=metadata_provider
            ):
                result.added = True
        return results

    def get_popular_shows(
        self, metadata_provider: str
    ) -> list[MetaDataProviderShowSearchResult]:
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
            if not self.check_if_show_exists(
                external_id=result.external_id, metadata_provider=metadata_provider
            ):
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
        season_files = self.tv_repository.get_season_files_by_season_id(
            season_id=season_id
        )
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
            torrent_file = self.torrent_service.get_torrent_by_id(
                torrent_id=season_file.torrent_id
            )
            if torrent_file.imported:
                print("Servas")
                return True
        return False

    def get_show_by_external_id(
        self, external_id: int, metadata_provider: str
    ) -> Show | None:
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
            seasons = self.tv_repository.get_seasons_by_torrent_id(
                torrent_id=show_torrent.id
            )
            season_files = self.torrent_service.get_season_files_of_torrent(
                torrent=show_torrent
            )
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

    def download_torrent(
        self,
        public_indexer_result_id: IndexerQueryResultId,
        show_id: ShowId,
        override_show_file_path_suffix: str = "",
    ) -> Torrent:
        """
        Download a torrent for a given indexer result and show.

        :param public_indexer_result_id: The ID of the indexer result.
        :param show_id: The ID of the show.
        :param override_show_file_path_suffix: Optional override for the file path suffix.
        :return: The downloaded torrent.
        """
        indexer_result = self.indexer_service.get_result(
            result_id=public_indexer_result_id
        )
        show_torrent = self.torrent_service.download(indexer_result=indexer_result)

        for season_number in indexer_result.season:
            season = self.tv_repository.get_season_by_number(
                season_number=season_number, show_id=show_id
            )
            season_file = SeasonFile(
                season_id=season.id,
                quality=indexer_result.quality,
                torrent_id=show_torrent.id,
                file_path_suffix=override_show_file_path_suffix,
            )
            self.tv_repository.add_season_file(season_file=season_file)
        return show_torrent

    def download_approved_season_request(
        self, season_request: SeasonRequest, show: Show
    ) -> bool:
        """
        Download an approved season request.

        :param season_request: The season request to download.
        :param show: The Show object.
        :return: True if the download was successful, False otherwise.
        :raises ValueError: If the season request is not authorized.
        """
        if not season_request.authorized:
            log.error(
                f"Season request {season_request.id} is not authorized for download"
            )
            raise ValueError(
                f"Season request {season_request.id} is not authorized for download"
            )

        log.info(f"Downloading approved season request {season_request.id}")

        season = self.get_season(season_id=season_request.season_id)
        torrents = self.get_all_available_torrents_for_a_season(
            season_number=season.number, show_id=show.id
        )
        available_torrents: list[IndexerQueryResult] = []

        for torrent in torrents:
            if (
                (torrent.quality.value < season_request.wanted_quality.value)
                or (torrent.quality.value > season_request.min_quality.value)
                or (torrent.seeders < 3)
            ):
                log.info(
                    f"Skipping torrent {torrent.title} with quality {torrent.quality} for season {season.id}, because it does not match the requested quality {season_request.wanted_quality}"
                )
            elif torrent.season != [season.number]:
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

        torrent = self.torrent_service.download(indexer_result=available_torrents[0])
        season_file = SeasonFile(
            season_id=season.id,
            quality=torrent.quality,
            torrent_id=torrent.id,
            file_path_suffix=QualityStrings[torrent.quality.name].value.upper(),
        )
        try:
            self.tv_repository.add_season_file(season_file=season_file)
        except IntegrityError:
            log.warning(
                f"Season file for season {season.id} and quality {torrent.quality} already exists, skipping."
            )
        self.delete_season_request(season_request.id)
        return True

    def import_torrent_files(self, torrent: Torrent, show: Show) -> None:
        """
        Organizes files from a torrent into the TV directory structure, mapping them to seasons and episodes.
        :param torrent: The Torrent object
        :param show: The Show object
        """

        video_files, subtitle_files = import_torrent(torrent=torrent)

        log.info(
            f"Importing these {len(video_files)} files:\n" + pprint.pformat(video_files)
        )

        show_file_path = (
            BasicConfig().tv_directory
            / f"{show.name} ({show.year})  [{show.metadata_provider}id-{show.external_id}]"
        )
        season_files = self.torrent_service.get_season_files_of_torrent(torrent=torrent)
        log.info(
            f"Found {len(season_files)} season files associated with torrent {torrent.title}"
        )

        for season_file in season_files:
            season = self.get_season(season_id=season_file.season_id)
            season_path = show_file_path / Path(f"Season {season.number}")
            try:
                season_path.mkdir(parents=True, exist_ok=True)
            except Exception as e:
                log.warning(f"Could not create path {season_path}: {e}")
            for episode in season.episodes:
                episode_file_name = (
                    f"{show.name} S{season.number:02d}E{episode.number:02d}"
                )
                if season_file.file_path_suffix != "":
                    episode_file_name += f" - {season_file.file_path_suffix}"
                pattern = (
                    r".*[.]S0?"
                    + str(season.number)
                    + r"E0?"
                    + str(episode.number)
                    + r"[.].*"
                )
                subtitle_pattern = pattern + r"[.]([A-Za-z]{2})[.]srt"
                target_file_name = season_path / episode_file_name

                # import subtitles
                for subtitle_file in subtitle_files:
                    log.debug(
                        f"Searching for pattern {subtitle_pattern} in subtitle file: {subtitle_file.name}"
                    )
                    regex_result = re.search(subtitle_pattern, subtitle_file.name)
                    if regex_result:
                        language_code = regex_result.group(1)
                        log.debug(
                            f"Found matching pattern: {subtitle_pattern} in subtitle file: {subtitle_file.name},"
                            + f" extracted language code: {language_code}"
                        )
                        target_subtitle_file = target_file_name.with_suffix(
                            f".{language_code}.srt"
                        )
                        import_file(
                            target_file=target_subtitle_file, source_file=subtitle_file
                        )
                    else:
                        log.debug(
                            f"Didn't find any pattern {subtitle_pattern} in subtitle file: {subtitle_file.name}"
                        )

                # import episode videos
                for file in video_files:
                    log.debug(
                        f"Searching for pattern {pattern} in video file: {file.name}"
                    )
                    if re.search(pattern, file.name):
                        log.debug(
                            f"Found matching pattern: {pattern} in file {file.name}"
                        )
                        target_video_file = target_file_name.with_suffix(file.suffix)
                        import_file(target_file=target_video_file, source_file=file)
                        break
                else:
                    # TODO: notify admin that no video file was found for this episode
                    log.warning(
                        f"S{season.number}E{episode.number} in Torrent {torrent.title}'s files not found."
                    )
        log.info(f"Finished organizing files for torrent {torrent.title}")


def auto_download_all_approved_season_requests() -> None:
    """
    Auto download all approved season requests.
    This is a standalone function as it creates its own DB session.
    """
    db: Session = SessionLocal()
    tv_repository = TvRepository(db=db)
    torrent_service = TorrentService(torrent_repository=TorrentRepository(db=db))
    indexer_service = IndexerService(indexer_repository=IndexerRepository(db=db))
    tv_service = TvService(
        tv_repository=tv_repository,
        torrent_service=torrent_service,
        indexer_service=indexer_service,
    )

    log.info("Auto downloading all approved season requests")
    season_requests = tv_repository.get_season_requests()
    log.info(f"Found {len(season_requests)} season requests to process")
    log.debug(f"Season requests:  {[x.model_dump() for x in season_requests]}")
    count = 0

    for season_request in season_requests:
        if season_request.authorized:
            log.info(f"Processing season request {season_request.id} for download")
            show = tv_repository.get_show_by_season_id(
                season_id=season_request.season_id
            )
            if tv_service.download_approved_season_request(
                season_request=season_request, show=show
            ):
                count += 1
            else:
                log.warning(
                    f"Failed to download season request {season_request.id} for show {show.name}"
                )

    log.info(f"Auto downloaded {count} approved season requests")
    db.close()


def import_all_torrents() -> None:
    db: Session = SessionLocal()
    tv_repository = TvRepository(db=db)
    torrent_service = TorrentService(torrent_repository=TorrentRepository(db=db))
    indexer_service = IndexerService(indexer_repository=IndexerRepository(db=db))
    tv_service = TvService(
        tv_repository=tv_repository,
        torrent_service=torrent_service,
        indexer_service=indexer_service,
    )
    log.info("Importing all torrents")
    torrents = torrent_service.get_all_torrents()
    log.info("Found %d torrents to import", len(torrents))
    imported_torrents = []
    for t in torrents:
        if t.imported == False and t.status == TorrentStatus.finished:
            show = torrent_service.get_show_of_torrent(torrent=t)
            if show is None:
                log.warning(
                    f"torrent {t.title} is not a tv torrent, skipping import."
                )
                continue
            imported_torrents.append(tv_service.import_torrent_files(torrent=t, show=show))
    log.info("Finished importing all torrents")