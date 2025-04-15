from sqlalchemy.orm import Session

from indexer.models import IndexerQueryResult
from indexer.schemas import IndexerQueryResultId, IndexerQueryResult as IndexerQueryResultSchema


def get_result(result_id: IndexerQueryResultId, db: Session) -> IndexerQueryResultSchema:
    return IndexerQueryResultSchema(**db.get(IndexerQueryResult, result_id).__dict__)


def save_result(result: IndexerQueryResultSchema, db: Session) -> IndexerQueryResultSchema:
    db.add(IndexerQueryResult(**result.__dict__))
    return result
