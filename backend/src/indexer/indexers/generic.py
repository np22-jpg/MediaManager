from backend.src.indexer.schemas import IndexerQueryResult


class GenericIndexer(object):
    name: str

    def __init__(self, name: str = None):
        if name:
            self.name = name
        else:
            raise ValueError('indexer name must not be None')

    def get_search_results(self, query: str) -> list[IndexerQueryResult]:
        """
        Sends a search request to the Indexer and returns the results.

        :param query: The search query to send to the Indexer.
        :return: A list of IndexerQueryResult objects representing the search results.
        """
        raise NotImplementedError()
