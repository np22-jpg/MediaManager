from config import DownloadClientConfig
from dowloadClients.qbittorrent import QbittorrentClient

config = DownloadClientConfig()

# TODO: add more elif when implementing more download clients
if config.client == "qbit":
    client = QbittorrentClient()
else:
    client = QbittorrentClient()
