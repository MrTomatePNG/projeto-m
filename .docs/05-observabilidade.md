# 5. Observabilidade

Observabilidade é a capacidade de entender o que está acontecendo dentro do sistema apenas observando suas saídas externas (logs, métricas, traces). Em um sistema distribuído e assíncrono como o MemeDroid, a observabilidade não é um luxo, mas uma necessidade fundamental para a depuração e monitoramento.

## 5.1 Logs Estruturados (Structured Logging)

*   **Padrão:** Todos os logs gerados pela aplicação devem ser **logs estruturados em formato JSON**. Em Go, utilizamos o pacote `log/slog`.
*   **Por que JSON?**
    *   **Parseabilidade:** Logs em formato JSON são facilmente parseáveis e indexáveis por ferramentas de análise de logs como ELK Stack (Elasticsearch, Logstash, Kibana), Grafana Loki, Splunk, ou Datadog.
    *   **Buscas Complexas:** Permitem buscas e agregações complexas. Por exemplo: "mostre todos os logs de erro para o `worker-id: 3` que ocorreram na última hora" ou "calcule a duração média do processamento de vídeo".
    *   **Clareza:** Cada parte da informação de log tem uma chave (`key`), o que elimina a ambiguidade de logs em texto plano.

**Exemplo de log estruturado vs. não estruturado:**

```
// Não estruturado (ruim)
log.Printf("Erro ao processar vídeo %d: %v", videoID, err)

// Estruturado (bom)
slog.Error("erro ao processar vídeo", "video_id", videoID, "error", err)
```
**Saída JSON do log estruturado:**
```json
{
  "time": "2023-10-27T10:00:00.000Z",
  "level": "ERROR",
  "msg": "erro ao processar vídeo",
  "video_id": 12345,
  "error": "codec not supported"
}
```

## 5.2 Rastreabilidade com Correlation ID

*   **O Problema:** Em um sistema assíncrono, uma única ação do usuário (como um upload) desencadeia uma cadeia de eventos que podem ocorrer em diferentes goroutines e até mesmo em diferentes serviços, muitas vezes minutos após a requisição original. Como conectar um erro no `Finalizer` com a requisição HTTP que o originou?
*   **A Solução: `Correlation ID`** (também conhecido como `Trace ID`).
    1.  Um ID único é gerado no início de cada requisição HTTP (geralmente em um middleware).
    2.  Este ID é inserido no `context.Context` da requisição.
    3.  O `context` é propagado através de todas as camadas da aplicação: do `handler` para o `service`, e então incluído no `MediaJob` que é enviado para o pipeline.
    4.  Cada log gerado em qualquer ponto dessa cadeia de eventos (no `Worker`, no `Finalizer`, etc.) **deve incluir o `Correlation ID`**.

*   **Resultado:** Ao procurar por um `Correlation ID` específico nos logs, é possível ver a história completa e ordenada de uma operação, não importa quão distribuída ou assíncrona ela seja. Isso transforma a depuração de "procurar uma agulha no palheiro" para "seguir uma trilha de migalhas de pão".

## 5.3 Níveis de Log

Utilizamos níveis de log padrão para categorizar a severidade e a importância das mensagens, permitindo filtrar o ruído em produção.

*   `DEBUG`: Informações detalhadas para depuração em ambiente de desenvolvimento. (Ex: "Worker pegou novo job").
*   `INFO`: Fluxo normal e eventos importantes da aplicação. (Ex: "Upload recebido", "Job finalizado com sucesso").
*   `WARN`: Eventos que não são erros, mas que podem indicar um problema potencial. (Ex: "Redis cache miss", "Tentativa de reprocessamento").
*   `ERROR`: Falhas recuperáveis ou eventos que indicam um problema que precisa de atenção, mas que não interrompem o funcionamento do serviço. (Ex: "Codec inválido", "Falha ao conectar ao cache Redis").
*   `FATAL` / `CRITICAL`: Falhas críticas que indicam um estado irrecuperável e que exigem a interrupção da aplicação. (Ex: "Banco de dados fora do ar", "Falha ao alocar memória").
