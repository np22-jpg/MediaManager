import pprint

import tvdb_v4_official
import logging

from pydantic_settings import BaseSettings

import media_manager.metadataProvider.utils
from media_manager.exceptions import InvalidConfigError
from media_manager.metadataProvider.abstractMetaDataProvider import (
    AbstractMetadataProvider,
    register_metadata_provider,
)
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.tv.schemas import Episode, Season, Show


class TvdbConfig(BaseSettings):
    TVDB_API_KEY: str | None = None


log = logging.getLogger(__name__)


class TvdbMetadataProvider(AbstractMetadataProvider):
    name = "tvdb"

    tvdb_client: tvdb_v4_official.TVDB

    def __init__(self, api_key: str = None):
        config = TvdbConfig()
        if config.TVDB_API_KEY is None:
            raise InvalidConfigError("TVDB_API_KEY is not set")
        self.tvdb_client = tvdb_v4_official.TVDB(config.TVDB_API_KEY)

    def download_show_poster_image(self, show: Show) -> bool:
        show_metadata = self.tvdb_client.get_series_extended(show.external_id)

        if show_metadata["image"] is not None:
            media_manager.metadataProvider.utils.download_poster_image(
                storage_path=self.storage_path,
                poster_url=show_metadata["image"],
                show=show,
            )
            log.info("Successfully downloaded poster image for show " + show.name)
            return True
        else:
            log.warning(f"image for show {show.name} could not be downloaded")
            return False

    def get_show_metadata(self, id: int = None) -> Show:
        """

        :param id: the external id of the show
        :type id: int
        :return: returns a ShowMetadata object
        :rtype: ShowMetadata
        """
        series = self.tvdb_client.get_series_extended(id)
        seasons = []
        for season in series["seasons"]:
            s = self.tvdb_client.get_season_extended(season["id"])
            episodes = [
                Episode(
                    number=episode["number"],
                    external_id=episode["id"],
                    title=episode["name"],
                )
                for episode in s["episodes"]
            ]
            seasons.append(
                Season(
                    number=s["number"],
                    name="TVDB doesn't provide Season Names",
                    overview="TVDB doesn't provide Season Overviews",
                    external_id=s["id"],
                    episodes=episodes,
                )
            )
        try:
            year = series["year"]
        except KeyError:
            year = None
        # NOTE: the TVDB API is fucking shit and seems to be very poorly documentated, I can't for the life of me
        #  figure out which statuses this fucking api returns
        show = Show(
            name=series["name"],
            overview=series["overview"],
            year=year,
            external_id=series["id"],
            metadata_provider=self.name,
            seasons=seasons,
            ended=False,
        )

        return show

    def search_show(
        self, query: str | None = None
    ) -> list[MetaDataProviderShowSearchResult]:
        if query is None:
            results = self.tvdb_client.get_all_series()
        else:
            results = self.tvdb_client.search(query)
        formatted_results = []
        for result in results:
            try:
                if result["type"] == "series":
                    try:
                        year = result["year"]
                    except KeyError:
                        year = None

                    formatted_results.append(
                        MetaDataProviderShowSearchResult(
                            poster_path=result["image_url"],
                            overview=result["overview"],
                            name=result["name"],
                            external_id=result["tvdb_id"],
                            year=year,
                            metadata_provider=self.name,
                            added=False,
                            vote_average=None,
                        )
                    )
            except Exception as e:
                log.warning(f"Error processing search result {result}: {e}")
        return formatted_results