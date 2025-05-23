from sqlalchemy import String, Integer
from sqlalchemy.dialects.postgresql import ARRAY
from sqlalchemy.orm import Mapped, mapped_column

from backend.src.database import Base
from indexer.schemas import IndexerQueryResultId
from torrent.schemas import Quality


class IndexerQueryResult(Base):
    __tablename__ = 'indexer_query_result'
    id: Mapped[IndexerQueryResultId] = mapped_column(primary_key=True)
    title: Mapped[str]
    download_url: Mapped[str]
    seeders: Mapped[int]
    flags = mapped_column(ARRAY(String))
    quality: Mapped[Quality]
    season = mapped_column(ARRAY(Integer))
    size = Mapped[int]
