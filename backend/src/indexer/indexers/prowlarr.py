import logging

import requests

from backend.src.indexer import GenericIndexer
from backend.src.indexer.config import ProwlarrConfig
from backend.src.indexer.schemas import IndexerQueryResult

log = logging.getLogger(__name__)


class Prowlarr(GenericIndexer):
    def __init__(self, **kwargs):
        """
        A subclass of GenericIndexer for interacting with the Prowlarr API.

        :param api_key: The API key for authenticating requests to Prowlarr.
        :param kwargs: Additional keyword arguments to pass to the superclass constructor.
        """
        super().__init__(name="prowlarr")
        config = ProwlarrConfig()
        self.api_key = config.api_key
        self.url = config.url
        log.debug("Registering Prowlarr as Indexer")

    def get_search_results(self, query: str) -> list[IndexerQueryResult]:
        log.debug("Searching for " + query)
        url = self.url + "/api/v1/search"
        headers = {"accept": "application/json", "X-Api-Key": self.api_key}

        params = {
            "query": query,
        }

        response = requests.get(url, headers=headers, params=params)
        if response.status_code == 200:
            result_list: list[IndexerQueryResult] = []
            for result in response.json():
                if result["protocol"] == "torrent":
                    log.debug("torrent result: " + result.__str__())
                    result_list.append(
                        IndexerQueryResult(
                            download_url=result["downloadUrl"],
                            title=result["sortTitle"],
                            seeders=result["seeders"],
                            flags=result["indexerFlags"],
                            size=result["size"],
                        )
                    )
            return result_list
        else:
            log.error(f"Prowlarr Error: {response.status_code}")
            return []
