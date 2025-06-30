#!/bin/sh
# This script is used to start the MediaManager backend service.


# text created with https://patorjk.com/software/taag/ font: Slanted
echo "
    __  ___         ___       __  ___                                     ____             __                  __
   /  |/  /__  ____/ (_)___ _/  |/  /___ _____  ____ _____ ____  _____   / __ )____ ______/ /_____  ____  ____/ /
  / /|_/ / _ \/ __  / / __ \`/ /|_/ / __ \`/ __ \/ __ \`/ __ \`/ _ \/ ___/  / __  / __ \`/ ___/ //_/ _ \/ __ \/ __  /
 / /  / /  __/ /_/ / / /_/ / /  / / /_/ / / / / /_/ / /_/ /  __/ /     / /_/ / /_/ / /__/ ,< /  __/ / / / /_/ /
/_/  /_/\___/\__,_/_/\__,_/_/  /_/\__,_/_/ /_/\__,_/\__, /\___/_/     /_____/\__,_/\___/_/|_|\___/_/ /_/\__,_/
                                                   /____/
"
echo "Buy me a coffee at https://buymeacoffee.com/maxdorninger"
echo "Running DB migrations..."

uv run alembic upgrade head

echo "Starting MediaManager backend service..."
uv run fastapi run /app/media_manager/main.py --port 8000
