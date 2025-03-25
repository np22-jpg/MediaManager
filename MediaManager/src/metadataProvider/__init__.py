import logging

import metadataProvider.tmdb
from database.tv import Show
from metadataProvider.abstractMetaDataProvider import metadata_providers

log = logging.getLogger(__name__)


def get_show_metadata(id: int = None, provider: str = "tmdb") -> Show:
    if id is None or provider is None:
        raise ValueError("Show Metadata requires id and provider")
    return metadata_providers[provider].get_show_metadata(id)
