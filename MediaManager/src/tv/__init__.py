import logging

import tmdbsimple as tmdb

from config import TvConfig

config = TvConfig()
log = logging.getLogger(__name__)

tmdb.API_KEY = config.api_key
