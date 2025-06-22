from typing import Annotated, Literal

from fastapi import Depends

from media_manager.exceptions import InvalidConfigError
from fastapi.exceptions import HTTPException
from media_manager.metadataProvider.tmdb import TmdbMetadataProvider
from media_manager.metadataProvider.abstractMetaDataProvider import (
    AbstractMetadataProvider,
)
from media_manager.metadataProvider.tvdb import TvdbMetadataProvider


def get_metadata_provider(
    metadata_provider: Literal["tmdb", "tvdb"] = "tmdb",
) -> AbstractMetadataProvider:
    try:
        if metadata_provider == "tmdb":
            return TmdbMetadataProvider()
        elif metadata_provider == "tvdb":
            return TvdbMetadataProvider()
        else:
            raise HTTPException(
                status_code=400,
                detail=f"Invalid metadata provider: {metadata_provider}. Supported providers are 'tmdb' and/or 'tvdb'.",
            )
    except InvalidConfigError as e:
        raise HTTPException(
            status_code=500,
            detail=f"Metadata provider '{metadata_provider}' not configured: {str(e)}",
        )


metadata_provider_dep = Annotated[
    AbstractMetadataProvider, Depends(get_metadata_provider)
]
