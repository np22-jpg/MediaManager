import logging

import config
from database.tv import Season
from indexer.generic import GenericIndexer, IndexerQueryResult
from indexer.prowlarr import Prowlarr

log = logging.getLogger(__name__)


def search(query: str | Season) -> list[IndexerQueryResult]:
    results = []

    if isinstance(query, Season):
        query = query.show.name + " s" + query.number.__str__()
        log.debug(f"Searching for Season {query}")

    for indexer in indexers:
        results.extend(indexer.get_search_results(query))

    return results


indexers: list[GenericIndexer] = []

if config.ProwlarrConfig().enabled:
    indexers.append(Prowlarr())
