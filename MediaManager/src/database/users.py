import uuid
from uuid import UUID

from sqlmodel import Field, SQLModel


class UserBase(SQLModel):
    name: str = Field()
    lastname: str
    email: str = Field(unique=True)

class UserPublic(UserBase):
    id: UUID = Field(primary_key=True, default_factory=uuid.uuid4)

class User(UserPublic, table=True):
    hashed_password: str

class UserCreate(UserBase):
    password: str