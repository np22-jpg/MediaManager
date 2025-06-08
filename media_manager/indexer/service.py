from sqlalchemy.orm import Session

import media_manager.indexer.repository
from media_manager.indexer import log, indexers
from media_manager.indexer.schemas import IndexerQueryResultId, IndexerQueryResult
from media_manager.tv.schemas import Show
from media_manager.indexer.repository import IndexerRepository


class IndexerService:
    def __init__(self, indexer_repository: IndexerRepository):
        self.repository = indexer_repository

    def get_result(self, result_id: IndexerQueryResultId) -> IndexerQueryResult:
        return self.repository.get_result(result_id=result_id)

    def search(
        self, query: str
    ) -> list[IndexerQueryResult]:
        """
        Search for results using the indexers based on a query.

        :param query: The search query.
        :param db: The database session.
        :return: A list of search results.
        """
        log.debug(f"Searching for: {query}")
        results = []

        for indexer in indexers:
            results.extend(indexer.search(query))

        for result in results:
            self.repository.save_result(result=result)

        log.debug(f"Found torrents: {results}")
        return results