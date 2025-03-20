import re
from typing import Literal

from pydantic import BaseModel, HttpUrl, computed_field


class IndexerQueryResult(BaseModel):
    title: str
    download_url: HttpUrl
    seeders: int
    flags: list[str]

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
        if self.seeders > other.seeders:
            return True
        else:
            return False

    def __lt__(self, other) -> bool:
        if self.seeders < other.seeders:
            return True
        else:
            return False


class GenericIndexer(object):
    url: str
    name: str

    def __init__(self, url: str = None, name: str = None):
        if url:
            self.url = url
        else:
            raise ValueError('indexer url must not be None')

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
