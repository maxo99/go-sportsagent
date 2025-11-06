#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

docker build -t maxo5499/sportsstack-go-sportsagent:latest .
