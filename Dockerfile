FROM ghcr.io/astral-sh/uv:debian-slim
ARG VERSION
LABEL version=${VERSION}
LABEL description="Docker image for the backend of MediaManager"
ENV MISC__IMAGE_DIRECTORY=/data/images \
    MISC__TV_DIRECTORY=/data/tv \
    MISC__MOVIE_DIRECTORY=/data/movies \
    MISC__TORRENT_DIRECTORY=/data/torrents \
    PUBLIC_VERSION=${VERSION} \
    MISC__API_BASE_PATH="/api/v1" \
    CONFIG_FILE="/app/config.toml"

WORKDIR /app

RUN apt-get update && \
    apt-get install -y ca-certificates gcc mime-support curl gzip unzip tar 7zip bzip2 unar

COPY pyproject.toml uv.lock ./

RUN uv sync --locked

COPY --chmod=755 mediamanager-backend-startup.sh .
COPY media_manager ./media_manager
COPY alembic ./alembic
COPY alembic.ini .

HEALTHCHECK CMD curl -f http://localhost:8000${MISC__API_BASE_PATH}/health || exit 1
EXPOSE 8000
CMD ["/app/mediamanager-backend-startup.sh"]
