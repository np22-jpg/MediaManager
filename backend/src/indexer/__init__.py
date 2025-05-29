import logging

from backend.src.indexer.config import ProwlarrConfig
from backend.src.indexer.indexers.generic import GenericIndexer, IndexerQueryResult
from backend.src.indexer.indexers.prowlarr import Prowlarr

log = logging.getLogger(__name__)

indexers: list[GenericIndexer] = []

if ProwlarrConfig().enabled:
    indexers.append(Prowlarr())
