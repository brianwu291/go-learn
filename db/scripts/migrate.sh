#!/bin/sh

# cd `/db` folder
# chmod +x ./scripts/migrate.sh
# Usage examples:
# ./scripts/migrate.sh up
# ./scripts/migrate.sh down
# ./scripts/migrate.sh force VERSION
# ./scripts/migrate.sh goto VERSION

# database connection values
DB_HOST="${DB_HOST}"
DB_PORT="${DB_PORT}"
DB_USER="${DB_USER}"
DB_NAME="${DB_NAME}"

# connection string
DB_URL="postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# command is the first argument
COMMAND=$1

# print usage if no command provided
if [ -z "$COMMAND" ]; then
    echo "Usage: $0 <up|down|version>"
    exit 1
fi

# run migration
migrate -path migrations -database "${DB_URL}" "${COMMAND}"
