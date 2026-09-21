# Implementation Plan: manga-acquisition-date

## Overview

Implementação incremental das quatro frentes da feature:

1. **Backend Go** — refatorar o handler `UpdateVolume` para suportar binding tipado, validação de ISO 8601 com timezone obrigatório e preservação de `acquired_at` existente.
2. **Frontend CSS** — substituir a paleta monocromática escura pela paleta "Trufa de Chocolate".
3. **Frontend JS** — ajustar a lógica de limpeza automática de filtro ao esvaziar um período.
4. **Testes** — ampliar `volume_test.go` com novos casos de exemplo e adicionar testes de propriedade com `pgregory.net/rapid`; adicionar testes unitários e de propriedade para funções puras do `app.js`.

---

## Tasks

- [x] 1. Refatorar o handler `UpdateVolume` no backend
  - [x] 1.1 Criar struct `updateVolumeRequest` com campos ponteiro e função `parseISO8601WithTZ`
    - Criar o struct `updateVolumeRequest` com `Collected *bool` e `AcquiredAt *string` em `handlers/volume.go`
    - Implementar a função `parseISO8601WithTZ(s string) (*time.Time, error)` que tenta `time.RFC3339` e `time.RFC3339Nano`
    - _Requirements: 1.1, 1.2, 1.5_

  - [x] 1.2 Reescrever a lógica principal de `UpdateVolume`
    - Substituir o binding atual pelo novo struct `updateVolumeRequest`
    - Retornar HTTP 400 quando `collected` for ausente (ponteiro nil)
    - Retornar HTTP 400 quando `acquired_at` estiver presente mas falhar em `parseISO8601WithTZ`
    - Buscar o volume existente com `DB.First` antes de atualizar (retorna HTTP 404 se não encontrado)
    - Aplicar o fluxo de decisão do design: valor explícito → preserva existente → `now().UTC()`
    - Retornar o volume atualizado com HTTP 200
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

  - [ ]* 1.3 Adicionar casos de teste de exemplo em `volume_test.go`
    - `TestUpdateVolume_CollectedAbsent` — payload sem `collected` → HTTP 400
    - `TestUpdateVolume_AcquiredAtInvalidFormat` — `acquired_at` sem timezone → HTTP 400
    - `TestUpdateVolume_PreservesExistingAcquiredAt` — volume com `acquired_at` preenchido + PATCH sem o campo → mesmo valor retornado
    - _Requirements: 1.1, 1.5_

  - [ ]* 1.4 Escrever testes de propriedade para o handler (Property 1 a 5 e 12)
    - Adicionar `pgregory.net/rapid` ao `go.mod` com `go get pgregory.net/rapid`
    - **Property 1: Preservação de `acquired_at` ao coletar sem valor explícito** — Validates: Requirements 1.1
    - **Property 2: Atribuição de timestamp quando `acquired_at` é nulo** — Validates: Requirements 1.1
    - **Property 3: Round trip de `acquired_at` explícito** — Validates: Requirements 1.2
    - **Property 4: Limpeza de `acquired_at` ao desmarcar coleta** — Validates: Requirements 1.3
    - **Property 5: Rejeição de payloads inválidos** — Validates: Requirements 1.5
    - **Property 12: Idempotência da sincronização** — Validates: Requirements 5.1
    - Incluir comentário `// Feature: manga-acquisition-date, Property N: <texto>` em cada teste
    - _Requirements: 1.1, 1.2, 1.3, 1.5, 5.1_

- [x] 2. Checkpoint — verificar backend
  - Executar `go test ./handlers/...` e garantir que todos os testes passem; tirar dúvidas antes de continuar.

- [x] 3. Atualizar a paleta de cores no `style.css`
  - [x] 3.1 Substituir as variáveis CSS na regra `:root`
    - Alterar os valores de `--bg`, `--surface`, `--border`, `--accent`, `--accent-hover`, `--text`, `--text-muted`, `--collected`, `--collected-border` conforme a tabela do design
    - Adicionar a nova variável `--collected-accent: #a8854a`
    - Alterar `--timeline-bg`, `--timeline-border`, `--pill-bg`, `--pill-active-text` conforme o design
    - _Requirements: 4.1_

  - [x] 3.2 Atualizar a regra `.volume-card.collected .volume-acquired`
    - Substituir o valor hardcoded `#81c784` por `color: var(--collected-accent)`
    - _Requirements: 2.4, 4.1_

- [x] 4. Corrigir lógica de limpeza automática de filtro no `app.js`
  - [x] 4.1 Adicionar verificação de período vazio em `toggleCollected`
    - Após `buildTimeline()`, inserir o bloco que verifica se `allVolumes` ainda contém volumes no período ativo (`activeMonthKey` / `activeDayKey`)
    - Chamar `clearFilter()` e retornar antecipado quando o período ficar vazio
    - Caso contrário, chamar `applyFilter()` normalmente
    - _Requirements: 3.9_

- [x] 5. Checkpoint — verificar frontend
  - Garantir que o CSS e o JS não possuem erros de sintaxe; tirar dúvidas antes de continuar.

- [x] 6. Adicionar testes para as funções puras do `app.js`
  - [x] 6.1 Extrair funções puras para um módulo `frontend/utils.js` e criar `frontend/utils.test.js`
    - Mover `formatAcquiredAt`, `toMonthKey`, `toDateKey`, `formatMonthLabel`, `formatDayLabel` para `frontend/utils.js` usando `module.exports`
    - Criar `frontend/utils.test.js` usando o módulo `node:test` e `node:assert` (zero dependências externas)
    - Adicionar import de `utils.js` em `app.js` (manter compatibilidade com o browser via verificação `typeof module !== 'undefined'`)
    - _Requirements: 2.1, 2.2_

  - [ ]* 6.2 Escrever testes unitários para `toMonthKey` e `toDateKey`
    - Cobrir strings ISO 8601 com timezone `Z`, offset positivo e offset negativo
    - Verificar o formato de saída `"YYYY-MM"` e `"YYYY-MM-DD"` respectivamente
    - _Requirements: 3.1, 3.3, 3.4_

  - [ ]* 6.3 Escrever teste de propriedade para `formatAcquiredAt` (Property 6)
    - Gerar strings ISO 8601 válidas aleatórias e verificar que a saída contém padrão `DD/MM/AAAA` e `HH:MM` separados por `"às"`
    - Verificar que retorna `''` para `null` ou `undefined`
    - **Property 6: Formatação de data no card** — Validates: Requirements 2.1
    - _Requirements: 2.1, 2.2_

  - [ ]* 6.4 Escrever testes de propriedade para lógica de filtro e timeline (Properties 7–11)
    - Criar função pura `buildTimelineData(volumes)` em `utils.js` que retorna `{ byMonth }` para teste isolado
    - Criar função pura `applyFilter(volumes, activeMonthKey, activeDayKey)` em `utils.js`
    - Criar função pura `shouldClearFilter(volumes, activeMonthKey, activeDayKey)` em `utils.js`
    - **Property 7: Ordenação cronológica dos meses** — Validates: Requirements 3.1
    - **Property 8: Corretude do filtro por mês** — Validates: Requirements 3.3
    - **Property 9: Corretude do filtro por dia** — Validates: Requirements 3.4
    - **Property 10: Corretude dos badges de contagem** — Validates: Requirements 3.8
    - **Property 11: Limpeza automática de filtro ao esvaziar período** — Validates: Requirements 3.9
    - _Requirements: 3.1, 3.3, 3.4, 3.8, 3.9_

- [x] 7. Checkpoint final — testes completos
  - Executar `go test ./handlers/...` (backend) e `node --test frontend/utils.test.js` (frontend); garantir que todos os testes passam; tirar dúvidas antes de encerrar.

---

## Notes

- Tasks marcadas com `*` são opcionais e podem ser puladas para um MVP mais rápido
- `pgregory.net/rapid` ainda não está no `go.mod` — a task 1.4 inclui o `go get`
- O model `models.Volume` e o schema do banco **não precisam de alteração** — `AcquiredAt *time.Time` já existe
- O HTML (`index.html`) **não precisa de alteração** — estrutura e classes já estão corretas
- As funções de UI que dependem do DOM (Requirements 2.2, 2.4, 3.2, 3.5–3.7, 3.10, 4.x) são verificadas manualmente
- A task 6.1 garante que as funções puras sejam testáveis tanto no Node.js quanto no browser

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "3.1"] },
    { "id": 1, "tasks": ["1.2", "3.2", "4.1"] },
    { "id": 2, "tasks": ["1.3", "1.4", "6.1"] },
    { "id": 3, "tasks": ["6.2", "6.3", "6.4"] }
  ]
}
```
