import hashlib
import logging
import mimetypes
import pprint
import re
from pathlib import Path

import bencoder
import qbittorrentapi
import requests
from pydantic_settings import BaseSettings, SettingsConfigDict
from sqlalchemy.orm import Session

import media_manager.torrent.repository
import media_manager.tv.repository
import media_manager.tv.service
from media_manager.config import BasicConfig
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.torrent.repository import (
    get_seasons_files_of_torrent,
    get_show_of_torrent,
    save_torrent,
)
from media_manager.torrent.schemas import Torrent, TorrentStatus, TorrentId
from media_manager.torrent.utils import (
    list_files_recursively,
    get_torrent_filepath,
    import_file,
    extract_archives,
)
from media_manager.tv.schemas import SeasonFile, Show

log = logging.getLogger(__name__)


class TorrentServiceConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="QBITTORRENT_")
    host: str = "localhost"
    port: int = 8080
    username: str = "admin"
    password: str = "admin"


class TorrentService:
    DOWNLOADING_STATE = (
        "allocating",
        "downloading",
        "metaDL",
        "pausedDL",
        "queuedDL",
        "stalledDL",
        "checkingDL",
        "forcedDL",
        "moving",
    )
    FINISHED_STATE = (
        "uploading",
        "pausedUP",
        "queuedUP",
        "stalledUP",
        "checkingUP",
        "forcedUP",
    )
    ERROR_STATE = ("missingFiles", "error", "checkingResumeData")
    UNKNOWN_STATE = ("unknown",)
    api_client = qbittorrentapi.Client(**TorrentServiceConfig().model_dump())

    def __init__(self, db: Session):
        try:
            self.api_client.auth_log_in()
            log.info("Successfully logged into qbittorrent")
            self.db = db
        except Exception as e:
            log.error(f"Failed to log into qbittorrent: {e}")
            raise
        finally:
            self.api_client.auth_log_out()

    def download(self, indexer_result: IndexerQueryResult) -> Torrent:
        log.info(f"Attempting to download torrent: {indexer_result.title}")
        torrent = Torrent(
            status=TorrentStatus.unknown,
            title=indexer_result.title,
            quality=indexer_result.quality,
            imported=False,
            hash="",
        )

        url = indexer_result.download_url
        torrent_filepath = BasicConfig().torrent_directory / f"{torrent.title}.torrent"
        with open(torrent_filepath, "wb") as file:
            content = requests.get(url).content
            file.write(content)

        with open(torrent_filepath, "rb") as file:
            content = file.read()
            try:
                decoded_content = bencoder.decode(content)
            except Exception as e:
                log.error(f"Failed to decode torrent file: {e}")
                raise e
            torrent.hash = hashlib.sha1(
                bencoder.encode(decoded_content[b"info"])
            ).hexdigest()
            answer = self.api_client.torrents_add(
                category="MediaManager", torrent_files=content, save_path=torrent.title
            )

        if answer == "Ok.":
            log.info(f"Successfully added torrent: {torrent.title}")
            return self.get_torrent_status(torrent=torrent)
        else:
            log.error(f"Failed to download torrent. API response: {answer}")
            raise RuntimeError(
                f"Failed to download torrent, API-Answer isn't 'Ok.'; API Answer: {answer}"
            )

    def get_torrent_status(self, torrent: Torrent) -> Torrent:
        log.info(f"Fetching status for torrent: {torrent.title}")
        info = self.api_client.torrents_info(torrent_hashes=torrent.hash)

        if not info:
            log.warning(f"No information found for torrent: {torrent.id}")
            torrent.status = TorrentStatus.unknown
        else:
            state: str = info[0]["state"]
            log.info(f"Torrent {torrent.id} is in state: {state}")

            if state in self.DOWNLOADING_STATE:
                torrent.status = TorrentStatus.downloading
            elif state in self.FINISHED_STATE:
                torrent.status = TorrentStatus.finished
            elif state in self.ERROR_STATE:
                torrent.status = TorrentStatus.error
            elif state in self.UNKNOWN_STATE:
                torrent.status = TorrentStatus.unknown
            else:
                torrent.status = TorrentStatus.error
        save_torrent(db=self.db, torrent_schema=torrent)
        return torrent

    def cancel_download(self, torrent: Torrent, delete_files: bool = False) -> Torrent:
        """
        cancels download of a torrent

        :param delete_files: Deletes the downloaded files of the torrent too, deactivated by default
        :param torrent: the torrent to cancel
        """
        log.info(f"Cancelling download for torrent: {torrent.title}")
        self.api_client.torrents_delete(delete_files=delete_files)
        return self.get_torrent_status(torrent=torrent)

    def pause_download(self, torrent: Torrent) -> Torrent:
        """
        pauses download of a torrent

        :param torrent: the torrent to pause
        """
        log.info(f"Pausing download for torrent: {torrent.title}")
        self.api_client.torrents_pause(torrent_hashes=torrent.hash)
        return self.get_torrent_status(torrent=torrent)

    def resume_download(self, torrent: Torrent) -> Torrent:
        """
        resumes download of a torrent

        :param torrent: the torrent to resume
        """
        log.info(f"Resuming download for torrent: {torrent.title}")
        self.api_client.torrents_resume(torrent_hashes=torrent.hash)
        return self.get_torrent_status(torrent=torrent)

    def import_torrent(self, torrent: Torrent) -> Torrent:
        log.info(f"importing torrent {torrent}")

        # get all files, extract archives if necessary and get all files (extracted) files again
        all_files = list_files_recursively(path=get_torrent_filepath(torrent=torrent))
        log.debug(f"Found {len(all_files)} files downloaded by the torrent")
        extract_archives(all_files)
        all_files = list_files_recursively(path=get_torrent_filepath(torrent=torrent))

        # Filter videos and subtitles from all files
        video_files = []
        subtitle_files = []
        for file in all_files:
            file_type = mimetypes.guess_file_type(file)
            if file_type[0] is not None:
                if file_type[0].startswith("video"):
                    video_files.append(file)
                    log.debug(f"File is a video, it will be imported: {file}")
                elif file_type[0].startswith("text") and file.suffix == ".srt":
                    subtitle_files.append(file)
                    log.debug(f"File is a subtitle, it will be imported: {file}")
                else:
                    log.debug(
                        f"File is neither a video nor a subtitle, will not be imported: {file}"
                    )
        log.info(
            f"Importing these {len(video_files)} files:\n" + pprint.pformat(video_files)
        )

        # Fetch show and season information
        show: Show = get_show_of_torrent(db=self.db, torrent_id=torrent.id)
        show_file_path = (
            BasicConfig().tv_directory
            / f"{show.name} ({show.year})  [{show.metadata_provider}id-{show.external_id}]"
        )
        season_files: list[SeasonFile] = get_seasons_files_of_torrent(
            db=self.db, torrent_id=torrent.id
        )
        log.info(
            f"Found {len(season_files)} season files associated with torrent {torrent.title}"
        )

        # creating directories and hard linking files
        for season_file in season_files:
            season = media_manager.tv.service.get_season(
                db=self.db, season_id=season_file.season_id
            )
            season_path = show_file_path / Path(f"Season {season.number}")

            try:
                season_path.mkdir(parents=True)
            except FileExistsError:
                log.warning(f"Path already exists: {season_path}")

            for episode in season.episodes:
                episode_file_name = (
                    f"{show.name} S{season.number:02d}E{episode.number:02d}"
                )
                if season_file.file_path_suffix != "":
                    episode_file_name += f" - {season_file.file_path_suffix}"

                pattern = (
                    r".*[.]S0?"
                    + str(season.number)
                    + r"E0?"
                    + str(episode.number)
                    + r"[.].*"
                )
                subtitle_pattern = pattern + r"[.]([A-Za-z]{2})[.]srt"
                target_file_name = season_path / episode_file_name

                # import subtitles
                for subtitle_file in subtitle_files:
                    log.debug(
                        f"Searching for pattern {subtitle_pattern} in subtitle file: {subtitle_file.name}"
                    )
                    regex_result = re.search(subtitle_pattern, subtitle_file.name)
                    if regex_result:
                        language_code = regex_result.group(1)
                        log.debug(
                            f"Found matching pattern: {subtitle_pattern} in subtitle file: {subtitle_file.name},"
                            + f" extracted language code: {language_code}"
                        )
                        target_subtitle_file = target_file_name.with_suffix(
                            f".{language_code}.srt"
                        )
                        import_file(
                            target_file=target_subtitle_file, source_file=subtitle_file
                        )
                    else:
                        log.debug(
                            f"Didn't find any pattern {subtitle_pattern} in subtitle file: {subtitle_file.name}"
                        )

                # import episode videos
                for file in video_files:
                    log.debug(
                        f"Searching for pattern {pattern} in video file: {file.name}"
                    )
                    if re.search(pattern, file.name):
                        log.debug(
                            f"Found matching pattern: {pattern} in file {file.name}"
                        )
                        target_video_file = target_file_name.with_suffix(file.suffix)
                        import_file(target_file=target_video_file, source_file=file)
                        break
                else:
                    log.warning(
                        f"S{season.number}E{episode.number} in Torrent {torrent.title}'s files not found."
                    )
        torrent.imported = True

        return self.get_torrent_status(torrent=torrent)

    def get_all_torrents(self) -> list[Torrent]:
        return [
            self.get_torrent_status(x)
            for x in media_manager.torrent.repository.get_all_torrents(db=self.db)
        ]

    def get_torrent_by_id(self, torrent_id: TorrentId) -> Torrent:
        return self.get_torrent_status(
            media_manager.torrent.repository.get_torrent_by_id(
                torrent_id=torrent_id, db=self.db
            )
        )

    def delete_torrent(self, torrent_id: TorrentId):
        t = media_manager.torrent.repository.get_torrent_by_id(
            torrent_id=torrent_id, db=self.db
        )
        if not t.imported:
            media_manager.tv.repository.remove_season_files_by_torrent_id(
                db=self.db, torrent_id=torrent_id
            )
        media_manager.torrent.repository.delete_torrent(db=self.db, torrent_id=t.id)

    def import_all_torrents(self) -> list[Torrent]:
        log.info("Importing all torrents")
        torrents = self.get_all_torrents()
        log.info("Found %d torrents to import", len(torrents))
        imported_torrents = []
        for t in torrents:
            if t.imported == False and t.status == TorrentStatus.finished:
                imported_torrents.append(self.import_torrent(t))
        log.info("Finished importing all torrents")
        return imported_torrents
