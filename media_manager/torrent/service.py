import hashlib
import logging

import bencoder
import qbittorrentapi
import requests
from pydantic_settings import BaseSettings, SettingsConfigDict

from media_manager.config import BasicConfig
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.torrent.repository import TorrentRepository
from media_manager.torrent.schemas import Torrent, TorrentStatus, TorrentId
from media_manager.tv.schemas import SeasonFile, Show
from media_manager.movies.schemas import Movie

log = logging.getLogger(__name__)


class TorrentServiceConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="QBITTORRENT_")
    host: str = "localhost"
    port: int = 8080
    username: str = "admin"
    password: str = "admin"


class TorrentService:
    DOWNLOADING_STATE = (
        "allocating",
        "downloading",
        "metaDL",
        "pausedDL",
        "queuedDL",
        "stalledDL",
        "checkingDL",
        "forcedDL",
        "moving",
    )
    FINISHED_STATE = (
        "uploading",
        "pausedUP",
        "queuedUP",
        "stalledUP",
        "checkingUP",
        "forcedUP",
    )
    ERROR_STATE = ("missingFiles", "error", "checkingResumeData")
    UNKNOWN_STATE = ("unknown",)
    api_client = qbittorrentapi.Client(**TorrentServiceConfig().model_dump())

    def __init__(self, torrent_repository: TorrentRepository):
        try:
            self.api_client.auth_log_in()
            log.info("Successfully logged into qbittorrent")
            self.torrent_repository = torrent_repository
        except Exception as e:
            log.error(f"Failed to log into qbittorrent: {e}")
            raise
        finally:
            self.api_client.auth_log_out()

    def get_season_files_of_torrent(self, torrent: Torrent) -> list[SeasonFile]:
        """
        Returns all season files of a torrent
        :param torrent: the torrent to get the season files of
        :return: list of season files
        """
        return self.torrent_repository.get_seasons_files_of_torrent(
            torrent_id=torrent.id
        )

    def get_show_of_torrent(self, torrent: Torrent) -> Show | None:
        """
        Returns the show of a torrent
        :param torrent: the torrent to get the show of
        :return: the show of the torrent
        """
        return self.torrent_repository.get_show_of_torrent(torrent_id=torrent.id)

    def get_movie_of_torrent(self, torrent: Torrent) -> Movie | None:
        """
        Returns the movie of a torrent
        :param torrent: the torrent to get the movie of
        :return: the movie of the torrent
        """
        return self.torrent_repository.get_movie_of_torrent(torrent_id=torrent.id)

    def download(self, indexer_result: IndexerQueryResult) -> Torrent:
        log.info(f"Attempting to download torrent: {indexer_result.title}")
        torrent = Torrent(
            status=TorrentStatus.unknown,
            title=indexer_result.title,
            quality=indexer_result.quality,
            imported=False,
            hash="",
        )

        url = indexer_result.download_url
        torrent_filepath = BasicConfig().torrent_directory / f"{torrent.title}.torrent"

        if torrent_filepath.exists():
            log.warning(f"Torrent already exists: {torrent_filepath}")
            return self.get_torrent_status(torrent=torrent)

        with open(torrent_filepath, "wb") as file:
            content = requests.get(url).content
            file.write(content)

        with open(torrent_filepath, "rb") as file:
            content = file.read()
            try:
                decoded_content = bencoder.decode(content)
            except Exception as e:
                log.error(f"Failed to decode torrent file: {e}")
                raise e
            torrent.hash = hashlib.sha1(
                bencoder.encode(decoded_content[b"info"])
            ).hexdigest()
            answer = self.api_client.torrents_add(
                category="MediaManager", torrent_files=content, save_path=torrent.title
            )

        if answer == "Ok.":
            log.info(f"Successfully added torrent: {torrent.title}")
            return self.get_torrent_status(torrent=torrent)
        else:
            log.error(f"Failed to download torrent. API response: {answer}")
            raise RuntimeError(
                f"Failed to download torrent, API-Answer isn't 'Ok.'; API Answer: {answer}"
            )

    def get_torrent_status(self, torrent: Torrent) -> Torrent:
        log.info(f"Fetching status for torrent: {torrent.title}")
        info = self.api_client.torrents_info(torrent_hashes=torrent.hash)

        if not info:
            log.warning(f"No information found for torrent: {torrent.id}")
            torrent.status = TorrentStatus.unknown
        else:
            state: str = info[0]["state"]
            log.info(f"Torrent {torrent.id} is in state: {state}")

            if state in self.DOWNLOADING_STATE:
                torrent.status = TorrentStatus.downloading
            elif state in self.FINISHED_STATE:
                torrent.status = TorrentStatus.finished
            elif state in self.ERROR_STATE:
                torrent.status = TorrentStatus.error
            elif state in self.UNKNOWN_STATE:
                torrent.status = TorrentStatus.unknown
            else:
                torrent.status = TorrentStatus.error
        self.torrent_repository.save_torrent(torrent=torrent)
        return torrent

    def cancel_download(self, torrent: Torrent, delete_files: bool = False) -> Torrent:
        """
        cancels download of a torrent

        :param delete_files: Deletes the downloaded files of the torrent too, deactivated by default
        :param torrent: the torrent to cancel
        """
        log.info(f"Cancelling download for torrent: {torrent.title}")
        self.api_client.torrents_delete(delete_files=delete_files)
        return self.get_torrent_status(torrent=torrent)

    def pause_download(self, torrent: Torrent) -> Torrent:
        """
        pauses download of a torrent

        :param torrent: the torrent to pause
        """
        log.info(f"Pausing download for torrent: {torrent.title}")
        self.api_client.torrents_pause(torrent_hashes=torrent.hash)
        return self.get_torrent_status(torrent=torrent)

    def resume_download(self, torrent: Torrent) -> Torrent:
        """
        resumes download of a torrent

        :param torrent: the torrent to resume
        """
        log.info(f"Resuming download for torrent: {torrent.title}")
        self.api_client.torrents_resume(torrent_hashes=torrent.hash)
        return self.get_torrent_status(torrent=torrent)

    def get_all_torrents(self) -> list[Torrent]:
        return [
            self.get_torrent_status(x)
            for x in self.torrent_repository.get_all_torrents()
        ]

    def get_torrent_by_id(self, torrent_id: TorrentId) -> Torrent:
        return self.get_torrent_status(
            self.torrent_repository.get_torrent_by_id(torrent_id=torrent_id)
        )

    # TODO: extract deletion logic to tv module
    # def delete_torrent(self, torrent_id: TorrentId):
    #    t = self.torrent_repository.get_torrent_by_id(torrent_id=torrent_id)
    #    if not t.imported:
    #        from media_manager.tv.repository import remove_season_files_by_torrent_id
    #        remove_season_files_by_torrent_id(db=self.db, torrent_id=torrent_id)
    #    media_manager.torrent.repository.delete_torrent(db=self.db, torrent_id=t.id)
    def get_movie_files_of_torrent(self, torrent: Torrent):
        return self.torrent_repository.get_movie_files_of_torrent(torrent_id=torrent.id)
