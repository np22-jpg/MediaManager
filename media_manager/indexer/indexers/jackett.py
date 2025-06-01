import logging
import xml.etree.ElementTree as ET
from xml.etree.ElementTree import Element

import requests

from media_manager.indexer.indexers.generic import GenericIndexer
from media_manager.indexer.config import JackettConfig
from media_manager.indexer.schemas import IndexerQueryResult

log = logging.getLogger(__name__)


class Jackett(GenericIndexer):
    def __init__(self, **kwargs):
        """
        A subclass of GenericIndexer for interacting with the Jacket API.

        """
        super().__init__(name="jackett")
        config = JackettConfig()
        self.api_key = config.api_key
        self.url = config.url
        self.indexers = config.indexers
        log.debug("Registering Jacket as Indexer")

    # TODO: change architecture to build query string in the torrent module, instead of tv module
    # NOTE: this could be done in parallel, but if there aren't more than a dozen indexers, it shouldn't matter
    def search(self, query: str) -> list[IndexerQueryResult]:
        log.debug("Searching for " + query)

        responses = []
        for indexer in self.indexers:
            log.debug(f"Searching in indexer: {indexer}")
            url = (
                    self.url
                    + f"/api/v2.0/indexers/{indexer}/results/torznab/api?apikey={self.api_key}&t=tvsearch&q={query}"
            )
            response = requests.get(url)
            responses.append(response)

        xmlns = {
            "torznab": "http://torznab.com/schemas/2015/feed",
            "atom": "http://www.w3.org/2005/Atom",
        }
        result_list: list[IndexerQueryResult] = []
        for response in responses:
            if response.status_code == 200:
                xml_tree = ET.fromstring(response.content)
                for item in xml_tree.findall("channel/item"):
                    attributes: list[Element] = [
                        x for x in item.findall("torznab:attr", xmlns)
                    ]
                    for attribute in attributes:
                        if attribute.attrib["name"] == "seeders":
                            seeders = int(attribute.attrib["value"])
                            break
                    else:
                        log.warning(
                            f"Seeders not found in torrent: {item.find('title').text}, skipping this torrent"
                        )
                        continue

                    result = IndexerQueryResult(
                        title=item.find("title").text,
                        download_url=item.find("link").text,
                        seeders=seeders,
                        flags=[],
                        size=int(item.find("size").text),
                    )
                    result_list.append(result)
                    log.debug(f"Raw result: {result.model_dump()}")
            else:
                log.error(f"Jacket Error: {response.status_code}")
                return []
        return result_list
