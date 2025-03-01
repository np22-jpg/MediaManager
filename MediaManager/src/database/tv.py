from database import PgDatabase, log

# TODO: add NOT NULL and default values to DB

def init_table():
    with PgDatabase() as db:
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_show (
                id UUID PRIMARY KEY,
                external_id TEXT,
                metadata_provider TEXT,
                name TEXT,
                episode_count INTEGER,
                season_count INTEGER,
                UNIQUE (external_id, metadata_provider)

            );""")
        log.info("tv_show Table initialized successfully")
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_season (
                show_id UUID  REFERENCES tv_show(id),
                season_number INTEGER,
                episode_count INTEGER,
                CONSTRAINT PK_season PRIMARY KEY (show_id,season_number)
            );""")
        log.info("tv_seasonTable initialized successfully")
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_episode (
                season_number  INTEGER,
                show_id uuid,
                episode_number INTEGER,
                title TEXT,
                CONSTRAINT PK_episode PRIMARY KEY (season_number,show_id,episode_number),
                FOREIGN KEY (season_number, show_id) REFERENCES tv_season(season_number,show_id)

            );""")
        log.info("tv_episode Table initialized successfully")
