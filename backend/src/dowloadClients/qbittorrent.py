import logging

import qbittorrentapi
from pydantic_settings import BaseSettings, SettingsConfigDict

from database.torrents import Torrent
from dowloadClients.genericDownloadClient import GenericDownloadClient

log = logging.getLogger(__name__)


class QbittorrentConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix='QBITTORRENT_')
    host: str = "localhost"
    port: int = 8080
    username: str = "admin"


class QbittorrentClient(GenericDownloadClient):
    DOWNLOADING_STATE = ("allocating", "downloading", "metaDL", "pausedDL", "queuedDL", "stalledDL", "checkingDL",
                         "forcedDL", "moving")
    FINISHED_STATE = ("uploading", "pausedUP", "queuedUP", "stalledUP", "checkingUP", "forcedUP")
    ERROR_STATE = ("missingFiles", "error", "checkingResumeData", "unknown")
    api_client = qbittorrentapi.Client(**QbittorrentConfig().model_dump())

    def __init__(self):
        super().__init__(name="qBittorrent")
        try:
            self.api_client.auth_log_in()
            log.info("Successfully logged into qbittorrent")
        except Exception as e:
            log.error(f"Failed to log into qbittorrent: {e}")
            raise
        finally:
            self.api_client.auth_log_out()

    def download(self, torrent: Torrent) -> Torrent:
        log.info(f"Attempting to download torrent: {torrent.torrent_filepath} with tag {torrent.id}")

        with open(torrent.torrent_filepath, "rb") as torrent_file:
            answer = self.api_client.torrents_add(category="MediaManager",
                                                  torrent_files=torrent_file,
                                                  tags=[torrent.id.__str__()])
        if answer == "Ok.":
            log.info(f"Successfully added torrent: {torrent.torrent_filepath}")
            return self.get_torrent_status(torrent=torrent)
        else:
            log.error(f"Failed to download torrent. API response: {answer}")
            raise RuntimeError(f"Failed to download torrent, API-Answer isn't 'Ok.'; API Answer: {answer}")

    def get_torrent_status(self, torrent: Torrent) -> Torrent:
        log.info(f"Fetching status for torrent: {torrent.id}")
        info = self.api_client.torrents_info(tag=f"{torrent.id}")

        if not info:
            log.warning(f"No information found for torrent: {torrent.id}")
            torrent.torrent_status = "error"
        else:
            state: str = info[0]["state"]
            log.info(f"Torrent {torrent.id} is in state: {state}")

            if state in self.DOWNLOADING_STATE:
                torrent.torrent_status = "downloading"
            elif state in self.FINISHED_STATE:
                torrent.torrent_status = "finished"
            elif state in self.ERROR_STATE:
                torrent.torrent_status = "error"
            else:
                torrent.torrent_status = "error"

        return torrent
