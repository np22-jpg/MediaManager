import logging
import os

import qbittorrentapi
from pydantic import BaseModel
from sqlmodel import Session

from database import engine
from database.tv import Season
from dowloadClients.genericDownloadClient import GenericDownloadClient

# Configure logging
log = logging.getLogger(__name__)


class QbittorrentConfig(BaseModel):
    host: str = os.getenv("QBITTORRENT_HOST") or "localhost"
    port: int = os.getenv("QBITTORRENT_PORT") or 8080
    username: str = os.getenv("QBITTORRENT_USERNAME") or "admin"
    password: str = os.getenv("QBITTORRENT_PASSWORD") or "adminadmin"


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

    def download(self, torrent: Season):
        log.info(f"Attempting to download torrent: {torrent.url} with tag {torrent.id}+{torrent.season_number}")
        answer = self.api_client.torrents_add(category="MediaManager",
                                              urls=torrent.url,
                                              tag=f"{torrent.id}+{torrent.season_number}")
        if answer == "Ok.":
            log.info(f"Successfully added torrent: {torrent.url}")
            return self.get_torrent_status(torrent=torrent)
        else:
            log.error(f"Failed to download torrent. API response: {answer}")
            raise RuntimeError(f"Failed to download torrent, API-Answer isn't 'Ok.'; API Answer: {answer}")

    def get_torrent_status(self, torrent: Season) -> Season:
        log.info(f"Fetching status for torrent: {torrent.id}+{torrent.season_number}")
        info = self.api_client.torrents_info(tag=f"{torrent.id}+{torrent.season_number}")

        if not info:
            log.warning(f"No information found for torrent: {torrent.id}+{torrent.season_number}")
            torrent.status = "error"
        else:
            state: str = info[0]["state"]
            log.info(f"Torrent {torrent.id} is in state: {state}")

            if state in self.DOWNLOADING_STATE:
                torrent.status = "downloading"
            elif state in self.FINISHED_STATE:
                torrent.status = "finished"
            elif state in self.ERROR_STATE:
                torrent.status = "error"
            else:
                torrent.status = "error"

        with Session(engine) as session:
            session.add(torrent)
            session.commit()
            log.info(f"Updated torrent {torrent.id} status to: {torrent.status} in database.")

        return torrent
