# Go Rate Limiter

Uma implementação flexível de limitador de taxa em Go que suporta limitação baseada em IP e em token, utilizando Redis como backend de armazenamento.

## Funcionalidades

- Limitação de taxa baseada em IP
- Limitação de taxa baseada em token (via cabeçalho API_KEY)
- Configuração de requisições por segundo
- Configuração de duração do bloqueio
- Armazenamento em Redis
- Implementação de middleware para fácil integração
- Configuração baseada em variáveis de ambiente
- Containerização com Docker

## Pré-requisitos

- Docker e Docker Compose

## Instalação e Execução

1. Clone o repositório:
```bash
git clone https://github.com/ntayoane/go-rate-limited.git
cd go-rate-limited
```

2. Inicie a aplicação usando Docker Compose:
```bash
docker-compose up --build
```

A aplicação estará disponível em `http://localhost:8080`

## Configuração

O limitador de taxa pode ser configurado usando variáveis de ambiente no arquivo `docker-compose.yml`:

```yaml
environment:
  - REDIS_HOST=redis
  - REDIS_PORT=6379
  - REDIS_PASSWORD=
  - REDIS_DB=0
  - IP_REQUESTS_PER_SECOND=5
  - IP_BLOCK_DURATION_SECONDS=300
  - DEFAULT_TOKEN_REQUESTS_PER_SECOND=10
  - DEFAULT_TOKEN_BLOCK_DURATION_SECONDS=300
  - SERVER_PORT=8080
```

## Uso

Faça requisições ao servidor:
```bash
# Sem token (limitação por IP)
curl http://localhost:8080

# Com token (limitação por token)
curl -H "API_KEY: seu-token" http://localhost:8080
```

## Comportamento da Limitação de Taxa

- Se um token for fornecido no cabeçalho `API_KEY`, a limitação baseada em token é aplicada
- Se nenhum token for fornecido, a limitação baseada em IP é aplicada
- Quando o limite de taxa é excedido, o servidor responde com:
  - Código de status: 429 (Too Many Requests)
  - Mensagem: "you have reached the maximum number of requests or actions allowed within a certain time frame"

## Arquitetura

O projeto segue uma arquitetura limpa com os seguintes componentes:

- `internal/config`: Gerenciamento de configuração
- `internal/limiter`: Implementação do limitador de taxa
- `internal/middleware`: Middleware HTTP
- `internal/storage`: Interface de armazenamento e implementação Redis

## Estrutura do Docker

O projeto utiliza dois containers:

1. `app`: Container da aplicação Go
   - Construído a partir do Dockerfile
   - Expõe a porta 8080
   - Conecta-se ao Redis

2. `redis`: Container do Redis
   - Usa a imagem oficial do Redis
   - Persiste dados em um volume
   - Expõe a porta 6379

## Testes

Para executar os testes:
```bash
go test ./...
```