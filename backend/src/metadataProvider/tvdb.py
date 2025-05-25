import pprint

import tvdb_v4_official
import logging
import mimetypes

import requests
import tmdbsimple
from pydantic_settings import BaseSettings
from tmdbsimple import TV, TV_Seasons

import metadataProvider.utils
from metadataProvider.abstractMetaDataProvider import AbstractMetadataProvider, register_metadata_provider
from metadataProvider.schemas import MetaDataProviderShowSearchResult
from tv.schemas import Episode, Season, Show, SeasonNumber, EpisodeNumber


class TvdbConfig(BaseSettings):
    TVDB_API_KEY: str | None = None


config = TvdbConfig()
log = logging.getLogger(__name__)


class TvdbMetadataProvider(AbstractMetadataProvider):
    name = "tvdb"
    tvdb_client: tvdb_v4_official.TVDB

    def __init__(self, api_key: str = None):
        self.tvdb_client = tvdb_v4_official.TVDB(api_key)

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
            episodes = [Episode(number=episode['number'], external_id=episode['id'], title=episode['name']) for episode
                        in s["episodes"]]
            seasons.append(Season(number=s['number'], name="TVDB doesn't provide Season Names",
                                  overview="TVDB doesn't provide Season Overviews", external_id=s['id'],
                                  episodes=episodes))
        try:
            year = series['year']
        except KeyError:
            year = None
        show = Show(name=series['name'], overview=series['overview'], year=year,
                    external_id=series['id'], metadata_provider=self.name, seasons=seasons)

        if series["image"] is not None:
            metadataProvider.utils.download_poster_image(storage_path=self.storage_path,
                                                         poster_url=series['image'], show=show)
        else:
            log.warning(f"image for show {show.name} could not be downloaded")

        return show

    def search_show(self, query: str | None = None) -> list[MetaDataProviderShowSearchResult]:
        if query is None:
            results = self.tvdb_client.get_all_series()
        else:
            results = self.tvdb_client.search(query)
        formatted_results = []
        for result in results:
            try:
                if result['type'] == 'series':
                    try:
                        year = result['year']
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
                            vote_average=None
                        )
                    )
            except Exception as e:
                log.warning(f"Error processing search result {result}: {e}")
        return formatted_results


if config.TVDB_API_KEY is not None:
    log.info("Registering TVDB as metadata provider")
    register_metadata_provider(metadata_provider=TvdbMetadataProvider(config.TVDB_API_KEY))

if __name__ == "__main__":
    tvdb = TvdbMetadataProvider(config.TVDB_API_KEY)
    # show_metadata = tvdb.get_show_metadata(id=328724)  # Replace with a valid TVDB ID
    # pprint.pprint(dict(show_metadata))
    search_results = tvdb.search_show("Simpsons Declassified")
    pprint.pprint(search_results)
