FROM ghcr.io/astral-sh/uv:debian-slim AS builder
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates gcc python3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY pyproject.toml uv.lock ./

RUN uv sync --locked

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

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.*/site-packages /usr/local/lib/python3.*/site-packages/

COPY --chmod=755 mediamanager-backend-startup.sh .
COPY media_manager ./media_manager
COPY alembic ./alembic
COPY alembic.ini .

EXPOSE 8000
CMD ["/app/mediamanager-backend-startup.sh"]