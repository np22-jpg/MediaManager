from pydantic import BaseModel


class MetaDataProviderShowSearchResult(BaseModel):
    poster_path: str | None
    overview: str | None
    name: str
    external_id: int
    year: int | None
    metadata_provider: str
    added: bool
