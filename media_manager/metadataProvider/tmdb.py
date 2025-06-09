import logging

import tmdbsimple
from pydantic_settings import BaseSettings
from tmdbsimple import TV, TV_Seasons

import media_manager.metadataProvider.utils
from media_manager.metadataProvider.abstractMetaDataProvider import (
    AbstractMetadataProvider,
    register_metadata_provider,
)
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.tv.schemas import Episode, Season, Show, SeasonNumber, EpisodeNumber


class TmdbConfig(BaseSettings):
    TMDB_API_KEY: str | None = None


ENDED_STATUS = {"Ended", "Canceled"}

config = TmdbConfig()
log = logging.getLogger(__name__)


class TmdbMetadataProvider(AbstractMetadataProvider):
    name = "tmdb"

    def __init__(self, api_key: str = None):
        tmdbsimple.API_KEY = api_key

    def download_show_poster_image(self, show: Show) -> bool:
        show_metadata = TV(show.external_id).info()
        # downloading the poster
        # all pictures from TMDB should already be jpeg, so no need to convert
        if show_metadata["poster_path"] is not None:
            poster_url = (
                "https://image.tmdb.org/t/p/original" + show_metadata["poster_path"]
            )
            if media_manager.metadataProvider.utils.download_poster_image(
                storage_path=self.storage_path, poster_url=poster_url, show=show
            ):
                log.info("Successfully downloaded poster image for show " + show.name)
            else:
                log.warning(f"download for image of show {show.name} failed")
                return False
        else:
            log.warning(f"image for show {show.name} could not be downloaded")
            return False
        return True


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
            season_metadata = TV_Seasons(
                tv_id=show_metadata["id"], season_number=season["season_number"]
            ).info()
            episode_list = []

            for episode in season_metadata["episodes"]:
                episode_list.append(
                    Episode(
                        external_id=int(episode["id"]),
                        title=episode["name"],
                        number=EpisodeNumber(episode["episode_number"]),
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

        year = media_manager.metadataProvider.utils.get_year_from_first_air_date(
            show_metadata["first_air_date"]
        )

        show = Show(
            external_id=id,
            name=show_metadata["name"],
            overview=show_metadata["overview"],
            year=year,
            seasons=season_list,
            metadata_provider=self.name,
            ended=show_metadata["status"] in ENDED_STATUS,
        )

        return show

    def search_show(
        self, query: str | None = None, max_pages: int = 5
    ) -> list[MetaDataProviderShowSearchResult]:
        """
        Search for shows using TMDB API.
        If no query is provided, it will return the most popular shows.
        """
        if query is None:
            result_factory = lambda page: tmdbsimple.Trending(media_type="tv").info()
        else:
            result_factory = lambda page: tmdbsimple.Search().tv(
                page=page, query=query, include_adult=True
            )

        results = []
        for i in range(1, max_pages + 1):
            result_page = result_factory(i)

            if not result_page["results"]:
                break
            else:
                results.extend(result_page["results"])

        formatted_results = []
        for result in results:
            try:
                if result["poster_path"] is not None:
                    poster_url = (
                        "https://image.tmdb.org/t/p/original" + result["poster_path"]
                    )
                else:
                    poster_url = None
                formatted_results.append(
                    MetaDataProviderShowSearchResult(
                        poster_path=poster_url,
                        overview=result["overview"],
                        name=result["name"],
                        external_id=result["id"],
                        year=media_manager.metadataProvider.utils.get_year_from_first_air_date(
                            result["first_air_date"]
                        ),
                        metadata_provider=self.name,
                        added=False,
                        vote_average=result["vote_average"],
                    )
                )
            except Exception as e:
                log.warning(f"Error processing search result {result}: {e}")
        return formatted_results



if config.TMDB_API_KEY is not None:
    log.info("Registering TMDB as metadata provider")
    register_metadata_provider(
        metadata_provider=TmdbMetadataProvider(config.TMDB_API_KEY)
    )
