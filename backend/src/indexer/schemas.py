import re
import typing
from uuid import UUID, uuid4

import pydantic
from pydantic import BaseModel, computed_field

from torrent.models import Quality, Torrent

IndexerQueryResultId = typing.NewType('IndexerQueryResultId', UUID)


# TODO: use something like strategy pattern to make sorting more user customizable
class IndexerQueryResult(BaseModel):
    id: IndexerQueryResultId = pydantic.Field(default_factory=uuid4)
    title: str
    download_url: str
    seeders: int
    flags: set[str]
    quality: Quality | None

    @computed_field
    @property
    def season(self) -> set[int]:
        pattern = r"\b[sS](\d+)\b"
        matches = re.findall(pattern, self.title, re.IGNORECASE)
        if matches.__len__() == 2:
            result = set()
            for i in range(int(matches[0]), int(matches[1]) + 1):
                result.add(i)
        elif matches.__len__() == 1:
            result = {int(matches[0])}
        else:
            result = {}
        return result

    def __gt__(self, other) -> bool:
        if self.quality.value != other.quality.value:
            return self.quality.value > other.quality.value
        return self.seeders < other.seeders

    def __lt__(self, other) -> bool:
        if self.quality.value != other.quality.value:
            return self.quality.value < other.quality.value
        return self.seeders > other.seeders

    def download(self) -> Torrent:
        """
        downloads a torrent file and returns the filepath
        """
        import requests
        url = self.download_url
        torrent_filepath = self.title + ".torrent"
        with open(torrent_filepath, 'wb') as out_file:
            content = requests.get(url).content
            out_file.write(content)

        return Torrent(status=None, title=self.title, quality=self.quality, id=self.id)
