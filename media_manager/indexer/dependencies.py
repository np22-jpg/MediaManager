from typing import Annotated

from fastapi import Depends, Path

from indexer.repository import IndexerRepository
from indexer.service import IndexerService
from media_manager.database import DbSessionDependency
from media_manager.tv.service import TvService


def get_indexer_repository(db_session: DbSessionDependency) -> IndexerRepository:
    return IndexerRepository(db_session)


indexer_repository_dep = Annotated[IndexerRepository, Depends(get_indexer_repository)]


def get_indexer_service(
    indexer_repository: IndexerRepository = indexer_repository_dep,
) -> IndexerService:
    return IndexerService(indexer_repository)


indexer_service_dep = Annotated[TvService, Depends(get_indexer_service)]