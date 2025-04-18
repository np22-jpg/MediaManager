from torrent.dowloadClients.config import DownloadClientConfig
from torrent.dowloadClients.qbittorrent import QbittorrentClient

config = DownloadClientConfig()

if config.download_client == "qbit":
    client = QbittorrentClient()
else:
    raise ValueError("Unknown download client")
