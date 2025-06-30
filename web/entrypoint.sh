#!/bin/sh
# This script is used to start the MediaManager frontend service.


# text created with https://patorjk.com/software/taag/ font: Slanted
cat << EOF
    __  ___         ___       __  ___                                     ______                 __                 __
   /  |/  /__  ____/ (_)___ _/  |/  /___ _____  ____ _____ ____  _____   / ____/________  ____  / /____  ____  ____/ /
  / /|_/ / _ \/ __  / / __ \`/ /|_/ / __ \`/ __ \/ __ \`/ __ \`/ _ \/ ___/  / /_  / ___/ __ \/ __ \/ __/ _ \/ __ \/ __  /
 / /  / /  __/ /_/ / / /_/ / /  / / /_/ / / / / /_/ / /_/ /  __/ /     / __/ / /  / /_/ / / / / /_/  __/ / / / /_/ /
/_/  /_/\___/\__,_/_/\__,_/_/  /_/\__,_/_/ /_/\__,_/\__, /\___/_/     /_/   /_/   \____/_/ /_/\__/\___/_/ /_/\__,_/
                                                   /____/
EOF
echo "Buy me a coffee at https://buymeacoffee.com/maxdorninger"
echo "Starting MediaManager frontend service..."
node build/index.js
