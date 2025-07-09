import hashlib
import logging

import bencoder
import qbittorrentapi
import requests
from pydantic_settings import BaseSettings, SettingsConfigDict

from media_manager.config import BasicConfig
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.torrent.download_clients.abstractDownloadClient import (
    AbstractDownloadClient,
)
from media_manager.torrent.schemas import TorrentStatus, Torrent

log = logging.getLogger(__name__)


class QbittorrentConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="QBITTORRENT_")
    host: str = "localhost"
    port: int = 8080
    username: str = "admin"
    password: str = "admin"


class QbittorrentDownloadClient(AbstractDownloadClient):
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

    def __init__(self):
        self.config = QbittorrentConfig()
        self.api_client = qbittorrentapi.Client(**self.config.model_dump())
        try:
            self.api_client.auth_log_in()
            log.info("Successfully logged into qbittorrent")
        except Exception as e:
            log.error(f"Failed to log into qbittorrent: {e}")
            raise
        finally:
            self.api_client.auth_log_out()

    def download_torrent(self, indexer_result: IndexerQueryResult) -> Torrent:
        """
        Add a torrent to the download client and return the torrent object.

        :param indexer_result: The indexer query result of the torrent file to download.
        :return: The torrent object with calculated hash and initial status.
        """
        log.info(f"Attempting to download torrent: {indexer_result.title}")

        torrent_filepath = (
            BasicConfig().torrent_directory / f"{indexer_result.title}.torrent"
        )

        if torrent_filepath.exists():
            log.warning(f"Torrent already exists: {torrent_filepath}")
            # Calculate hash from existing file
            with open(torrent_filepath, "rb") as file:
                content = file.read()
                decoded_content = bencoder.decode(content)
                torrent_hash = hashlib.sha1(
                    bencoder.encode(decoded_content[b"info"])
                ).hexdigest()
        else:
            # Download the torrent file
            with open(torrent_filepath, "wb") as file:
                content = requests.get(str(indexer_result.download_url)).content
                file.write(content)

            # Calculate hash and add to qBittorrent
            with open(torrent_filepath, "rb") as file:
                content = file.read()
                try:
                    decoded_content = bencoder.decode(content)
                except Exception as e:
                    log.error(f"Failed to decode torrent file: {e}")
                    raise e

                torrent_hash = hashlib.sha1(
                    bencoder.encode(decoded_content[b"info"])
                ).hexdigest()

                try:
                    self.api_client.auth_log_in()
                    answer = self.api_client.torrents_add(
                        category="MediaManager",
                        torrent_files=content,
                        save_path=indexer_result.title,
                    )
                finally:
                    self.api_client.auth_log_out()

            if answer != "Ok.":
                log.error(f"Failed to download torrent. API response: {answer}")
                raise RuntimeError(
                    f"Failed to download torrent, API-Answer isn't 'Ok.'; API Answer: {answer}"
                )

        log.info(f"Successfully processed torrent: {indexer_result.title}")

        # Create and return torrent object
        torrent = Torrent(
            status=TorrentStatus.unknown,
            title=indexer_result.title,
            quality=indexer_result.quality,
            imported=False,
            hash=torrent_hash,
        )

        # Get initial status from download client
        torrent.status = self.get_torrent_status(torrent)

        return torrent

    def remove_torrent(self, torrent: Torrent, delete_data: bool = False) -> None:
        """
        Remove a torrent from the download client.

        :param torrent: The torrent to remove.
        :param delete_data: Whether to delete the downloaded data.
        """
        log.info(f"Removing torrent: {torrent.title}")
        try:
            self.api_client.auth_log_in()
            self.api_client.torrents_delete(
                torrent_hashes=torrent.hash, delete_files=delete_data
            )
        finally:
            self.api_client.auth_log_out()

    def get_torrent_status(self, torrent: Torrent) -> TorrentStatus:
        """
        Get the status of a specific torrent.

        :param torrent: The torrent to get the status of.
        :return: The status of the torrent.
        """
        log.info(f"Fetching status for torrent: {torrent.title}")
        try:
            self.api_client.auth_log_in()
            info = self.api_client.torrents_info(torrent_hashes=torrent.hash)
        finally:
            self.api_client.auth_log_out()

        if not info:
            log.warning(f"No information found for torrent: {torrent.id}")
            return TorrentStatus.unknown
        else:
            state: str = info[0]["state"]
            log.info(f"Torrent {torrent.id} is in state: {state}")

            if state in self.DOWNLOADING_STATE:
                return TorrentStatus.downloading
            elif state in self.FINISHED_STATE:
                return TorrentStatus.finished
            elif state in self.ERROR_STATE:
                return TorrentStatus.error
            elif state in self.UNKNOWN_STATE:
                return TorrentStatus.unknown
            else:
                return TorrentStatus.error

    def pause_torrent(self, torrent: Torrent) -> None:
        """
        Pause a torrent download.

        :param torrent: The torrent to pause.
        """
        log.info(f"Pausing torrent: {torrent.title}")
        try:
            self.api_client.auth_log_in()
            self.api_client.torrents_pause(torrent_hashes=torrent.hash)
        finally:
            self.api_client.auth_log_out()

    def resume_torrent(self, torrent: Torrent) -> None:
        """
        Resume a torrent download.

        :param torrent: The torrent to resume.
        """
        log.info(f"Resuming torrent: {torrent.title}")
        try:
            self.api_client.auth_log_in()
            self.api_client.torrents_resume(torrent_hashes=torrent.hash)
        finally:
            self.api_client.auth_log_out()
