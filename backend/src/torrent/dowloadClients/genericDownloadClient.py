from abc import ABCMeta, abstractmethod

from torrent.schemas import TorrentBase


class GenericDownloadClient(metaclass=ABCMeta):
    name: str

    @abstractmethod
    def __init__(cls, name: str = None, **kwargs):
        super().__init__(**kwargs)
        if name is None:
            raise ValueError('name cannot be None')
        cls.name = name

    @abstractmethod
    def download(self, torrent: TorrentBase) -> TorrentBase:
        """
        downloads a torrent

        :param torrent: id of the torrent to download
        """
        raise NotImplementedError()

    @abstractmethod
    def get_torrent_status(self, torrent: TorrentBase) -> TorrentBase:
        """
        updates a torrents 'status' field

        :param torrent: id of the media to update
        """
        raise NotImplementedError()

    @abstractmethod
    def cancel_download(self, torrent: TorrentBase) -> TorrentBase:
        """
        cancels download of a torrent

        :param torrent: id of the torrent to download
        """
        raise NotImplementedError()

    @abstractmethod
    def pause_download(self, torrent: TorrentBase) -> TorrentBase:
        """
        pauses download of a torrent

        :param torrent: id of the torrent to download
        """
        raise NotImplementedError()
