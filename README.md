# CampusWorld — Backend

API Go do **CampusWorld** — fonte da verdade para jogadores, convites, trust, guildas e território.

## Arquitetura

Segue o padrão **handler → service → repository** do projeto [Lingo](../../Lokra/lingo/backend), documentado em:

- [`minecraft/ARCHITECTURE.md`](../../ARCHITECTURE.md) — arquitetura de referência + plano CampusWorld
- [`faculdade/CAMPUSWORLD.md`](../CAMPUSWORLD.md) — especificação do produto

## Stack

- Go 1.22+
- `net/http` (stdlib, Go 1.22+ routing)
- GORM + PostgreSQL 16
- Migrações SQL em `migrations/`
- Docker Compose (Postgres + API)

## Layout (planejado)

```text
backend/
├── docs/
├── migrations/
├── server/
│   ├── cmd/server/main.go
│   └── internal/
│       ├── httpserver/
│       ├── player/
│       ├── invite/
│       └── ...
├── docker-compose.yml
└── .env.example
```

## Fase 1 — escopo

- Whitelist (plugin consulta antes do login)
- Registro/upsert de jogadores
- Convites (`/invite` in-game → API)
- Perfil e árvore de convites (read-only web)
- Bootstrap automático do fundador (`BOOTSTRAP_MINECRAFT_UUID`)

## Fase 2 — confiança e guildas

- Trust Score v1 + eventos (`trust_events`)
- Sponsor Score (média dos convidados + penalidade por ban)
- Graduação automática de probation
- Guildas (criar, entrar, sair, listar)
- Perfil enriquecido com guilda
- Grafo de moderação (`/v1/players/{id}/sponsor-tree`)

## Fase 3 — território

- Cidades (vinculáveis a guildas)
- Claims com proteção de blocos (plugin consulta API)
- Alianças entre guildas
- Progressão territorial por tempo de conta + trust score
- Zonas: `urban`, `rural`, `industrial`, `historic`

## Fase 4 — auditoria e rollback

- `audit_events` append-only (blocos em claims)
- Ingest em batch via plugin (`POST /v1/internal/audit/events`)
- Rollback por jogador + janela temporal (até 24h)
- Itens de rollback: `restore` (desfaz break) / `remove` (desfaz place)
- Trust penalty `rollback_applied` (-5) no alvo

## Fase 5 — operações e visibilidade

- Profiles: `dev` | `budget` (default) | `production` — ver [BUDGET-DEPLOY](../docs/BUDGET-DEPLOY.md)
- Worker (`cmd/worker`) — backups, alertas, purge de audit (opt-in via `WORKER_ENABLED`)
- Backup Postgres local (`BACKUP_STORAGE=local`) — sem AWS obrigatório
- `GET /v1/metrics/overview`, `/v1/metrics/territory`
- `GET /v1/alerts`, `POST /v1/alerts/{id}/acknowledge`
- Alertas de griefing (threshold configurável)

```bash
# Stack barata (sem worker)
docker compose up -d

# Com worker + backup local
WORKER_ENABLED=1 BACKUP_ENABLED=1 BACKUP_STORAGE=local \
  docker compose --profile worker up -d
```

## Desenvolvimento

```bash
cp .env.example .env

# Postgres only (host port 15432)
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d postgres

# API local
cd server
export DATABASE_URL="postgres://campus:campus@127.0.0.1:15432/campusworld?sslmode=disable"
export PLUGIN_API_KEY="dev-plugin-key"
go run ./cmd/server
```

Ou stack completa:

```bash
docker compose up -d --build
curl http://127.0.0.1:8080/health
```

## Rotas Fase 1

| Método | Rota | Auth |
|--------|------|------|
| GET | `/health` | — |
| GET | `/ready` | — |
| GET | `/v1/internal/whitelist/{minecraftUuid}?username=` | `X-Plugin-Key` |
| POST | `/v1/internal/players/upsert` | `X-Plugin-Key` |
| POST | `/v1/internal/invites` | `X-Plugin-Key` |
| GET | `/v1/players/{id}` | — |
| GET | `/v1/players/minecraft/{minecraftUuid}` | — |
| GET | `/v1/players/{id}/invites` | — |
| GET | `/v1/invites/{code}` | — |
| GET | `/v1/guilds` | — |
| GET | `/v1/guilds/{id}` | — |
| GET | `/v1/guilds/slug/{slug}` | — |
| GET | `/v1/guilds/{id}/members` | — |
| GET | `/v1/players/{id}/trust-events` | — |
| GET | `/v1/players/{id}/sponsor-tree` | — |
| POST | `/v1/internal/guilds` | `X-Plugin-Key` |
| POST | `/v1/internal/guilds/{id}/join` | `X-Plugin-Key` |
| POST | `/v1/internal/guilds/{id}/leave` | `X-Plugin-Key` |
| POST | `/v1/internal/trust/events` | `X-Plugin-Key` |
| GET | `/v1/cities` | — |
| GET | `/v1/cities/{id}` | — |
| GET | `/v1/cities/slug/{slug}` | — |
| GET | `/v1/cities/{id}/claims` | — |
| GET | `/v1/claims/{id}` | — |
| GET | `/v1/guilds/{id}/alliances` | — |
| POST | `/v1/internal/cities` | `X-Plugin-Key` |
| POST | `/v1/internal/claims` | `X-Plugin-Key` |
| DELETE | `/v1/internal/claims/{id}` | `X-Plugin-Key` |
| GET | `/v1/internal/claims/permission` | `X-Plugin-Key` |
| POST | `/v1/internal/alliances` | `X-Plugin-Key` |
| GET | `/v1/players/{id}/audit-events` | — |
| GET | `/v1/rollbacks/{id}` | — |
| POST | `/v1/internal/audit/events` | `X-Plugin-Key` |
| POST | `/v1/internal/rollbacks` | `X-Plugin-Key` |
| GET | `/v1/internal/rollbacks/{id}/items` | `X-Plugin-Key` |
| POST | `/v1/internal/rollbacks/{id}/complete` | `X-Plugin-Key` |
| GET | `/v1/metrics/overview` | — |
| GET | `/v1/metrics/territory` | — |
| GET | `/v1/alerts` | — |
| POST | `/v1/alerts/{id}/acknowledge` | — |

## Bootstrap do fundador

```env
BOOTSTRAP_MINECRAFT_UUID=<uuid-minecraft>
BOOTSTRAP_USERNAME=Fundador
```

Cria um jogador `active` no primeiro boot se ainda não existir.

## Documentação

- [Arquitetura](../docs/ARCHITECTURE.md)
- [Setup Paper](../docs/SETUP-PAPER.md)
- [CampusWorld spec](../CAMPUSWORLD.md)
