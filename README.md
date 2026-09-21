# ☠️ Collection One Piece

Aplicação web para gerenciar uma coleção pessoal dos volumes do mangá **One Piece**.  
O catálogo é importado da [Jikan API](https://jikan.moe/) (wrapper da MyAnimeList) e armazenado localmente em PostgreSQL.

---

## Funcionalidades

- **Catálogo completo** de volumes importado via Jikan API (manga ID 13)
- **Toggle de coleção** — marcar/desmarcar volumes como coletados
- **Rastreamento de aquisição** — `acquired_at` preenchido automaticamente ao marcar como coletado e zerado ao desmarcar
- **Timeline na sidebar** — navegação hierárquica por mês → dia com contadores
- **Progresso da coleção** — total coletado vs. total disponível no header
- **Layout assimétrico** — sidebar de navegação + área principal com grid de volumes

---

## Stack

| Camada     | Tecnologia                                      |
|------------|-------------------------------------------------|
| Backend    | Go 1.22 · Gin v1.9.1 · GORM v1.25              |
| Banco      | PostgreSQL 16-alpine                            |
| Frontend   | HTML5 · CSS3 · JavaScript vanilla (sem npm)     |
| Servidor   | nginx 1.25-alpine                               |
| Container  | Docker · Docker Compose v3.9                    |
| Testes     | Go `testing` · Node.js `node:test` nativo       |

---

## Estrutura do projeto

```
collection-onepiece/
├── src/
│   ├── backend/
│   │   ├── main.go              # Entry point: DB + Gin + rotas
│   │   ├── go.mod / go.sum
│   │   ├── Dockerfile           # Multi-stage: golang:1.22-alpine → alpine:3.19
│   │   ├── database/
│   │   │   └── database.go      # Singleton DB, Connect(), AutoMigrate
│   │   ├── handlers/
│   │   │   ├── volume.go        # GetVolumes, UpdateVolume, SyncVolumes
│   │   │   └── volume_test.go   # Testes de integração com SQLite in-memory
│   │   └── models/
│   │       └── volume.go        # Struct Volume com tags GORM e JSON
│   └── frontend/
│       ├── index.html           # Shell HTML
│       ├── style.css            # Dark theme com CSS custom properties
│       ├── app.js               # Estado global, fetch, renderização, timeline
│       ├── utils.js             # Funções puras de formatação de data
│       ├── utils.test.js        # Testes com node:test nativo
│       ├── nginx.conf           # Gzip + cache de assets + SPA fallback
│       └── Dockerfile           # nginx:1.25-alpine
└── docker-compose.yml           # Serviços: db · backend · frontend
```

---

## Como rodar

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)

### Subir o ambiente completo

```bash
docker compose up --build
```

| Serviço    | URL                       |
|------------|---------------------------|
| Frontend   | http://localhost:3000     |
| Backend    | http://localhost:8080     |
| PostgreSQL | localhost:5432            |

### Sincronizar os volumes (primeira execução)

Via frontend — clique no botão **🔄 Sincronizar volumes**, ou via terminal:

```bash
curl -X POST http://localhost:8080/api/sync
```

---

## API

| Método  | Rota                 | Descrição                                              |
|---------|----------------------|--------------------------------------------------------|
| `GET`   | `/api/volumes`       | Lista todos os volumes ordenados por `volume_number`   |
| `PATCH` | `/api/volumes/:id`   | Atualiza `collected` e opcionalmente `acquired_at`     |
| `POST`  | `/api/sync`          | Importa volumes da Jikan API (idempotente)             |
| `GET`   | `/health`            | Health check                                           |

### PATCH `/api/volumes/:id`

**Campos do payload:**

| Campo        | Tipo      | Obrigatório | Descrição                                    |
|--------------|-----------|-------------|----------------------------------------------|
| `collected`  | `boolean` | ✅          | Marca ou desmarca o volume                   |
| `acquired_at`| `string`  | ❌          | ISO 8601 com timezone (ex: `2024-03-15T10:00:00-03:00`) |

**Lógica de `acquired_at`:**
- `collected: false` → campo zerado automaticamente
- `collected: true` + `acquired_at` no payload → usa o valor fornecido
- `collected: true` + volume já tinha data no banco → preserva a data existente
- `collected: true` + sem data alguma → usa `time.Now().UTC()` do servidor

**Exemplo:**

```bash
# Marcar como coletado
curl -X PATCH http://localhost:8080/api/volumes/1 \
  -H "Content-Type: application/json" \
  -d '{"collected": true}'

# Marcar com data de aquisição específica
curl -X PATCH http://localhost:8080/api/volumes/1 \
  -H "Content-Type: application/json" \
  -d '{"collected": true, "acquired_at": "2024-03-15T10:00:00-03:00"}'
```

---

## Variáveis de ambiente

O backend aceita as seguintes variáveis (todas com fallback padrão):

| Variável      | Padrão      | Descrição              |
|---------------|-------------|------------------------|
| `DB_HOST`     | `localhost` | Host do PostgreSQL     |
| `DB_PORT`     | `5432`      | Porta do PostgreSQL    |
| `DB_USER`     | `onepiece`  | Usuário do banco       |
| `DB_PASSWORD` | `onepiece`  | Senha do banco         |
| `DB_NAME`     | `onepiece`  | Nome do banco          |
| `PORT`        | `8080`      | Porta HTTP do backend  |

---

## Testes

```bash
# Backend — SQLite in-memory (não requer PostgreSQL)
cd src/backend && go test ./...

# Frontend — Node.js test runner nativo
node src/frontend/utils.test.js
```

---

## Roadmap

- [ ] Imagens de capa individuais por volume via Jikan API
- [ ] Deploy no Kubernetes (EKS / k3s / kind)
- [ ] Helm Chart
- [ ] CI/CD com ArgoCD
