# Product — Collection One Piece

Aplicação web para gerenciar uma coleção pessoal dos volumes do mangá **One Piece**.

## Funcionalidades principais

- **Catálogo de volumes** importado da [Jikan API](https://jikan.moe/) (wrapper da MyAnimeList, manga ID 13)
- **Marcar volumes** como coletados/não coletados via toggle
- **Rastreamento de aquisição**: `acquired_at` é preenchido automaticamente ao marcar como coletado
- **Timeline de aquisição**: navegação filtrada por mês e dia
- **Progresso da coleção**: exibição do total coletado vs total disponível

## Roadmap

- Deploy em Kubernetes (EKS / k3s / kind)
- Helm Chart
- CI/CD com ArgoCD
- Imagens de capa individuais por volume via Jikan API
