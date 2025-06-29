#!/bin/bash
# This script is used to start the MediaManager backend service.
uv run alembic upgrade head
uv run fastapi run /app/media_manager/main.py