import logging

import metadataProvider.tmdb
from metadataProvider.abstractMetaDataProvider import metadata_providers
from metadataProvider.schemas import MetaDataProviderShowSearchResult
from tv.schemas import Show

log = logging.getLogger(__name__)


def get_show_metadata(id: int = None, provider: str = "tmdb") -> Show:
    if id is None or provider is None:
        raise ValueError("Show Metadata requires id and provider")
    return metadata_providers[provider].get_show_metadata(id)


def search_show(query: str, provider: str = "tmdb") -> list[MetaDataProviderShowSearchResult]:
    return metadata_providers[provider].search_show(query)
