from abc import ABC, abstractmethod

from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.torrent.schemas import TorrentId, TorrentStatus, Torrent


class AbstractDownloadClient(ABC):
    """
    Abstract base class for download clients.
    Defines the interface that all download clients must implement.
    """

    @abstractmethod
    def download_torrent(self, torrent: IndexerQueryResult) -> Torrent:
        """
        Add a torrent to the download client and return the torrent object.

        :param torrent: The indexer query result of the torrent file to download.
        :return: The torrent object with calculated hash and initial status.
        """
        pass

    @abstractmethod
    def remove_torrent(self, torrent: Torrent, delete_data: bool = False) -> None:
        """
        Remove a torrent from the download client.

        :param torrent: The torrent to remove.
        :param delete_data: Whether to delete the downloaded data.
        """
        pass

    @abstractmethod
    def get_torrent_status(self, torrent: Torrent) -> TorrentStatus:
        """
        Get the status of a specific torrent.

        :param torrent: The torrent to get the status of.
        :return: The status of the torrent.
        """
        pass