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
        raise NotImplementedError()

    @abstractmethod
    def search_show(
        self, query: str | None = None
    ) -> list[MetaDataProviderShowSearchResult]:
        raise NotImplementedError()

    @abstractmethod
    def download_show_poster_image(self, show: Show) -> bool:
        """
        Downloads the poster image for a show.
        :param show: The show to download the poster image for.
        :return: True if the image was downloaded successfully, False otherwise.
        """
        raise NotImplementedError()


metadata_providers = {}


def register_metadata_provider(metadata_provider: AbstractMetadataProvider):
    log.info("Registering metadata provider:" + metadata_provider.name)
    metadata_providers[metadata_provider.name] = metadata_provider
