from sqlalchemy.orm import Session

from indexer.models import IndexerQueryResult
from indexer.schemas import (
    IndexerQueryResultId,
    IndexerQueryResult as IndexerQueryResultSchema,
)


def get_result(
        result_id: IndexerQueryResultId, db: Session
) -> IndexerQueryResultSchema:
    return IndexerQueryResultSchema.model_validate(
        db.get(IndexerQueryResult, result_id)
    )


def save_result(
        result: IndexerQueryResultSchema, db: Session
) -> IndexerQueryResultSchema:
    db.add(IndexerQueryResult(**result.model_dump()))
    return result
