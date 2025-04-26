from sqlalchemy.orm import Session

import indexer.repository
from indexer import IndexerQueryResult, log, indexers
from indexer.repository import save_result
from indexer.schemas import IndexerQueryResultId


def search(query: str, db: Session) -> list[IndexerQueryResult]:
    results = []

    log.debug(f"Searching for Torrent: {query}")

    for indexer in indexers:
        results.extend(indexer.get_search_results(query))
    for result in results:
        save_result(result=result, db=db)
    return results


def get_indexer_query_result(result_id: IndexerQueryResultId, db: Session) -> IndexerQueryResult:
    return indexer.repository.get_result(result_id=result_id, db=db)
