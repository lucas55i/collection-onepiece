# ☠️ Collection One Piece

Aplicação para gerenciar sua coleção de volumes do mangá One Piece.  
Os volumes são importados da [Jikan API](https://jikan.moe/) e armazenados em PostgreSQL.

## Stack

| Camada    | Tecnologia              |
|-----------|-------------------------|
| Backend   | Go + Gin + GORM         |
| Banco     | PostgreSQL 16           |
| Frontend  | HTML + CSS + JS vanilla |
| Container | Docker + Docker Compose |
| Futuro    | Helm + ArgoCD (k8s)     |

## Estrutura

```
collection-onepiece/
├── backend/
│   ├── main.go
│   ├── handlers/
│   │   ├── volume.go
│   │   └── volume_test.go
│   ├── models/
│   │   └── volume.go
│   ├── database/
│   │   └── database.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── index.html
│   ├── style.css
│   ├── app.js
│   ├── nginx.conf
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

## Como rodar

### Pré-requisitos

- Docker
- Docker Compose

### Subir o ambiente

```bash
docker compose up --build
```

| Serviço  | URL                          |
|----------|------------------------------|
| Frontend | http://localhost:3000        |
| Backend  | http://localhost:8080        |
| Postgres | localhost:5432               |

### Sincronizar volumes

Acesse o frontend e clique em **"🔄 Sincronizar volumes"**, ou via API:

```bash
curl -X POST http://localhost:8080/api/sync
```

### Endpoints da API

| Método | Rota                | Descrição                          |
|--------|---------------------|------------------------------------|
| GET    | /api/volumes        | Lista todos os volumes             |
| PATCH  | /api/volumes/:id    | Atualiza campo `collected`         |
| POST   | /api/sync           | Importa volumes da Jikan API       |
| GET    | /health             | Health check                       |

### Exemplo PATCH

```bash
curl -X PATCH http://localhost:8080/api/volumes/1 \
  -H "Content-Type: application/json" \
  -d '{"collected": true}'
```

## Testes

```bash
cd backend
go test ./...
```

## Roadmap

- [ ] Deploy no Kubernetes (k8s)
- [ ] Helm Chart
- [ ] Pipeline CI/CD com ArgoCD
- [ ] Imagem de capa por volume via Jikan API
