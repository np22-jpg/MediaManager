from typing import Literal

from pydantic import BaseModel, HttpUrl

from indexer.prowlarr import Prowlarr


class IndexerQueryResult(BaseModel):
    title: str
    download_url: HttpUrl
    seeders: int
    protocol: Literal["usenet", "torrent"]
    flags: list[str]

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
            raise ValueError("indexer url must not be None")

        if name:
            self.name = name
        else:
            raise ValueError("indexer name must not be None")

    def get_search_results(self, query: str) -> list[IndexerQueryResult]:
        """
        Sends a search request to the Indexer and returns the results.

        :param query: The search query to send to the Indexer.
        :return: A list of IndexerQueryResult objects representing the search results.
        """
        raise NotImplementedError()
