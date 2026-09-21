# Project Structure

```
collection-onepiece/
├── backend/
│   ├── main.go              # Entry point: DB connect + configuração Gin + registro de rotas
│   ├── go.mod / go.sum
│   ├── Dockerfile           # Multi-stage build
│   ├── database/
│   │   └── database.go      # Singleton database.DB, Connect(), AutoMigrate, getEnv()
│   ├── handlers/
│   │   ├── volume.go        # GetVolumes, UpdateVolume, SyncVolumes, syncFromJikan
│   │   └── volume_test.go   # Testes de integração com SQLite in-memory
│   └── models/
│       └── volume.go        # Struct Volume com tags GORM e JSON
├── frontend/
│   ├── index.html           # Shell HTML — carrega utils.js e app.js
│   ├── style.css            # Dark theme com CSS custom properties
│   ├── utils.js             # Funções puras de formatação de data (testáveis via Node)
│   ├── utils.test.js        # Testes com node:test nativo
│   ├── app.js               # Estado global, fetch, renderização DOM, timeline
│   ├── nginx.conf           # Gzip + cache de assets estáticos
│   └── Dockerfile           # nginx:1.25-alpine
└── docker-compose.yml       # Serviços: db + backend + frontend + volume pgdata
```

## Convenções de organização

### Backend (camadas)
- `models/` — structs de domínio puras, sem lógica
- `database/` — conexão e migrations; expõe variável global `database.DB`
- `handlers/` — lógica de negócio + HTTP handlers; acessa `database.DB` diretamente
- `main.go` — wiring apenas (sem lógica de domínio)
- Novos recursos seguem o padrão: struct em `models/`, handler em `handlers/`, rota em `main.go`

### Frontend (separação de responsabilidades)
- `utils.js` — apenas funções puras exportáveis (sem acesso ao DOM)
- `app.js` — todo estado global e lógica de UI; estado declarado no topo do arquivo

## Convenções de código

### Go
- Comentários em **português brasileiro**
- Nomes exportados em PascalCase (`GetVolumes`, `UpdateVolume`)
- Structs internas de request em camelCase (`updateVolumeRequest`)
- Tags GORM e JSON explícitas em todos os campos do model
- Usar `*bool` / `*string` em request structs para distinguir zero-value de campo ausente
- Erros HTTP retornados como `gin.H{"error": "mensagem em português"}`
- Testes usam helper `setupTestDB(t)` com `t.Helper()` que substitui `database.DB` por SQLite in-memory

### JavaScript
- Comentários de seção com `// ─── Nome da Seção ────────`
- Funções declaradas com `function nome()` (não arrow functions no topo nível)
- `async/await` para todas as chamadas `fetch`
- Estado global declarado e comentado no topo de `app.js`

### CSS
- CSS Custom Properties no `:root` para todo o tema
- Classes em kebab-case (`.volume-card`, `.timeline-btn`)
- Utilitário `.hidden` com `display: none !important`
