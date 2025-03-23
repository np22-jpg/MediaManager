import logging

import config
from indexer.generic import GenericIndexer, IndexerQueryResult
from indexer.prowlarr import Prowlarr

log = logging.getLogger(__name__)


def search(query: str) -> list[IndexerQueryResult]:
    results = []
    for indexer in indexers:
        results.extend(indexer.get_search_results(query))
    return results


indexers: list[GenericIndexer] = []

if config.ProwlarrConfig().enabled:
    indexers.append(Prowlarr())
