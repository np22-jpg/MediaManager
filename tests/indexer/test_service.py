import uuid
from unittest.mock import MagicMock
import pytest

from media_manager.indexer.schemas import IndexerQueryResult, IndexerQueryResultId
from media_manager.indexer.repository import IndexerRepository
from media_manager.indexer.service import IndexerService


class DummyIndexer:
    def search(self, query):
        return [
            IndexerQueryResult(
                id=IndexerQueryResultId(uuid.uuid4()),
                title=f"{query} S01 1080p",
                download_url="http://example.com/1",
                seeders=10,
                flags=["test"],
                size=123456,
            )
        ]


@pytest.fixture
def mock_indexer_repository():
    repo = MagicMock(spec=IndexerRepository)
    repo.save_result.side_effect = lambda result: result
    return repo


@pytest.fixture
def indexer_service(monkeypatch, mock_indexer_repository):
    # Patch the global indexers list
    monkeypatch.setattr("media_manager.indexer.service.indexers", [DummyIndexer()])
    return IndexerService(indexer_repository=mock_indexer_repository)


def test_search_returns_results(indexer_service, mock_indexer_repository):
    query = "TestShow"
    results = indexer_service.search(query)
    assert len(results) == 1
    assert results[0].title == f"{query} S01 1080p"
    mock_indexer_repository.save_result.assert_called_once()


def test_get_result_returns_result(mock_indexer_repository):
    result_id = IndexerQueryResultId(uuid.uuid4())
    expected_result = IndexerQueryResult(
        id=result_id,
        title="Test S01 1080p",
        download_url="http://example.com/1",
        seeders=10,
        flags=["test"],
        size=123456,
    )
    mock_indexer_repository.get_result.return_value = expected_result
    service = IndexerService(indexer_repository=mock_indexer_repository)
    result = service.get_result(result_id)
    assert result == expected_result
    mock_indexer_repository.get_result.assert_called_once_with(result_id=result_id)
