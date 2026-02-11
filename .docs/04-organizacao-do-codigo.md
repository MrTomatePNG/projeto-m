# 4. Organização do Código (Project Structure)

A estrutura do projeto segue o **"Standard Go Project Layout"**, uma convenção da comunidade Go que promove modularidade, encapsulamento e clareza. Isso torna o projeto mais fácil de navegar e manter.

## 4.1 Estrutura de Diretórios Principal

```
/home/tomate/Lixeira/Go/webflix/
├── cmd/               # Executáveis principais da aplicação.
│   └── api/           # O ponto de entrada (main.go) da API.
├── internal/          # Código privado do projeto.
├── sql/               # Definições e migrations SQL.
├── .docs/             # Documentação do projeto.
├── go.mod             # Gerenciamento de dependências.
└── docker-compose.yml # Orquestração de serviços para desenvolvimento.
```

## 4.2 Detalhamento dos Diretórios

### `cmd/`
Contém as aplicações executáveis do projeto. Cada subdiretório representa um ponto de entrada (`main.go`) distinto.
*   **`cmd/api/`**: É o "entrypoint" para o servidor da API web. É aqui que os servidores HTTP são iniciados, os `workers` do pipeline são "spawnados" e as configurações são carregadas.

### `internal/`
Destinado a todo o código que é **privado** a este projeto.
*   **Por que `internal`?** O compilador Go impõe uma regra: pacotes dentro de `internal/` **não podem** ser importados por outros projetos. Isso garante um forte encapsulamento e deixa claro quais partes do código são específicas da aplicação e não devem ser reutilizadas externamente.
*   **Subdiretórios:**
    *   **`database/`**: Contém o código gerado pelo `sqlc` e as interfaces para interagir com o banco de dados.
    *   **`handlers/`**: Camada HTTP. Responsável por decodificar requisições HTTP, validar a entrada, chamar os `services` e formatar as respostas JSON.
    *   **`middleware/`**: Contém middlewares HTTP, como logging de requisições, recuperação de `panics` e manipulação de `Correlation ID`.
    *   **`services/`**: Representa a camada de regras de negócio. Este código deve ser "puro", ou seja, não deve ter conhecimento direto sobre HTTP ou detalhes de implementação do banco de dados. Ele orquestra a lógica de negócio principal.
    *   **`utils/`**: Pacote para funções utilitárias genéricas que podem ser usadas em várias partes do projeto (ex: manipulação de strings, helpers de criptografia).
    *   **`workers/`**: Contém toda a lógica do pipeline de processamento assíncrono: o `Dispatcher`, o `Worker Pool`, o `Finalizer` e o `Janitor`.

### `sql/`
Armazena todas as definições e scripts SQL, separando-os claramente do código Go.
*   **`migrations/`**: Contém os scripts de migração do banco de dados (esquema `up` e `down`), gerenciados por ferramentas como `golang-migrate`.
*   **`queries/`**: Arquivos `.sql` contendo as queries que o `sqlc` usará para gerar o código Go tipado.
*   **`schemas/`**: Arquivos `.sql` que definem o schema do banco de dados, usados para referência e para a geração do `sqlc`.

### `.docs/`
Este diretório contém a documentação detalhada do projeto (como o arquivo que você está lendo agora), incluindo diagramas de arquitetura e explicações conceituais.
