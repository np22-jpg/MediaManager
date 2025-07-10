import logging
from pydantic_settings import BaseSettings, SettingsConfigDict

from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.torrent.download_clients.abstractDownloadClient import (
    AbstractDownloadClient,
)
from media_manager.torrent.schemas import Torrent, TorrentStatus
import sabnzbd_api

log = logging.getLogger(__name__)


class SabnzbdConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="SABNZBD_")
    host: str = "localhost"
    port: int = 8080
    api_key: str = ""


class SabnzbdDownloadClient(AbstractDownloadClient):
    DOWNLOADING_STATE = (
        "Downloading",
        "Queued",
        "Paused",
        "Extracting",
        "Moving",
        "Running",
    )
    FINISHED_STATE = ("Completed",)
    ERROR_STATE = ("Failed",)
    UNKNOWN_STATE = ("Unknown",)

    def __init__(self):
        self.config = SabnzbdConfig()
        self.client = sabnzbd_api.SabnzbdClient(
            host=self.config.host,
            port=str(self.config.port),
            api_key=self.config.api_key,
        )
        self.client._base_url = f"{self.config.host.rstrip('/')}:{self.config.port}/api" # the library expects a /sabnzbd prefix for whatever reason
        try:
            # Test connection
            version = self.client.version()

            log.info(f"Successfully connected to SABnzbd version: {version}")
        except Exception as e:
            log.error(f"Failed to connect to SABnzbd: {e}")
            raise

    def download_torrent(self, indexer_result: IndexerQueryResult) -> Torrent:
        """
        Add a NZB/torrent to SABnzbd and return the torrent object.

        :param indexer_result: The indexer query result of the NZB file to download.
        :return: The torrent object with calculated hash and initial status.
        """
        log.info(f"Attempting to download NZB: {indexer_result.title}")

        try:
            # Add NZB to SABnzbd queue
            response = self.client.add_uri(
                url=str(indexer_result.download_url),
            )
            if not response["status"] == True:
                error_msg = response
                log.error(f"Failed to add NZB to SABnzbd: {error_msg}")
                raise RuntimeError(f"Failed to add NZB to SABnzbd: {error_msg}")

            # Generate a hash for the NZB (using title and download URL)
            nzo_id = response["nzo_ids"][0]

            log.info(f"Successfully added NZB: {indexer_result.title}")

            # Create and return torrent object
            torrent = Torrent(
                status=TorrentStatus.unknown,
                title=indexer_result.title,
                quality=indexer_result.quality,
                imported=False,
                hash=nzo_id,
                usenet=True,
            )

            # Get initial status from SABnzbd
            torrent.status = self.get_torrent_status(torrent)

            return torrent

        except Exception as e:
            log.error(f"Failed to download NZB {indexer_result.title}: {e}")
            raise

    def remove_torrent(self, torrent: Torrent, delete_data: bool = False) -> None:
        """
        Remove a torrent from SABnzbd.

        :param torrent: The torrent to remove.
        :param delete_data: Whether to delete the downloaded files.
        """
        log.info(f"Removing torrent: {torrent.title} (Delete data: {delete_data})")
        try:
            self.client.delete_job(nzo_id=torrent.hash, delete_files=delete_data)
            log.info(f"Successfully removed torrent: {torrent.title}")
        except Exception as e:
            log.error(f"Failed to remove torrent {torrent.title}: {e}")
            raise

    def pause_torrent(self, torrent: Torrent) -> None:
        """
        Pause a torrent in SABnzbd.

        :param torrent: The torrent to pause.
        """
        log.info(f"Pausing torrent: {torrent.title}")
        try:
            self.client.pause_job(nzo_id=torrent.hash)
            log.info(f"Successfully paused torrent: {torrent.title}")
        except Exception as e:
            log.error(f"Failed to pause torrent {torrent.title}: {e}")
            raise

    def resume_torrent(self, torrent: Torrent) -> None:
        """
        Resume a paused torrent in SABnzbd.

        :param torrent: The torrent to resume.
        """
        log.info(f"Resuming torrent: {torrent.title}")
        try:
            self.client.resume_job(nzo_id=torrent.hash)
            log.info(f"Successfully resumed torrent: {torrent.title}")
        except Exception as e:
            log.error(f"Failed to resume torrent {torrent.title}: {e}")
            raise

    def get_torrent_status(self, torrent: Torrent) -> TorrentStatus:
        """
        Get the status of a specific download from SABnzbd.

        :param torrent: The torrent to get the status of.
        :return: The status of the torrent.
        """
        log.info(f"Fetching status for download: {torrent.title}")
        response = self.client.get_downloads(nzo_ids=torrent.hash)
        log.debug("SABnzbd response: %s", response)
        status = response["queue"]["status"]
        log.info(f"Download status for NZB {torrent.title}: {status}")
        return self._map_status(status)

    def _map_status(self, sabnzbd_status: str) -> TorrentStatus:
        """
        Map SABnzbd status to TorrentStatus.

        :param sabnzbd_status: The status from SABnzbd.
        :return: The corresponding TorrentStatus.
        """
        if sabnzbd_status in self.DOWNLOADING_STATE:
            return TorrentStatus.downloading
        elif sabnzbd_status in self.FINISHED_STATE:
            return TorrentStatus.finished
        elif sabnzbd_status in self.ERROR_STATE:
            return TorrentStatus.error
        else:
            return TorrentStatus.unknown
