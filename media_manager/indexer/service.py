from media_manager.indexer import log, indexers
from media_manager.indexer.schemas import IndexerQueryResultId, IndexerQueryResult
from media_manager.indexer.repository import IndexerRepository
from media_manager.notification.manager import notification_manager


class IndexerService:
    def __init__(self, indexer_repository: IndexerRepository):
        self.repository = indexer_repository

    def get_result(self, result_id: IndexerQueryResultId) -> IndexerQueryResult:
        return self.repository.get_result(result_id=result_id)

    def search(self, query: str) -> list[IndexerQueryResult]:
        """
        Search for results using the indexers based on a query.

        :param query: The search query.
        :param db: The database session.
        :return: A list of search results.
        """
        log.debug(f"Searching for: {query}")
        results = []
        failed_indexers = []

        for indexer in indexers:
            try:
                indexer_results = indexer.search(query)
                results.extend(indexer_results)
                log.debug(f"Indexer {indexer.__class__.__name__} returned {len(indexer_results)} results for query: {query}")
            except Exception as e:
                failed_indexers.append(indexer.__class__.__name__)
                log.error(f"Indexer {indexer.__class__.__name__} failed for query '{query}': {e}")

        # Send notification if indexers failed
        if failed_indexers and notification_manager.is_configured():
            notification_manager.send_notification(
                title="Indexer Failure",
                message=f"The following indexers failed for query '{query}': {', '.join(failed_indexers)}. Check indexer configuration and connectivity."
            )

        # Send notification if no results found from any indexer
        if not results and notification_manager.is_configured():
            notification_manager.send_notification(
                title="No Search Results",
                message=f"No torrents found for query '{query}' from any configured indexer. Consider checking the search terms or indexer availability."
            )

        for result in results:
            self.repository.save_result(result=result)

        log.debug(f"Found torrents: {results}")
        return results
