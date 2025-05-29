import mimetypes

import requests


def get_year_from_first_air_date(first_air_date: str | None) -> int | None:
    if first_air_date:
        return int(first_air_date.split("-")[0])
    else:
        return None


def download_poster_image(storage_path=None, poster_url=None, show=None) -> bool:
    res = requests.get(poster_url, stream=True)
    content_type = res.headers["content-type"]
    file_extension = mimetypes.guess_extension(content_type)
    if res.status_code == 200:
        with open(storage_path.joinpath(str(show.id) + file_extension), "wb") as f:
            f.write(res.content)
        return True
    else:
        return False
