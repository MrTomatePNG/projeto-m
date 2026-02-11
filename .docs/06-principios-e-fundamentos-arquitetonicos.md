# 6. Princípios e Fundamentos Arquitetônicos

As decisões de design e arquitetura do MemeDroid são embasadas em princípios e padrões amplamente aceitos na engenharia de software, bem como em um profundo conhecimento das ferramentas escolhidas. Este documento serve como um guia de estudo e referência para esses conceitos.

## 6.1 Design de Sistemas (Teoria)

*   **[The Twelve-Factor App](https://12factor.net/pt_br/)**
    *   **O que é?** Uma metodologia para construir aplicações SaaS robustas e escaláveis.
    *   **Como aplicamos?**
        *   **III. Configurações:** Todas as configurações (senhas de banco, hosts, etc.) são injetadas através de variáveis de ambiente, nunca "hard-coded".
        *   **VI. Processos:** A aplicação é executada como um ou mais processos *stateless*, o que é fundamental para a escalabilidade horizontal.
        *   **VII. Backing Services:** Serviços como PostgreSQL e Redis são tratados como recursos conectados via rede, cujas configurações podem ser trocadas facilmente.

*   **[System Design Primer](https://github.com/donnemartin/system-design-primer)**
    *   **O que é?** Uma coleção de conceitos e exemplos de arquitetura para sistemas em larga escala.
    *   **Conceitos chave utilizados:**
        *   **Cache-Aside:** O padrão que usamos com o Redis para aliviar a carga do banco de dados.
        *   **Sharding:** O conceito por trás da nossa estratégia de armazenamento de arquivos em disco e um possível caminho para escalar o banco de dados no futuro.
        *   **Load Balancing:** Essencial para distribuir o tráfego entre múltiplas instâncias da nossa API.

## 6.2 Concorrência em Go (Prática)

A "sala de máquinas" do MemeDroid depende de um uso correto e idiomático da concorrência em Go.

*   **[Go Concurrency Patterns: Pipelines and Cancellation](https://go.dev/blog/pipelines)**
    *   **Por que ler?** É a "bíblia" para o nosso projeto. Explica exatamente como montar o pipeline de processamento (fan-out/fan-in) e, crucialmente, como usar o `context.Context` para cancelar operações de longa duração (como um processamento de vídeo) de forma segura.

*   **[Visualizing Concurrency in Go](https://divan.dev/posts/go_concurrency_visualize/)**
    *   **Por que ler?** Se você tem dificuldade em mentalizar como goroutines e canais interagem, este artigo mostra animações que clarificam os conceitos de `fan-out`, `fan-in` e `backpressure`.

## 6.3 Performance de Banco de Dados

*   **[Use The Index, Luke!](https://use-the-index-luke.com/)**
    *   **O que é?** O guia definitivo sobre indexação em bancos de dados relacionais.
    *   **Aplicação:** Essencial para entender por que criamos certos índices e como otimizar queries lentas. Por exemplo, a criação de índices parciais para posts com status `completed` ou índices compostos para ordenação e filtragem.

*   **[PostgreSQL JSON Functions](https://devhints.io/postgresql-json)**
    *   **O que é?** Um "cheat sheet" para as poderosas funções JSON do PostgreSQL.
    *   **Aplicação:** Referência rápida para entender o poder do `json_agg` e `json_build_object`, que são usados para resolver o problema de N+1 queries.

## 6.4 Ferramentas Específicas

*   **[sqlc](https://sqlc.dev/)**
    *   **O que é?** Ferramenta que gera código Go tipado a partir de SQL.
    *   **Por que usar?** Para combinar a performance do SQL nativo com a segurança de tipos do Go, evitando ORMs pesados. O [Playground do sqlc](https://play.sqlc.dev/) é útil para testar queries antes de integrá-las ao projeto.

*   **[FFmpeg](https://ffmpeg.org/documentation.html)**
    *   **O que é?** A biblioteca "canivete suíço" para manipulação de áudio e vídeo.
    *   **Aplicação:** Usada no coração dos nossos workers para transcodificar vídeos para formatos web-friendly (como `mp4` com codec H.264) и gerar thumbnails.
    *   *Exemplo útil:* `ffmpeg -i input.mov -vcodec h264 -acodec aac output.mp4`

*   **[Redis](https://redis.io/docs/latest/develop/get-started/)**
    *   **O que é?** Um armazenamento de estrutura de dados em memória, usado como banco de dados, cache e message broker.
    *   **Aplicação:** Usamos para caching, com comandos como `SET`, `GET` e `EXPIRE` para garantir que nosso cache não fique obsoleto.
