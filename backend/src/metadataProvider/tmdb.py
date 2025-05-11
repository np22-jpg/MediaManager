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


class TmdbConfig(BaseSettings):
    TMDB_API_KEY: str | None = None


config = TmdbConfig()
log = logging.getLogger(__name__)


class TmdbMetadataProvider(AbstractMetadataProvider):
    name = "tmdb"

    def get_show_metadata(self, id: int = None) -> Show:
        """

        :param id: the external id of the show
        :type id: int
        :return: returns a ShowMetadata object
        :rtype: ShowMetadata
        """
        show_metadata = TV(id).info()
        season_list = []
        # inserting all the metadata into the objects
        for season in show_metadata["seasons"]:
            season_metadata = TV_Seasons(tv_id=show_metadata["id"], season_number=season["season_number"]).info()
            episode_list = []

            for episode in season_metadata["episodes"]:
                episode_list.append(
                    Episode(
                        external_id=int(episode["id"]),
                        title=episode["name"],
                        number=EpisodeNumber(episode["episode_number"])
                    )
                )

            season_list.append(
                Season(
                    external_id=int(season_metadata["id"]),
                    name=season_metadata["name"],
                    overview=season_metadata["overview"],
                    number=SeasonNumber(season_metadata["season_number"]),
                    episodes=episode_list,

                )
            )

        year = metadataProvider.utils.get_year_from_first_air_date(show_metadata["first_air_date"])

        show = Show(
            external_id=id,
            name=show_metadata["name"],
            overview=show_metadata["overview"],
            year=year,
            seasons=season_list,
            metadata_provider=self.name,
        )

        # TODO: convert images automatically to .jpg
        # downloading the poster
        if show_metadata["poster_path"] is not None:
            poster_url = "https://image.tmdb.org/t/p/original" + show_metadata["poster_path"]
            res = requests.get(poster_url, stream=True)
            content_type = res.headers["content-type"]
            file_extension = mimetypes.guess_extension(content_type)
            if res.status_code == 200:
                with open(self.storage_path.joinpath(str(show.id) + file_extension), 'wb') as f:
                    f.write(res.content)
                log.info(f"image for show {show.name} successfully downloaded")
        else:
            log.warning(f"image for show {show.name} could not be downloaded")

        return show

    def search_show(self, query: str) -> list[MetaDataProviderShowSearchResult]:
        results = tmdbsimple.Search().tv(query=query)
        formatted_results = []
        for result in results["results"]:
            if result["poster_path"] is not None:
                poster_url = "https://image.tmdb.org/t/p/original" + result["poster_path"]
            else:
                poster_url = None
            formatted_results.append(
                MetaDataProviderShowSearchResult(
                    poster_path=poster_url,
                    overview=result["overview"],
                    name=result["name"],
                    external_id=result["id"],
                    year=metadataProvider.utils.get_year_from_first_air_date(result["first_air_date"]),
                    metadata_provider=self.name,
                    added=False,
                )
            )
        return formatted_results

    def __init__(self, api_key: str = None):
        tmdbsimple.API_KEY = api_key


if config.TMDB_API_KEY is not None:
    log.info("Registering TMDB as metadata provider")
    register_metadata_provider(metadata_provider=TmdbMetadataProvider(config.TMDB_API_KEY))
