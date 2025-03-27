import logging

import requests

from config import ProwlarrConfig
from indexer import GenericIndexer, IndexerQueryResult

log = logging.getLogger(__name__)


class Prowlarr(GenericIndexer):
    def __init__(self, **kwargs):
        """
        A subclass of GenericIndexer for interacting with the Prowlarr API.

        :param api_key: The API key for authenticating requests to Prowlarr.
        :param kwargs: Additional keyword arguments to pass to the superclass constructor.
        """
        super().__init__(name='prowlarr')
        config = ProwlarrConfig()
        self.api_key = config.api_key
        self.url = config.url

    def get_search_results(self, query: str) -> list[IndexerQueryResult]:
        url = self.url + '/api/v1/search'
        headers = {
            'accept': 'application/json',
            'X-Api-Key': self.api_key
        }

        params = {
            'query': query,
        }

        response = requests.get(url, headers=headers, params=params)
        if response.status_code == 200:
            result_list: list[IndexerQueryResult] = []
            for result in response.json():
                if result['protocol'] == 'torrent':
                    log.debug("torrent result: " + result.__str__())
                    result_list.append(
                        IndexerQueryResult(
                            _download_url=result['downloadUrl'],
                            title=result['sortTitle'],
                            seeders=result['seeders'],
                            flags=result['indexerFlags'],
                        )
                    )
            return result_list
        else:
            print(f'Error: {response.status_code}')
            return []
