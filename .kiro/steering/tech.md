# Tech Stack

## Backend (Go)

- **Linguagem:** Go 1.22
- **Módulo:** `github.com/collection-onepiece/backend`
- **Framework HTTP:** Gin (`github.com/gin-gonic/gin` v1.9.1) + CORS via `gin-contrib/cors`
- **ORM:** GORM v1.25 com driver `postgres` (pgx v5) em produção
- **Testes:** driver `sqlite` in-memory (`gorm.io/driver/sqlite`) — apenas nos testes

## Frontend (Vanilla)

- HTML5 + CSS3 + JavaScript puro — sem framework, sem bundler, sem npm
- Servidor em produção: **nginx 1.25-alpine**
- Testes: Node.js built-in test runner (`node:test` + `node:assert`)

## Infraestrutura

- **Banco de dados:** PostgreSQL 16-alpine
- **Containerização:** Docker + Docker Compose v3.9
- Backend: Dockerfile multi-stage (`golang:1.22-alpine` → `alpine:3.19`)
- Frontend: `nginx:1.25-alpine`

## API Backend

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/volumes` | Lista todos os volumes ordenados por `volume_number` |
| PATCH | `/api/volumes/:id` | Atualiza `collected` e `acquired_at` |
| POST | `/api/sync` | Importa volumes da Jikan API (idempotente) |
| GET | `/health` | Health check |

## Comandos comuns

```bash
# Subir todo o ambiente (db + backend + frontend)
docker compose up --build

# Rodar testes do backend (SQLite in-memory, sem PostgreSQL)
cd backend && go test ./...

# Rodar testes do frontend
node frontend/utils.test.js

# Sincronizar volumes via API
curl -X POST http://localhost:8080/api/sync
```

**URLs locais:** frontend → http://localhost:3000 | backend → http://localhost:8080 | PostgreSQL → localhost:5432
