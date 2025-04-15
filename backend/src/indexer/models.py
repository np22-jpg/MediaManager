from sqlalchemy.orm import Mapped, mapped_column

from database import Base
from indexer.schemas import IndexerQueryResultId
from torrent.schemas import Quality


class IndexerQueryResult(Base):
    __tablename__ = 'indexer_query_result'
    id: Mapped[IndexerQueryResultId] = mapped_column(primary_key=True)
    title: Mapped[str]
    download_url: Mapped[str]
    seeders: Mapped[int]
    flags: Mapped[set[str]]
    quality: Mapped[Quality | None]
    season: Mapped[set[int]]
