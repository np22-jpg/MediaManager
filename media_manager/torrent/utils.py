import logging
import mimetypes
from pathlib import Path

import patoolib

from config import BasicConfig
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
