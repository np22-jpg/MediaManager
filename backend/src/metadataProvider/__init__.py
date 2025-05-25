import logging

import metadataProvider.tmdb
import metadataProvider.tvdb
from metadataProvider.abstractMetaDataProvider import metadata_providers
from metadataProvider.schemas import MetaDataProviderShowSearchResult
from tv.schemas import Show

log = logging.getLogger(__name__)


def get_show_metadata(id: int = None, provider: str = "tmdb") -> Show:
    if id is None or provider is None:
        raise ValueError("Show Metadata requires id and provider")
    return metadata_providers[provider].get_show_metadata(id)


def search_show(query: str | None = None, provider: str = "tmdb") -> list[MetaDataProviderShowSearchResult]:
    """
    If no query is provided, it will return the most popular shows.
    """
    return metadata_providers[provider].search_show(query)

