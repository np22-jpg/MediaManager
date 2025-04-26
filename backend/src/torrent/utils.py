import logging
from pathlib import Path

from config import BasicConfig
from torrent.schemas import Torrent


def list_files_recursively(path: Path = Path(".")) -> list[Path]:
    files = list(path.glob("**/*"))
    logging.debug(f"Found {len(files)} entries via glob")
    valid_files = []
    for x in files:
        if x.is_dir():
            logging.debug(f"'{x}' is a directory")
        elif x.is_symlink():
            logging.debug(f"'{x}' is a symlink")
        else:
            valid_files.append(x)
    logging.debug(f"Returning {len(valid_files)} files after filtering")
    return valid_files


def get_torrent_filepath(torrent: Torrent):
    return BasicConfig().torrent_directory / torrent.title
