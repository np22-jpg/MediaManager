from database.tv import Season


class GenericDownloadClient(object):
    name: str

    def __init__(self, name: str = None):
        if name is None:
            raise ValueError('name cannot be None')
        self.name = name

    # TODO: change Torrents type to SeasonTorrents|MovieTorrents
    def download(self, torrent: Season) -> Season:
        """
        downloads a torrent

        :param torrent: object of the media to download
        """
        raise NotImplementedError()

    def get_torrent_status(self, torrent: Season) -> Season:
        """
        updates a torrents 'status' field

        :param torrent: object of the media to update
        """
        raise NotImplementedError()
