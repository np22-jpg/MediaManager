from sqlalchemy import select
from sqlalchemy.orm import Session

from media_manager.torrent.models import Torrent
from media_manager.torrent.schemas import TorrentId, Torrent as TorrentSchema
from media_manager.tv.models import SeasonFile, Show, Season
from media_manager.tv.schemas import SeasonFile as SeasonFileSchema, Show as ShowSchema


def get_seasons_files_of_torrent(
    db: Session, torrent_id: TorrentId
) -> list[SeasonFileSchema]:
    stmt = select(SeasonFile).where(SeasonFile.torrent_id == torrent_id)
    result = db.execute(stmt).scalars().all()
    return [SeasonFileSchema.model_validate(season_file) for season_file in result]


def get_show_of_torrent(db: Session, torrent_id: TorrentId) -> ShowSchema:
    stmt = (
        select(Show)
        .join(SeasonFile.season)
        .join(Season.show)
        .where(SeasonFile.torrent_id == torrent_id)
    )
    result = db.execute(stmt).unique().scalar_one_or_none()
    return ShowSchema.model_validate(result)


def save_torrent(db: Session, torrent_schema: TorrentSchema) -> TorrentSchema:
    db.merge(Torrent(**torrent_schema.model_dump()))
    db.commit()
    return TorrentSchema.model_validate(torrent_schema)


def get_all_torrents(db: Session) -> list[TorrentSchema]:
    stmt = select(Torrent)
    result = db.execute(stmt).scalars().all()

    return [TorrentSchema.model_validate(torrent_schema) for torrent_schema in result]


def get_torrent_by_id(db: Session, torrent_id: TorrentId) -> TorrentSchema:
    return TorrentSchema.model_validate(db.get(Torrent, torrent_id))


def delete_torrent(db: Session, torrent_id: TorrentId):
    db.delete(db.get(Torrent, torrent_id))
