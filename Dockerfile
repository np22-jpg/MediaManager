FROM ghcr.io/astral-sh/uv:debian-slim
ARG VERSION
LABEL version=${VERSION}
LABEL description="Docker image for the backend of MediaManager"

ENV IMAGE_DIRECTORY=/data/images
ENV TV_SHOW_DIRECTORY=/data/tv
ENV MOVIE_DIRECTORY=/data/movies
ENV TORRENT_DIRECTORY=/data/torrents
ENV OPENID_ENABLED=FALSE

RUN apt update && apt install -y ca-certificates gcc python3-dev

WORKDIR /app
COPY --chmod=755 mediamanager-backend-startup.sh .
COPY media_manager ./media_manager
COPY alembic ./alembic
COPY alembic.ini .
COPY pyproject.toml .
COPY uv.lock .
RUN uv sync --locked
EXPOSE 8000
CMD ["/app/mediamanager-backend-startup.sh"]