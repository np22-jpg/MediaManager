import logging
from abc import ABC, abstractmethod

import media_manager.config
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.tv.schemas import Show

log = logging.getLogger(__name__)


class AbstractMetadataProvider(ABC):
    storage_path = media_manager.config.BasicConfig().image_directory

    @property
    @abstractmethod
    def name(self) -> str:
        pass

    @abstractmethod
    def get_show_metadata(self, id: int = None) -> Show:
        pass

    @abstractmethod
    def search_show(self, query) -> list[MetaDataProviderShowSearchResult]:
        pass


metadata_providers = {}


def register_metadata_provider(metadata_provider: AbstractMetadataProvider):
    log.info("Registering metadata provider:" + metadata_provider.name)
    metadata_providers[metadata_provider.name] = metadata_provider
