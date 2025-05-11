def get_year_from_first_air_date(first_air_date: str | None) -> int | None:
    if first_air_date:
        return int(first_air_date.split('-')[0])
    else:
        return None
