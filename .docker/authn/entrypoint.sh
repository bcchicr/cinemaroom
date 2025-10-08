#!/bin/sh

set -e

if [ -z "$APP_ENV" ]; then
    echo "Error: APP_ENV is not set."
    exit 1
fi

if [ "$APP_ENV" = "dev" ]; then
    echo "Starting server in dev mode"
    go run main.go
elif [ "$APP_ENV" = "test" ]; then
    echo "Building for test mode"
    go build -o app main.go
    echo "Executing server in test mode."
    ./app
else
    echo "Error: unsupported APP_ENV value: $APP_ENV"
    exit 1
fi
