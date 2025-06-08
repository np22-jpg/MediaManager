import logging
import mimetypes
from pathlib import Path

import patoolib

from media_manager.config import BasicConfig
from media_manager.torrent.schemas import Torrent

log = logging.getLogger(__name__)


def list_files_recursively(path: Path = Path(".")) -> list[Path]:
    files = list(path.glob("**/*"))
    log.debug(f"Found {len(files)} entries via glob")
    valid_files = []
    for x in files:
        if x.is_dir():
            log.debug(f"'{x}' is a directory")
        elif x.is_symlink():
            log.debug(f"'{x}' is a symlink")
        else:
            valid_files.append(x)
    log.debug(f"Returning {len(valid_files)} files after filtering")
    return valid_files


def extract_archives(files):
    for file in files:
        file_type = mimetypes.guess_type(file)
        log.debug(f"File: {file}, Size: {file.stat().st_size} bytes, Type: {file_type}")
        if file_type[0] == "application/x-compressed":
            log.debug(
                f"File {file} is a compressed file, extracting it into directory {file.parent}"
            )
            patoolib.extract_archive(str(file), outdir=str(file.parent))


def get_torrent_filepath(torrent: Torrent):
    return BasicConfig().torrent_directory / torrent.title


def import_file(target_file: Path, source_file: Path):
    if target_file.exists():
        target_file.unlink()
    target_file.hardlink_to(source_file)


def import_torrent(torrent: Torrent) -> (list[Path], list[Path], list[Path]):
    """
    Extracts all files from the torrent download directory, including extracting archives.
    Returns a tuple containing: seperated video files, subtitle files, and all files found in the torrent directory.
    """
    log.info(f"Importing torrent {torrent}")
    all_files = list_files_recursively(path=get_torrent_filepath(torrent=torrent))
    log.debug(f"Found {len(all_files)} files downloaded by the torrent")
    extract_archives(all_files)
    all_files = list_files_recursively(path=get_torrent_filepath(torrent=torrent))

    video_files = []
    subtitle_files = []
    for file in all_files:
        file_type, _ = mimetypes.guess_type(str(file))
        if file_type is not None:
            if file_type.startswith("video"):
                video_files.append(file)
                log.debug(f"File is a video, it will be imported: {file}")
            elif file_type.startswith("text") and Path(file).suffix == ".srt":
                subtitle_files.append(file)
                log.debug(f"File is a subtitle, it will be imported: {file}")
            else:
                log.debug(
                    f"File is neither a video nor a subtitle, will not be imported: {file}"
                )

    log.info(
        f"Found {len(all_files)} files ({len(video_files)} video files, {len(subtitle_files)} subtitle files) for further processing."
    )
    return video_files, subtitle_files, all_files
