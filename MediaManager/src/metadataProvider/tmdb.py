import logging
import mimetypes

import requests
import tmdbsimple
from tmdbsimple import TV, TV_Seasons

import config
from database.tv import Episode, Season, Show
from metadataProvider.abstractMetaDataProvider import MetadataProvider, register_metadata_provider

config = config.TmdbConfig()
log = logging.getLogger(__name__)


class TmdbMetadataProvider(MetadataProvider):
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
                        number=int(episode["episode_number"])
                    )
                )

            season_list.append(
                Season(
                    external_id=int(season_metadata["id"]),
                    name=season_metadata["name"],
                    overview=season_metadata["overview"],
                    number=int(season_metadata["season_number"]),
                    episodes=episode_list
                )
            )

        year: str | None = show_metadata["first_air_date"]
        if year:
            year: int = int(year.split('-')[0])
        else:
            year = None

        show = Show(
            external_id=id,
            name=show_metadata["name"],
            overview=show_metadata["overview"],
            year=year,
            seasons=season_list,
            metadata_provider=self.name,
        )

        # downloading the poster
        poster_url = "https://image.tmdb.org/t/p/original" + show_metadata["poster_path"]
        res = requests.get(poster_url, stream=True)
        content_type = res.headers["content-type"]
        file_extension = mimetypes.guess_extension(content_type)
        if res.status_code == 200:
            with open(f"{self.storage_path}/images/{show.id}{file_extension}", 'wb') as f:
                f.write(res.content)
            log.info(f"image for show {show.name} successfully downloaded")

        else:
            log.warning(f"image for show {show.name} could not be downloaded")

        return show

    def __init__(self, api_key: str = None):
        tmdbsimple.API_KEY = api_key


if config.api_key is not None:
    log.info("Registerng TMDB as metadata provider")
    register_metadata_provider(metadata_provider=TmdbMetadataProvider(config.api_key))
