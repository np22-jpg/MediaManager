import logging

from indexer.config import ProwlarrConfig
from indexer.indexers.generic import GenericIndexer, IndexerQueryResult
from indexer.indexers.prowlarr import Prowlarr

log = logging.getLogger(__name__)

indexers: list[GenericIndexer] = []

if ProwlarrConfig().enabled:
    indexers.append(Prowlarr())
