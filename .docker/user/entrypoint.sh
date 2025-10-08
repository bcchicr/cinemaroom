#!/bin/sh

set -e

echo "Starting server with config: ${RR_CONFIG_FILE}"
rr serve -c "${RR_CONFIG_FILE}" --debug