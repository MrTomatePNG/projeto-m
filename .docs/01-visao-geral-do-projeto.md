# 1. Visão Geral do Projeto "MemeDroid (Webflix)"

**Stack Principal:** Go (Golang), PostgreSQL, Redis, Docker.

## 1.1 Objetivo Primário

O MemeDroid é projetado como uma plataforma de mídia social otimizada para o consumo intensivo de conteúdo visual (memes, vídeos). A arquitetura foi concebida para suportar uma carga de trabalho predominantemente **Read-Heavy** (muito mais leituras do que escritas).

O objetivo principal é construir uma plataforma **performática e resiliente**, capaz de:
*   Gerenciar o ciclo de vida completo de uploads de mídia de forma assíncrona.
*   Otimizar mídias para consumo rápido na web (geração de `webp`/`mp4` e thumbnails).
*   Manter uma alta disponibilidade e responsividade da API, mesmo sob picos de carga.
*   Garantir a consistência dos dados e a recuperação automática de falhas de processamento.

## 1.2 O Desafio Central: Processamento de Mídia

O gargalo mais significativo em uma aplicação como esta é o processamento de arquivos de mídia (vídeo e imagem). São operações que consomem muitos recursos computacionais (CPU e RAM) e podem levar de segundos a minutos para serem concluídas.

Se esse processamento fosse feito de forma síncrona (durante a requisição HTTP), a experiência do usuário seria severamente degradada, com longos tempos de espera e possíveis timeouts. Para resolver isso, o MemeDroid desacopla completamente o upload inicial do processamento subsequente.

## 1.3 Princípios Arquitetônicos Chave

Estes princípios são a espinha dorsal do projeto e guiam todas as decisões de design e implementação:

*   **Assincronicidade:** Operações de longa duração **não devem** bloquear a resposta ao usuário. O upload deve retornar uma confirmação rápida (`202 Accepted`), e o processamento real ocorre em segundo plano (background).
    *   **Por que?** Para garantir uma excelente experiência do usuário e manter a API responsiva e escalável.

*   **Concorrência Gerenciada:** A utilização de `Worker Pools` de forma controlada é fundamental para lidar com múltiplas tarefas de processamento simultaneamente sem esgotar os recursos do servidor.
    *   **Por que?** Para otimizar o uso do hardware, evitar sobrecarga e garantir estabilidade sob picos de demanda.

*   **Resiliência e Tolerância a Falhas:** O sistema deve ser capaz de se recuperar automaticamente de falhas inesperadas e lidar com "estados zumbis" (ex: posts travados em "processamento").
    *   **Por que?** Para garantir a integridade dos dados, evitar o acúmulo de "lixo" no sistema (arquivos órfãos) e manter a funcionalidade mesmo em face de erros. O processo `Janitor` é o principal ator aqui.

*   **Observabilidade:** A capacidade de monitorar, rastrear e entender o comportamento do sistema em tempo real e de forma retroativa é crucial.
    *   **Por que?** Para identificar e diagnosticar problemas rapidamente, entender fluxos de usuário e otimizar a performance. Logs estruturados e `Correlation IDs` são as ferramentas para isso.
