import re
from typing import Literal

from pydantic import BaseModel, computed_field


class IndexerQueryResult(BaseModel):
    title: str
    download_url: str
    seeders: int
    flags: list[str]

    # TODO: make system to detect quality more sophisticated
    @computed_field
    @property
    def quality(self) -> Literal['high', 'medium', 'low']:
        high_quality_pattern = r'\b(4k|4K)\b'
        medium_quality_pattern = r'\b(1080p|1080P)\b'

        if re.search(high_quality_pattern, self.title):
            return 'high'
        elif re.search(medium_quality_pattern, self.title):
            return 'medium'
        else:
            return 'low'

    def __gt__(self, other) -> bool:
        if self.seeders < other.seeders:
            return True
        else:
            return False

    def __lt__(self, other) -> bool:
        if self.seeders > other.seeders:
            return True
        else:
            return False

    def download(self) -> str:
        import requests
        url = self.download_url
        torrent_filepath = self.title + ".torrent"
        with open(torrent_filepath, 'wb') as out_file:
            content = requests.get(url).content
            out_file.write(content)
        return torrent_filepath


class GenericIndexer(object):
    name: str

    def __init__(self, name: str = None):
        if name:
            self.name = name
        else:
            raise ValueError('indexer name must not be None')

    def get_search_results(self, query: str) -> list[IndexerQueryResult]:
        """
        Sends a search request to the Indexer and returns the results.

        :param query: The search query to send to the Indexer.
        :return: A list of IndexerQueryResult objects representing the search results.
        """
        raise NotImplementedError()
