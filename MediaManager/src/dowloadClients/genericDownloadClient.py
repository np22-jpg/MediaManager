from dowloadClients import Torrents


class GenericDownloadClient(object):
    name: str

    def __init__(self, name: str = None):
        if name is None:
            raise ValueError('name cannot be None')
        self.name = name

    # TODO: change Torrents type to SeasonTorrents|MovieTorrents
    def download(self, torrent: Torrents) -> Torrents:
        """
        downloads a torrent

        :param torrent: object of the media to download
        """
        raise NotImplementedError()

    def get_torrent_status(self, torrent: Torrents) -> Torrents:
        """
        updates a torrents 'status' field

        :param torrent: object of the media to update
        """
        raise NotImplementedError()
