from pydantic_settings import BaseSettings

from media_manager.torrent.download_clients.qbittorrent import QbittorrentConfig
from media_manager.torrent.download_clients.sabnzbd import SabnzbdConfig


class TorrentConfig(BaseSettings):
    qbittorrent: QbittorrentConfig = QbittorrentConfig()
    sabnzbd: SabnzbdConfig = SabnzbdConfig()
