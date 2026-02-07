Com certeza. Aqui está o **Documento de Arquitetura e Design Técnico (TDD)** do Projeto "MemeDroid".

Este documento consolida todas as decisões, teorias e padrões que discutimos. Ele serve como o seu "Manual de Bordo" para o desenvolvimento.

---

# 📘 MemeDroid Backend: Documento de Design Técnico

**Versão:** 1.0
**Stack Principal:** Go (Golang), PostgreSQL, Redis, Docker.
**Objetivo:** Plataforma de rede social focada em mídia (memes/vídeos) com alta performance de leitura e processamento assíncrono resiliente.

---

## 1. Visão Geral e Objetivos

O sistema é projetado para suportar uma carga de trabalho **Read-Heavy** (muito mais leitura do que escrita). O gargalo principal identificado é o processamento de mídia (vídeo/imagem), que deve ser desacoplado da resposta HTTP para não degradar a experiência do usuário.

### Princípios Chave:

* **Assincronicidade:** O upload devolve `202 Accepted` imediatamente; o processamento ocorre em background.
* **Concorrência Gerenciada:** Uso estrito de Worker Pools para evitar exaustão de recursos (CPU/RAM).
* **Resiliência:** O sistema deve ser capaz de se recuperar de falhas (via Janitor) e não deixar "lixo" (arquivos órfãos) no disco.
* **Observabilidade:** Logs estruturados para rastrear a vida útil de um request.

---

## 2. Arquitetura de Dados (Persistence Layer)

### 2.1. Tecnologias

* **PostgreSQL:** Fonte da verdade. Escolhido pela integridade referencial, tipos complexos (JSONB, Arrays) e suporte excelente via `sqlc`.
* **Redis:** Cache-Aside. Armazena timelines e contadores para aliviar o banco relacional nas leituras frequentes.

### 2.2. Modelagem Relacional

#### Entidades Principais:

1. **Users:**
* Foco em segurança (`password` hash via Bcrypt).
* Constraints `UNIQUE` em email/username para integridade.


2. **Posts:**
* **Status Enum:** `pending` -> `processing` -> `completed` (ou `failed`).
* **Media Hash:** SHA-256 do arquivo original para deduplicação (evitar reprocessar o mesmo meme viral).
* **Denormalização:** Uso de Agregação JSON (`json_agg`) na leitura para trazer Posts + Tags em uma única query, evitando o problema N+1.


3. **Tags (Sistema N:N):**
* Tabela `post_tags` atua como ponte.
* Índices em ambas as direções (`post_id` e `tag_id`) para buscas bidirecionais rápidas.



### 2.3. Query Pattern (sqlc)

Utilização de SQL puro gerando código Go tipado.

* *Vantagem:* Performance de SQL nativo com segurança de tipos do Go.
* *Estratégia:* Uso de parâmetros posicionais (`$1`, `$2`) e transações para operações críticas.

---

## 3. Modelo de Concorrência (Processing Layer)

A "Sala de Máquinas" do backend utiliza o padrão de **Pipeline** com estágios definidos.

### 3.1. O Fluxo (Pipeline)

1. **Ingestion (Handler):** Valida request, salva `raw file` no disco temporário, insere `status: pending` no DB. Envia `MediaJob` para o canal.
2. **Dispatcher (Channel):** Buffer (Fila) que segura os jobs até que um worker esteja livre. Atua como *Backpressure*.
3. **Fan-Out (Workers):**
* `N` Goroutines rodando em paralelo (`runtime.NumCPU()`).
* Cada worker possui `Context` com Timeout (ex: 5 min) para evitar processos zumbis.
* Executa tarefas pesadas (FFmpeg, Resize).


4. **Fan-In (Finalizer):**
* Uma única goroutine que coleta resultados.
* **Responsabilidade:** Atualizar DB (`completed`), deletar arquivo `raw`, invalidar Cache.
* Garante consistência e evita *Race Conditions* no banco.



### 3.2. O Sistema de Manutenção (Janitor)

Um processo independente (baseado em `time.Ticker`) que roda periodicamente.

* **Função:** Buscar posts travados em `processing` por tempo excedente (ex: > 30min).
* **Ação:** Marcar como `failed` e limpar arquivos temporários.
* **Objetivo:** Garantir a consistência eventual do sistema em caso de *Hard Crash* (queda de energia/servidor).

---

## 4. Organização do Código (Project Structure)

Estrutura baseada no "Standard Go Project Layout", focada em modularidade e encapsulamento.

```text
/memedroid
├── cmd/api/          # Entrypoint (main.go). Onde os Workers sobem.
├── internal/
│   ├── database/     # Código gerado pelo sqlc.
│   ├── workers/      # Lógica do Pipeline (Pool, Finalizer, Janitor).
│   ├── services/     # Regras de negócio (Auth, FileSystem, Hash).
│   └── handlers/     # Camada HTTP (decodifica JSON, chama services).
├── sql/              # Schemas e Queries SQL.
└── storage/          # Armazenamento local (com Sharding de pastas).

```

---

## 5. Observabilidade e Logs

* **Padrão:** Structured Logging (JSON).
* **Correlation ID:** Um ID único gerado no Request HTTP que viaja dentro do `MediaJob` e aparece em todos os logs (do worker ao finalizer).
* **Níveis:**
* `INFO`: Fluxo normal (Upload recebido, Job finalizado).
* `ERROR`: Falhas recuperáveis (Codec inválido).
* `FATAL`: Falhas de infra (Banco fora do ar).



---

## 6. Referências e Material de Estudo

### Design e Arquitetura

* **The Twelve-Factor App:** Metodologia para construção de apps SaaS (foco em Configuração e Processos).
* **System Design Primer:** Conceitos de Cache-aside, Sharding e Load Balancing.

### Go Concurrency

* **Go Blog - Pipelines and Cancellation:** A bíblia para entender como cancelar jobs no meio.
* **Go Memory Model:** Para entender por que Channels são preferíveis a Mutexes na maioria dos casos.

### Banco de Dados

* **Use The Index, Luke:** Guia definitivo sobre indexação e performance SQL.
* **PostgreSQL Documentation (JSON Types):** Para entender o poder do `json_agg`.

---

### ✅ Checklist para o MVP (Produto Mínimo Viável)

1. [ ] Configurar `docker-compose` (Go + Postgres + Redis).
2. [ ] Definir Schema SQL e rodar `sqlc generate` sem erros.
3. [ ] Implementar Cadastro/Login (JWT + Bcrypt).
4. [ ] Criar Worker Pool básico (Upload -> Print Log -> Delete).
5. [ ] Integrar FFmpeg/Imaging no Worker.
6. [ ] Implementar Finalizer para atualizar status no banco.
7. [ ] Implementar Janitor para limpeza.

---
Aqui está uma curadoria de materiais focada exatamente nas tecnologias e padrões que definimos para o **MemeDroid**.

Separei por categorias para facilitar seu estudo conforme você for implementando cada parte.

---

### 1. 🐹 Go Concurrency (O Motor do Sistema)

Esses são leituras obrigatórias para entender como fazer o seu **Worker Pool** e o **Pipeline** funcionarem sem vazamento de memória.

* **[Go Concurrency Patterns: Pipelines and Cancellation (Oficial)](https://go.dev/blog/pipelines)**
* *Por que ler:* É a bíblia para o seu projeto. Explica exatamente como usar `Context` para cancelar uploads e como montar o pipeline de processamento de imagem.


* **[Visualizing Concurrency in Go (Divan)](https://divan.dev/posts/go_concurrency_visualize/)**
* *Por que ler:* Se você gosta de ver as coisas funcionando, esse artigo mostra animações 3D de como as goroutines e channels interagem. Ajuda muito a mentalizar o Fan-out/Fan-in.


* **[Robust & Efficient Concurrency with Go (Video)](https://www.youtube.com/watch?v=5zXAHh5tJqQ)**
* *Por que ver:* Palestra excelente sobre como criar workers que não "engasgam" e como lidar com timeouts (essencial para o seu FFmpeg não rodar pra sempre).



---

### 2. 🗄️ Banco de Dados & SQL (A Persistência)

Como estamos usando **PostgreSQL** com **sqlc**, você precisa entender de performance e agregação.

* **[Use The Index, Luke!](https://use-the-index-luke.com/)**
* *Por que ler:* O guia definitivo sobre índices. Leia a seção sobre `WHERE` e `ORDER BY` para entender por que criamos aquele índice composto na coluna `created_at`.


* **[PostgreSQL JSON Functions (Cheat Sheet)](https://devhints.io/postgresql-json)**
* *Por que ler:* Resumo rápido de como usar `json_agg` e `json_build_object`. Vai te salvar quando você quiser mexer na query de buscar posts com tags.


* **[sqlc Playground](https://www.google.com/search?q=https://play.sqlc.dev/)**
* *Por que usar:* Antes de rodar no seu projeto e ter erro de sintaxe, teste suas queries aqui. Ele mostra na hora como o código Go vai ficar.



---

### 3. 🏗️ Arquitetura e Organização (A Estrutura)

Para manter suas pastas `internal`, `cmd` e `services` organizadas.

* **[Standard Go Project Layout](https://github.com/golang-standards/project-layout)**
* *Por que ler:* É o padrão de mercado. Explica o que deve ir dentro de `/internal` (código privado) e `/cmd` (entrypoints).


* **[The 12-Factor App (Em Português)](https://12factor.net/pt_br/)**
* *Por que ler:* Foca nos capítulos **III. Configurações** (variáveis de ambiente) e **VIII. Concorrência** (processos). É a base para rodar sua aplicação no Docker sem dor de cabeça.


* **[Clean Architecture in Go (Artigo Prático)](https://threedots.tech/post/introducing-clean-architecture/)**
* *Por que ler:* Mostra como separar o `Handler` (HTTP) do `Service` (Lógica), exatamente como desenhamos.



---

### 4. 🛠️ Ferramentas Específicas

Materiais sobre as ferramentas que vão processar seus memes.

* **[FFmpeg for Meme Makers (Guia Informal)](https://ffmpeg.org/documentation.html)**
* *Dica:* Não leia a documentação inteira (é gigante). Foque em comandos de conversão para `mp4` (H.264) e geração de thumbnail (`-vframes 1`).
* *Exemplo útil:* `ffmpeg -i input.mov -vcodec h264 -acodec aac output.mp4`


* **[Redis Crash Course](https://redis.io/docs/latest/develop/get-started/)**
* *Por que ler:* Para entender o básico de comandos como `SET`, `GET` e `EXPIRE` (para o cache dos posts não ficar velho).



---

### 5. 🔍 Observabilidade (Logs)

Para não programar no escuro.

* **[A Guide to Structured Logging in Go](https://betterstack.com/community/guides/logging/logging-in-go/)**
* *Por que ler:* Um tutorial moderno sobre como usar o novo pacote `log/slog` do Go 1.21+ para gerar JSON logs.


* **[Distributed Tracing - The Mental Model](https://www.google.com/search?q=https://www.honeycomb.io/blog/distributed-tracing-guide)**
* *Por que ler:* Explica o conceito de "Correlation ID" que mencionei, fundamental para saber que o erro no banco foi causado pelo Request X.



---

### 💡 Minha sugestão de ordem de estudo:

1. Comece pelo **Standard Go Project Layout** para criar as pastas certas.
2. Leia sobre **Pipelines** no blog do Go para implementar o Worker.
3. Use o **sqlc Playground** para validar suas queries.
4. Deixe o **FFmpeg** e **Redis** por último (primeiro faça o sistema funcionar só movendo arquivos e salvando no banco).
🏗️ Arquitetura Sistêmica: MemeDroid

O sistema é dividido em três domínios: Ingestão (Síncrona), Processamento (Assíncrona/Pipeline) e Manutenção (Estado).
1. Modelo de Fluxo e Componentes (Mermaid)
2. Definição dos Modelos de Dados (Data Domain)

Abaixo, a explicação teórica dos modelos que sustentam a aplicação:
A. Modelo de Identidade (Users)

    Propósito: Autenticação e Autoridade.

    Atributos Chave: Identificadores únicos (username, email) e segredo criptográfico (password_hash).

    Relacionamento: Um usuário é o "Owner" de múltiplos posts.

B. Modelo de Conteúdo (Posts)

    Propósito: Registro central da mídia.

    Máquina de Estados (Status):

        pending: Registro criado, arquivo recebido.

        processing: Worker assumiu a tarefa.

        completed: Mídia otimizada e disponível.

        failed: Erro no processamento (log disponível).

    Metadados: Caminhos para mídia original vs. otimizada e thumbnails.

C. Modelo de Taxonomia (Tags)

    Propósito: Indexação e descoberta.

    Arquitetura: Relacionamento N:N (Muitos-para-Muitos). Permite que um post tenha várias tags e uma tag organize vários posts.

3. Os Motores de Execução (The Engines)
I. O Pipeline de Processamento (Workers)

Em vez de processar o vídeo/imagem dentro da requisição HTTP (o que travaria o usuário), o sistema usa um Pipeline Assíncrono.

    Teoria: O Worker é um consumidor de recursos. Ele é isolado para que, se o FFmpeg travar, o servidor API continue aceitando novos uploads.

    Concorrência: Limitada pelo hardware (CPU cores), garantindo que o servidor nunca sofra de exaustão de memória.

II. O Finalizador (Consistência Eventual)

O Finalizer é o componente que traz ordem ao caos dos Workers.

    Teoria: Ele garante que o banco de dados só seja atualizado quando o arquivo final estiver seguro no disco. Ele é o responsável por deletar o "lixo" (arquivos brutos de upload).

III. O Janitor (Resiliência)

O Janitor resolve o problema do "Estado Zumbi".

    Teoria: Em computação distribuída, falhas são inevitáveis. Se um Worker sofrer um panic, o Post nunca sairia do estado processing. O Janitor é o auditor que limpa essas falhas após um tempo de tolerância (TTL).

4. O Sistema de Rastreabilidade (Observability)

    TraceID Contextual: Cada ação (do clique do usuário até o Janitor) carrega um identificador único no context.Context.

    Teoria: Isso permite o Log Correlation. Você consegue ver o "DNA" de uma falha através de todos os serviços, mesmo que o erro aconteça 10 minutos após o upload original.
