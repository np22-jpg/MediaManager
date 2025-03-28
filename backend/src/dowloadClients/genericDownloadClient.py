from database.tv import TorrentMixin


class GenericDownloadClient(object):
    name: str

    def __init__(self, name: str = None):
        if name is None:
            raise ValueError('name cannot be None')
        self.name = name

    def download(self, torrent: TorrentMixin) -> TorrentMixin:
        """
        downloads a torrent

        :param torrent: object of the media to download
        """
        raise NotImplementedError()

    def get_torrent_status(self, torrent: TorrentMixin) -> TorrentMixin:
        """
        updates a torrents 'status' field

        :param torrent: object of the media to update
        """
        raise NotImplementedError()
