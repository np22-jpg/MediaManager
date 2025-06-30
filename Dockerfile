FROM ghcr.io/astral-sh/uv:debian-slim
ARG VERSION
LABEL version=${VERSION}
LABEL description="Docker image for the backend of MediaManager"
ENV IMAGE_DIRECTORY=/data/images \
    TV_SHOW_DIRECTORY=/data/tv \
    MOVIE_DIRECTORY=/data/movies \
    TORRENT_DIRECTORY=/data/torrents \
    OPENID_ENABLED=FALSE \
    PUBLIC_VERSION=${VERSION}
WORKDIR /app

RUN apt-get update && \
    apt-get install -y ca-certificates gcc mime-support

COPY pyproject.toml uv.lock ./

RUN uv sync --locked

COPY --chmod=755 mediamanager-backend-startup.sh .
COPY media_manager ./media_manager
COPY alembic ./alembic
COPY alembic.ini .

EXPOSE 8000
CMD ["/app/mediamanager-backend-startup.sh"]