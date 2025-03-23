import re
from enum import Enum

from pydantic import BaseModel, computed_field


class Quality(Enum):
    high = 1
    medium = 2
    low = 3
    very_low = 4
    unknown = 5


class IndexerQueryResult(BaseModel):
    title: str
    download_url: str
    seeders: int
    flags: list[str]

    # TODO: make system to detect quality more sophisticated
    @computed_field
    def quality(self) -> Quality:
        high_quality_pattern = r'\b(4k|4K)\b'
        medium_quality_pattern = r'\b(1080p|1080P)\b'
        low_quality_pattern = r'\b(720p|720P)\b'
        very_low_quality_pattern = r'\b(480p|480P|360p|360P)\b'

        if re.search(high_quality_pattern, self.title):
            return Quality.high
        elif re.search(medium_quality_pattern, self.title):
            return Quality.medium
        elif re.search(low_quality_pattern, self.title):
            return Quality.low
        elif re.search(very_low_quality_pattern, self.title):
            return Quality.very_low
        else:
            return Quality.unknown

    def __gt__(self, other) -> bool:
        if self.quality.value != other.quality.value:
            return self.quality.value > other.quality.value
        return self.seeders < other.seeders

    def __lt__(self, other) -> bool:
        if self.quality.value != other.quality.value:
            return self.quality.value < other.quality.value
        return self.seeders > other.seeders

    def download(self) -> str:
        """
        downloads a torrent file and returns the filepath
        """
        import requests
        url = self.download_url
        torrent_filepath = self.title + ".torrent"
        with open(torrent_filepath, 'wb') as out_file:
            content = requests.get(url).content
            out_file.write(content)
        return torrent_filepath

    @computed_field
    @property
    def season(self) -> list[int]:
        pattern = r"\b[sS](\d+)\b"
        matches = re.findall(pattern, self.title, re.IGNORECASE)
        if matches.__len__() == 2:
            result = []
            for i in range(int(matches[0]), int(matches[1]) + 1):
                result.append(i)
        elif matches.__len__() == 1:
            result = [int(matches[0])]
        else:
            result = []
        return result


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
