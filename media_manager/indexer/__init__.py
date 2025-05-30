import logging

from media_manager.indexer.config import ProwlarrConfig
from media_manager.indexer.indexers.generic import GenericIndexer, IndexerQueryResult
from media_manager.indexer.indexers.prowlarr import Prowlarr

log = logging.getLogger(__name__)

indexers: list[GenericIndexer] = []

if ProwlarrConfig().enabled:
    indexers.append(Prowlarr())
