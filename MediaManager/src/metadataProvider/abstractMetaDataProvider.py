import logging
from abc import ABC, abstractmethod

import config
from database.tv import Show

log = logging.getLogger(__name__)


class MetadataProvider(ABC):
    storage_path = config.BasicConfig().storage_directory
    @property
    @abstractmethod
    def name(self) -> str:
        pass

    @abstractmethod
    def get_show_metadata(self, id: int = None) -> Show:
        pass

    @abstractmethod
    def search_show(self, query):
        pass


metadata_providers = {}


def register_metadata_provider(metadata_provider: MetadataProvider):
    log.info("Registering metadata provider:" + metadata_provider.name)
    metadata_providers[metadata_provider.name] = metadata_provider
