from dowloadClients.config import DownloadClientConfig
from dowloadClients.qbittorrent import QbittorrentClient

config = DownloadClientConfig()

if config.download_client == "qbit":
    client = QbittorrentClient()
else:
    raise ValueError("Unknown download client")
