import logging

from sqlalchemy.orm import Session

from media_manager.indexer.models import IndexerQueryResult
from media_manager.indexer.schemas import (
    IndexerQueryResultId,
    IndexerQueryResult as IndexerQueryResultSchema,
)

log = logging.getLogger(__name__)

def get_result(
    result_id: IndexerQueryResultId, db: Session
) -> IndexerQueryResultSchema:
    return IndexerQueryResultSchema.model_validate(
        db.get(IndexerQueryResult, result_id)
    )


def save_result(
    result: IndexerQueryResultSchema, db: Session
) -> IndexerQueryResultSchema:
    log.debug("Saving indexer query result: %s", result)
    db.add(IndexerQueryResult(**result.model_dump()))
    db.commit()
    return result
