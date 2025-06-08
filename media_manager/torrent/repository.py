from sqlalchemy import select

from media_manager.database import DbSessionDependency
from media_manager.torrent.models import Torrent
from media_manager.torrent.schemas import TorrentId, Torrent as TorrentSchema
from media_manager.tv.models import SeasonFile, Show, Season
from media_manager.tv.schemas import SeasonFile as SeasonFileSchema, Show as ShowSchema
from tv.exceptions import NotFoundError


class TorrentRepository:
    def __init__(self, db: DbSessionDependency):
        self.db = db

    def get_seasons_files_of_torrent(
        self, torrent_id: TorrentId
    ) -> list[SeasonFileSchema]:
        stmt = select(SeasonFile).where(SeasonFile.torrent_id == torrent_id)
        result = self.db.execute(stmt).scalars().all()
        return [SeasonFileSchema.model_validate(season_file) for season_file in result]

    def get_show_of_torrent(self, torrent_id: TorrentId) -> ShowSchema | None:
        stmt = (
            select(Show)
            .join(SeasonFile.season)
            .join(Season.show)
            .where(SeasonFile.torrent_id == torrent_id)
        )
        result = self.db.execute(stmt).unique().scalar_one_or_none()
        if result is None:
            return None
        return ShowSchema.model_validate(result)

    def save_torrent(self, torrent: TorrentSchema) -> TorrentSchema:
        self.db.merge(Torrent(**torrent.model_dump()))
        self.db.commit()
        return TorrentSchema.model_validate(torrent)

    def get_all_torrents(self) -> list[TorrentSchema]:
        stmt = select(Torrent)
        result = self.db.execute(stmt).scalars().all()

        return [
            TorrentSchema.model_validate(torrent_schema) for torrent_schema in result
        ]

    def get_torrent_by_id(self, torrent_id: TorrentId) -> TorrentSchema:
        result = self.db.get(Torrent, torrent_id)
        if result is None:
            raise NotFoundError(f"Torrent with ID {torrent_id} not found.")
        return TorrentSchema.model_validate(result)

    def delete_torrent(self, torrent_id: TorrentId):
        self.db.delete(self.db.get(Torrent, torrent_id))
