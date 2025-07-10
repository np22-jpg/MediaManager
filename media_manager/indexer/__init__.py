import logging

from media_manager.config import AllEncompassingConfig
from media_manager.indexer.indexers.jackett import Jackett
from media_manager.indexer.indexers.generic import GenericIndexer
from media_manager.indexer.indexers.prowlarr import Prowlarr

log = logging.getLogger(__name__)

indexers: list[GenericIndexer] = []

config = AllEncompassingConfig()
if config.indexers.prowlarr.enabled:
    indexers.append(Prowlarr())
if config.indexers.jackett.enabled:
    indexers.append(Jackett())
