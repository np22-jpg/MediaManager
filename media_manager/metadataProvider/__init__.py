import logging
from cachetools import TTLCache, cached

import media_manager.metadataProvider.tmdb as tmdb
import media_manager.metadataProvider.tvdb as tvdb
from media_manager.metadataProvider.abstractMetaDataProvider import metadata_providers
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.tv.schemas import Show

_ = tvdb, tmdb

log = logging.getLogger(__name__)
search_show_cache = TTLCache(maxsize=128, ttl=24 * 60 * 60)  # Cache for 24 hours


def get_show_metadata(id: int = None, provider: str = "tmdb") -> Show:
    if id is None or provider is None:
        raise ValueError("Show Metadata requires id and provider")
    return metadata_providers[provider].get_show_metadata(id)


def download_show_poster_image(show: Show) -> bool:
    """
    Downloads the poster image for a show.
    :param show: The show to download the poster image for.
    :return: True if the image was downloaded successfully, False otherwise.
    """
    return metadata_providers[show.metadata_provider].download_show_poster_image(show)

@cached(search_show_cache)
def search_show(
    query: str | None = None, provider: str = "tmdb"
) -> list[MetaDataProviderShowSearchResult]:
    """
    If no query is provided, it will return the most popular shows.
    """
    return metadata_providers[provider].search_show(query)
