FROM node:24-alpine AS frontend-build
WORKDIR /frontend
ARG VERSION
ARG BASE_PATH=""

COPY web/package*.json ./
RUN npm ci --production && npm cache clean --force

COPY web/ ./
RUN env PUBLIC_VERSION=${VERSION} PUBLIC_API_URL=${BASE_PATH}/api/v1 BASE_PATH=${BASE_PATH}/web npm run build

FROM ghcr.io/astral-sh/uv:debian-slim
ARG VERSION
ARG BASE_PATH=""
LABEL author="github.com/maxdorninger"
LABEL version=${VERSION}
LABEL description="Docker image for MediaManager"

ENV MISC__IMAGE_DIRECTORY=/data/images \
    MISC__TV_DIRECTORY=/data/tv \
    MISC__MOVIE_DIRECTORY=/data/movies \
    MISC__TORRENT_DIRECTORY=/data/torrents \
    PUBLIC_VERSION=${VERSION} \
    CONFIG_FILE="/app/config.toml"\
    BASE_PATH=${BASE_PATH}\
    FRONTEND_FILES_DIR="/app/web/build"


WORKDIR /app

RUN apt-get update && \
    apt-get install -y ca-certificates gcc mime-support curl gzip unzip tar 7zip bzip2 unar && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY pyproject.toml uv.lock ./
RUN uv sync --locked

COPY --chmod=755 mediamanager-backend-startup.sh .
COPY media_manager ./media_manager
COPY alembic ./alembic
COPY alembic.ini .

COPY --from=frontend-build /frontend/build /app/web/build

HEALTHCHECK CMD curl -f http://localhost:8000${BASE_PATH}/api/v1/health || exit 1
EXPOSE 8000
CMD ["/app/mediamanager-backend-startup.sh"]
