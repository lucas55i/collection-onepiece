# Requirements Document

## Introduction

Este documento descreve os requisitos para duas melhorias no projeto **One Piece Collection** — uma aplicação web para rastrear a coleção pessoal de volumes do mangá One Piece.

A primeira melhoria consolida e formaliza o registro preciso de data e horário de aquisição de cada volume, garantindo que a informação seja persistida no backend, exibida no card e navegável por uma timeline de meses/dias. Embora a estrutura de dados (`acquired_at`) e parte da lógica de timeline já existam no código, os requisitos aqui estabelecidos definem o comportamento esperado de forma completa e verificável.

A segunda melhoria aplica uma nova identidade visual baseada na paleta **"Trufa de Chocolate"** — tons de marrom-escuro, caramelo e creme — substituindo a paleta atual monocromática escura, sem alterar a estrutura HTML ou a lógica de negócio.

---

## Glossary

- **Volume**: Um exemplar físico de volume do mangá One Piece, representado pelo model `Volume` no backend.
- **Backend**: Serviço Go/Gin que expõe a API REST na porta 8080.
- **Frontend**: Interface web estática (HTML/CSS/JS) servida pelo Nginx.
- **API**: Interface REST exposta pelo Backend em `/api`.
- **Jikan_API**: API pública `api.jikan.moe/v4` usada para sincronizar metadados dos volumes.
- **DB**: Banco de dados PostgreSQL gerenciado pelo GORM.
- **Collector**: O usuário que interage com a interface para registrar sua coleção.
- **Timeline**: Componente de navegação na interface que agrupa volumes adquiridos por mês e por dia.
- **Card**: Elemento visual da grade que representa um único Volume.
- **Paleta_Chocolate**: Conjunto de variáveis CSS definidas com base em tons de trufa de chocolate (marrom-escuro, caramelo, creme).

---

## Requirements

### Requirement 1: Registro da data e horário de aquisição

**User Story:** Como Collector, quero que a data e o horário exatos sejam registrados automaticamente quando marco um volume como adquirido, para que eu tenha um histórico preciso de quando completei cada parte da minha coleção.

#### Acceptance Criteria

1. WHEN o Collector envia `PATCH /api/volumes/:id` com `{"collected": true}` e o campo `acquired_at` não está presente no payload, THE Backend SHALL verificar se o Volume já possui `acquired_at` no DB: se o Volume não possuir `acquired_at`, o Backend SHALL persistir `acquired_at` com o timestamp UTC do momento exato da requisição com precisão mínima de segundos; se o Volume já possuir `acquired_at`, o Backend SHALL preservar o valor existente sem sobrescrevê-lo. Em ambos os casos, THE Backend SHALL retornar o Volume atualizado completo com HTTP 200, com o campo `acquired_at` serializado em formato ISO 8601 com timezone.

2. WHEN o Collector envia `PATCH /api/volumes/:id` com `{"collected": true, "acquired_at": "<valor>"}` e o valor de `acquired_at` é uma string ISO 8601 válida contendo timezone (`Z` ou offset `±HH:MM`), THE Backend SHALL persistir `acquired_at` no DB com o valor fornecido, sem sobrescrever com o timestamp do servidor, e SHALL retornar o Volume atualizado completo com HTTP 200.

3. WHEN o Collector envia `PATCH /api/volumes/:id` com `{"collected": false}`, THE Backend SHALL persistir `acquired_at` como `NULL` no DB e SHALL retornar o Volume atualizado completo com HTTP 200, incluindo `"acquired_at": null`.

4. IF o `id` fornecido em `PATCH /api/volumes/:id` não corresponder a nenhum Volume no DB, THEN THE Backend SHALL retornar HTTP 404 com corpo `{"error": "Volume não encontrado"}`.

5. IF o payload de `PATCH /api/volumes/:id` se enquadrar em qualquer uma das condições a seguir, THEN THE Backend SHALL retornar HTTP 400 com corpo contendo `"error"` e mensagem indicando a condição:
   - (a) JSON malformado ou não parseável;
   - (b) campo `collected` ausente do payload;
   - (c) campo `acquired_at` presente no payload com valor que não seja ISO 8601 válido com timezone obrigatório (`Z` ou offset `±HH:MM`);
   - (d) campo `collected` ou `acquired_at` presentes com tipo de dado incorreto (ex.: `collected` como string, `acquired_at` como número).

---

### Requirement 2: Exibição da data de aquisição no Card

**User Story:** Como Collector, quero ver a data em que adquiri cada volume exibida no seu Card, para que eu consiga identificar rapidamente quando cada volume foi adicionado à minha coleção.

#### Acceptance Criteria

1. WHEN o Frontend renderiza um Card de Volume cujo `acquired_at` não é `null`, THE Frontend SHALL exibir a data e o horário formatados em português brasileiro no padrão `DD/MM/AAAA às HH:MM` dentro do elemento `.volume-acquired` do Card.

2. WHEN o Frontend renderiza um Card de Volume cujo `acquired_at` é `null` ou ausente, THE Frontend SHALL exibir o elemento `.volume-acquired` sem texto, mantendo altura mínima igual ao line-height definido para o elemento, evitando colapso do espaço reservado e deslocamento de layout.

3. WHEN o Collector altera o estado `collected` de um Volume via checkbox e a API retorna resposta com sucesso, THE Frontend SHALL re-renderizar o Card com a data atualizada retornada pela API em no máximo 500ms após o recebimento da resposta, sem recarregar a página inteira. IF a API retornar erro, THEN THE Frontend SHALL manter o estado anterior do Card e exibir uma indicação de falha visível ao Collector.

4. IF um Volume está marcado como coletado E o campo `acquired_at` do Volume não é `null`, THEN THE Frontend SHALL aplicar `color: var(--collected-accent)` ao texto do elemento `.volume-acquired` do Card correspondente.

---

### Requirement 3: Timeline de navegação por data de aquisição

**User Story:** Como Collector, quero navegar pelos volumes que adquiri organizados por mês e por dia em uma timeline, para que eu possa visualizar meu progresso de coleção ao longo do tempo.

#### Acceptance Criteria

1. WHEN o Frontend carrega os volumes e pelo menos um Volume possui `acquired_at` não nulo, THE Frontend SHALL exibir o componente Timeline com botões agrupados por mês, ordenados cronologicamente do mais antigo ao mais recente.

2. WHEN nenhum Volume possui `acquired_at` não nulo, THE Frontend SHALL ocultar completamente o componente Timeline.

3. WHEN o Collector clica em um botão de mês na Timeline, THE Frontend SHALL filtrar a grade de Cards exibindo apenas os Volumes cujo `acquired_at` pertence ao mês selecionado, e SHALL exibir os botões de dias disponíveis naquele mês.

4. WHEN o Collector clica em um botão de dia na Timeline, THE Frontend SHALL filtrar a grade de Cards exibindo apenas os Volumes cujo `acquired_at` pertence ao dia selecionado dentro do mês ativo.

5. WHEN o Collector clica em um botão de mês já ativo, THE Frontend SHALL limpar o filtro e exibir todos os Volumes.

6. WHEN o Collector clica em um botão de dia já ativo, THE Frontend SHALL retornar ao filtro de mês, exibindo todos os Volumes do mês ativo.

7. WHEN o Collector clica no botão "Limpar filtro", THE Frontend SHALL remover todos os filtros ativos e exibir todos os Volumes.

8. THE Frontend SHALL exibir em cada botão de mês e de dia um badge numérico indicando a quantidade total absoluta de Volumes adquiridos naquele período, calculada sempre sobre o conjunto completo de Volumes — não sobre o subconjunto filtrado no momento.

9. WHEN o Collector altera o estado `collected` de um Volume, THE Frontend SHALL reconstruir a Timeline para refletir o novo estado de `acquired_at`. IF o Volume desmarcado possuía `acquired_at` não nulo E pertencia ao mês ou dia atualmente filtrado na Timeline, THEN THE Frontend SHALL manter o filtro ativo e remover o Volume da grade de Cards. IF o mês ou dia filtrado ficar sem Volumes após a remoção, THEN THE Frontend SHALL limpar automaticamente o filtro ativo e exibir todos os Volumes.

10. THE Frontend SHALL aplicar `background: var(--pill-active)` e `color: var(--pill-active-text)` ao botão de mês que estiver ativo. O botão de dia ativo SHALL receber o mesmo tratamento visual.

---

### Requirement 4: Redesign visual — Paleta "Trufa de Chocolate"

**User Story:** Como Collector, quero que a interface tenha uma identidade visual baseada em tons de trufa de chocolate (marrom-escuro, caramelo e creme), para que a experiência de navegação seja mais acolhedora e elegante.

#### Acceptance Criteria

1. THE Frontend SHALL aplicar uma paleta de cores de tons de marrom-escuro, caramelo e creme via variáveis CSS definidas na regra `:root`, com as seguintes variáveis substituindo a paleta anterior:

   | Variável                | Valor sugerido  | Semântica                              |
   |-------------------------|-----------------|----------------------------------------|
   | `--bg`                  | `#1a0f0a`       | Fundo principal (marrom muito escuro)  |
   | `--surface`             | `#2c1a10`       | Superfície de cards e header           |
   | `--border`              | `#3d2516`       | Bordas padrão                          |
   | `--accent`              | `#c8864a`       | Caramelo — cor de destaque principal   |
   | `--accent-hover`        | `#e09a5a`       | Caramelo claro — hover                 |
   | `--text`                | `#f5e6d3`       | Creme — texto principal                |
   | `--text-muted`          | `#9e7a5e`       | Tom médio — texto secundário           |
   | `--collected`           | `#2a1f0e`       | Fundo do card coletado                 |
   | `--collected-border`    | `#8b6914`       | Borda dourada do card coletado         |
   | `--collected-accent`    | `#a8854a`       | Cor do texto de data no card coletado  |
   | `--timeline-bg`         | `#221309`       | Fundo da barra de timeline             |
   | `--timeline-border`     | `#3a2010`       | Borda da timeline                      |
   | `--pill-bg`             | `#2e1a0e`       | Fundo dos botões pill da timeline      |
   | `--pill-active`         | `var(--accent)` | Fundo do pill ativo                    |
   | `--pill-active-text`    | `#0d0601`       | Texto no pill ativo                    |

2. WHEN o Collector passa o cursor sobre um Card não coletado, THE Frontend SHALL aplicar `border-color: var(--accent)` ao Card. WHEN o cursor deixar o Card, THE Frontend SHALL restaurar `border-color` ao valor padrão `var(--border)`.

3. WHEN o Collector passa o cursor sobre um botão da Timeline, THE Frontend SHALL aplicar `border-color: var(--accent)` e `color: var(--accent)` ao botão. WHEN o cursor deixar o botão, THE Frontend SHALL restaurar `border-color` e `color` aos seus valores padrão.

4. THE Frontend SHALL manter inalterados os seguintes atributos visuais e estruturais: família tipográfica, tamanhos de fonte, valores de padding e margin, gap da grade de Cards, border-radius dos elementos, estrutura HTML (tags, classes e hierarquia) e número de colunas da grade.

5. THE Frontend SHALL garantir contraste mínimo de 4.5:1 entre `--text` e `--bg` (conformidade WCAG AA para texto normal).

---

### Requirement 5: Sincronização de volumes (comportamento existente preservado)

**User Story:** Como Collector, quero que a sincronização com a Jikan API continue funcionando após as mudanças, para que eu possa importar novos volumes sem perder dados de aquisição já registrados.

#### Acceptance Criteria

1. WHEN o Collector aciona `POST /api/sync`, THE Backend SHALL criar registros de Volume no DB exclusivamente para números de volume ainda inexistentes no DB, sem alterar `collected` nem `acquired_at` dos Volumes já existentes.

2. WHEN `POST /api/sync` for acionado e todos os volumes já existirem no DB, THE Backend SHALL completar a operação sem modificar nenhum registro existente e SHALL retornar HTTP 200.

3. WHEN o Backend executa a sincronização e a Jikan_API retorna `volumes: 0`, THE Backend SHALL usar o valor de fallback `114` como total de volumes a sincronizar.

4. IF a Jikan_API retornar erro ou não responder dentro de 10 segundos, THEN THE Backend SHALL retornar HTTP 500 com mensagem de erro indicando falha na comunicação com a API externa, sem modificar nenhum Volume já existente no DB.
