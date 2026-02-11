# 2. Arquitetura de Dados (Persistence Layer)

A camada de persistência é o "coração" da aplicação, responsável por armazenar a "fonte da verdade" e otimizar o acesso aos dados para performance.

## 2.1 Tecnologias Escolhidas

A estratégia de dados utiliza duas tecnologias complementares:

*   **PostgreSQL:** É a **fonte da verdade** principal do sistema.
    *   **Por que?** Escolhido por sua robustez, integridade referencial (garantindo relacionamentos consistentes), e suporte a tipos de dados complexos como `JSONB` e `Arrays`. É ideal para dados transacionais e relacionais críticos. A integração com `sqlc` também é um fator decisivo.

*   **Redis:** Utilizado como **cache de propósito específico** no padrão Cache-Aside.
    *   **Por que?** Para aliviar a carga do PostgreSQL em operações de leitura muito frequentes e que não precisam ser sempre consistentes em tempo real, como timelines de posts ou contadores. Sua alta velocidade o torna perfeito para esse fim.

## 2.2 Modelagem Relacional (PostgreSQL)

A modelagem reflete as entidades centrais e seus relacionamentos:

### 2.2.1 `Users` (Usuários)
*   **Propósito:** Gerenciar autenticação e autorização.
*   **Atributos Chave:**
    *   `username` e `email`: Devem ter constraint `UNIQUE` para garantir a unicidade do usuário.
    *   `password_hash`: Armazena o hash da senha gerado via `Bcrypt`, nunca a senha em texto claro.
*   **Relacionamento:** Um `User` é o proprietário (`Owner`) de múltiplos `Posts`.

### 2.2.2 `Posts` (Publicações)
*   **Propósito:** Registrar e gerenciar todo o conteúdo de mídia carregado.
*   **Máquina de Estados (`status` ENUM):** O campo `status` é fundamental para rastrear o ciclo de vida da mídia:
    *   `pending`: Registro criado, aguardando processamento.
    *   `processing`: Um `Worker` assumiu a tarefa.
    *   `completed`: Mídia otimizada com sucesso e disponível.
    *   `failed`: Erro no processamento (logs devem conter a causa).
*   **`media_hash` (SHA-256):**
    *   **Por que?** Permite a **deduplicação de mídia**. Evita reprocessar o mesmo meme viral repetidamente, economizando recursos de CPU e armazenamento.
*   **Denormalização para Leitura (`json_agg`):**
    *   **Por que?** Na leitura de posts, utilizamos a função de agregação `json_agg` do PostgreSQL para retornar um post juntamente com suas tags em uma única query. Isso resolve o clássico **problema de N+1 queries**, otimizando drasticamente a performance de leitura.
*   **Metadados:** Inclui caminhos para a mídia original, otimizada (`webp`/`mp4`) e thumbnails.

### 2.2.3 `Tags` e `post_tags` (Etiquetas)
*   **Propósito:** Indexar e categorizar o conteúdo, facilitando a descoberta.
*   **Arquitetura:** Relacionamento **N:N (Muitos-para-Muitos)**. A tabela `post_tags` atua como uma tabela ponte.
*   **Otimização:** Índices em ambas as direções (`post_id` e `tag_id`) garantem buscas bidirecionais rápidas (ex: encontrar todas as tags de um post, ou todos os posts com uma tag).

## 2.3 Padrão de Acesso a Dados (sqlc)

Adotamos o `sqlc` para interagir com o PostgreSQL.

*   **Como funciona?** `sqlc` lê arquivos SQL puros (schemas e queries) e gera código Go idiomático e tipado para essas operações.
*   **Por que?**
    1.  **Performance:** As queries são executadas como SQL nativo, sem a sobrecarga de um ORM complexo.
    2.  **Segurança de Tipos:** O código gerado é totalmente tipado, o que previne uma classe inteira de bugs em tempo de compilação.
    3.  **Prevenção de SQL Injection:** O uso de parâmetros posicionais (`$1`, `$2`) e a geração de código pelo `sqlc` tornam a injeção de SQL muito mais difícil.
    4.  **Clareza:** O SQL fica separado do código Go, em arquivos `.sql`, facilitando a análise e otimização por DBAs.

## 2.4 Estratégia de Armazenamento (Sharding de Pastas)

Para o armazenamento de arquivos em disco, a estrutura de pastas segue um padrão de **Sharding**.
*   **Exemplo:** Em vez de salvar tudo em `/storage`, um arquivo é salvo em um caminho como `/storage/u/e9/9a/...`
*   **Por que?** A maioria dos sistemas de arquivos (filesystems) começa a ter degradação de performance quando um único diretório contém dezenas ou centenas de milhares de arquivos. O sharding distribui os arquivos em uma árvore de diretórios, mantendo a performance do sistema de arquivos alta e previsível.
