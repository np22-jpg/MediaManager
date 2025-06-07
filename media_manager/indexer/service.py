from sqlalchemy.orm import Session

import media_manager.indexer.repository
from media_manager.indexer import log, indexers
from media_manager.indexer.repository import save_result
from media_manager.indexer.schemas import IndexerQueryResultId, IndexerQueryResult


def search(query: str, db: Session) -> list[IndexerQueryResult]:
    results = []

    log.debug(f"Searching for Torrent: {query}")

    for i in indexers:
        results.extend(i.search(query))
    for result in results:
        save_result(result=result, db=db)
    log.debug(f"Found Torrents: {results}")
    return results


def get_indexer_query_result(
    result_id: IndexerQueryResultId, db: Session
) -> IndexerQueryResult:
    return media_manager.indexer.repository.get_result(result_id=result_id, db=db)
