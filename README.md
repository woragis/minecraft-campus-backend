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
