from datetime import datetime
from uuid import UUID

from sqlalchemy import ForeignKey, PrimaryKeyConstraint, UniqueConstraint, DateTime
from sqlalchemy.orm import Mapped, mapped_column, relationship

from media_manager.auth.db import User
from media_manager.database import Base
from media_manager.torrent.models import Quality


class Notification(Base):
    __tablename__ = "notification"

    id: Mapped[UUID] = mapped_column(primary_key=True)
    message: Mapped[str]
    read: Mapped[bool]
    timestamp = mapped_column(type=DateTime)

