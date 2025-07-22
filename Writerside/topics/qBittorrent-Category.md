# qBittorrent Category

qBittorrent supports saving Torrents to subdirectories based on the category of the Torrent.
The default category name that MediaManager uses is `MediaManager`.

With the variable `torrents.qbittorrent.category_name` you can change the category name that MediaManager uses when
adding Torrents to qBittorrent.

With the variable `torrents.qbittorrent.category_save_path` you can change the path where the Torrents are saved to. By
default, no subdirectory is used. Note that qBittorrent saves torrents to this path, so it must be a
valid path that qBittorrent can write to. Example value: `/data/torrents/MediaManager`. Note that for MediaManager to be
able to successfully import torrents, you must add the subdirectory to the `misc.torrent_directory` variable.
