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

# Initialize config if it doesn't exist
CONFIG_DIR=${CONFIG_DIR:-/app/config}
CONFIG_FILE="$CONFIG_DIR/config.toml"
EXAMPLE_CONFIG="/app/config.example.toml"

echo "Checking configuration setup..."

# Create config directory if it doesn't exist
if [ ! -d "$CONFIG_DIR" ]; then
    echo "Creating config directory: $CONFIG_DIR"
    mkdir -p "$CONFIG_DIR"
fi

# Copy example config if config.toml doesn't exist
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found. Copying example config to: $CONFIG_FILE"
    if [ -f "$EXAMPLE_CONFIG" ]; then
        cp "$EXAMPLE_CONFIG" "$CONFIG_FILE"
        echo "Example config copied successfully!"
        echo "Please edit $CONFIG_FILE to configure MediaManager for your environment."
        echo "Important: Make sure to change the token_secret value!"
    else
        echo "ERROR: Example config file not found at $EXAMPLE_CONFIG"
        exit 1
    fi
else
    echo "Config file found at: $CONFIG_FILE"
fi

echo "Running DB migrations..."

uv run alembic upgrade head

echo "Starting MediaManager backend service..."
echo ""
echo "🔐 LOGIN INFORMATION:"
echo "   If this is a fresh installation, a default admin user will be created automatically."
echo "   Check the application logs for the login credentials."
echo "   You can also register a new user and it will become admin if the email"
echo "   matches one of the admin_emails in your config.toml"
echo ""
uv run fastapi run /app/media_manager/main.py --port 8000
