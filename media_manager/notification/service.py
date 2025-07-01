import re

from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session

from media_manager.exceptions import InvalidConfigError
from media_manager.indexer.repository import IndexerRepository
from media_manager.database import SessionLocal
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.indexer.schemas import IndexerQueryResultId
from media_manager.metadataProvider.schemas import MetaDataProviderSearchResult
from media_manager.notification.repository import NotificationRepository
from media_manager.notification.schemas import NotificationId, Notification
from media_manager.torrent.schemas import Torrent, TorrentStatus
from media_manager.torrent.service import TorrentService
from media_manager.movies import log
from media_manager.movies.schemas import (
    Movie,
    MovieId,
    MovieRequest,
    MovieFile,
    RichMovieTorrent,
    PublicMovie,
    PublicMovieFile,
    MovieRequestId,
    RichMovieRequest,
)
from media_manager.torrent.schemas import QualityStrings
from media_manager.movies.repository import MovieRepository
from media_manager.exceptions import NotFoundError
import pprint
from media_manager.config import BasicConfig
from media_manager.torrent.repository import TorrentRepository
from media_manager.torrent.utils import import_file, import_torrent
from media_manager.indexer.service import IndexerService
from media_manager.metadataProvider.abstractMetaDataProvider import (
    AbstractMetadataProvider,
)
from media_manager.metadataProvider.tmdb import TmdbMetadataProvider
from media_manager.metadataProvider.tvdb import TvdbMetadataProvider


class NotificationService:
    def __init__(
        self,
        notification_repository: NotificationRepository,
    ):
        self.notification_repository = notification_repository

    def get_notification(self, id: NotificationId) -> Notification:
        return self.notification_repository.get_notification(id=id)

    def get_unread_notifications(self) -> list[Notification]:
        return self.notification_repository.get_unread_notifications()

    def get_all_notifications(self) -> list[Notification]:
        return self.notification_repository.get_all_notifications()

    def save_notification(self, notification: Notification) -> None:
        return self.notification_repository.save_notification(notification)

    def mark_notification_as_read(self, id: NotificationId) -> None:
        return self.notification_repository.mark_notification_as_read(id=id)

    def mark_notification_as_unread(self, id: NotificationId) -> None:
        return self.notification_repository.mark_notification_as_unread(id=id)

    def delete_notification(self, id: NotificationId) -> None:
        return self.notification_repository.delete_notification(id=id)

