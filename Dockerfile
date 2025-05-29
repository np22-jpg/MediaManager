FROM ghcr.io/astral-sh/uv:debian-slim
ARG VERSION
LABEL version=${VERSION}
LABEL description="Docker image for the backend of MediaManager"
WORKDIR /app
COPY media_manager ./media_manager
COPY alembic ./alembic
COPY alembic.ini .
COPY pyproject.toml .
COPY uv.lock .
RUN uv sync --locked
EXPOSE 8000
CMD ["uv", "run", "fastapi", "run", "/app/media_manager/main.py"]
