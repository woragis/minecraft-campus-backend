#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ ! -f .env ]]; then
  cp .env.example .env
  echo "Created .env from .env.example"
fi

docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build postgres api

echo ""
echo "API:      http://127.0.0.1:8080/health"
echo "Postgres: localhost:15432"
echo ""
echo "Logs: docker compose logs -f api"
