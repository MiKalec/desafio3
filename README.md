# Desafio 3 - FullCycle

## Enunciado

Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

## Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.24+ (apenas se quiser executar localmente sem Docker)

## Como executar

### Subindo dependências com Docker Compose

1. Clone o repositório:
```bash
git clone <repository-url>
cd desafio3
```

2. Execute o docker compose para subir as dependências (MySQL e RabbitMQ):
```bash
docker compose up
```

Isso irá iniciar:
- **MySQL** na porta `3306`
- **RabbitMQ** nas portas `5672` (AMQP) e `15672` (Management UI)

### Executando a aplicação (go run)

1. Certifique-se de que o MySQL e RabbitMQ estão rodando (via `docker compose up` ou localmente).

2. Configure as variáveis de ambiente no arquivo `cmd/ordersystem/.env`:
```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=:8000
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```

3. Em outro terminal (na pasta do projeto), execute a aplicação:
```bash
go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go
```

## Portas dos Serviços

- **REST API**: `http://localhost:8000`
- **gRPC Server**: `localhost:50051`
- **GraphQL Server**: `http://localhost:8080`
  - GraphQL Playground: `http://localhost:8080`
  - GraphQL Endpoint: `http://localhost:8080/query`
- **MySQL**: `localhost:3306`
- **RabbitMQ Management UI**: `http://localhost:15672` (usuário: guest, senha: guest)

## Endpoints REST

### Criar uma Order

**HTTP:**
```http
POST http://localhost:8000/order
Content-Type: application/json

{
  "id": "123",
  "price": 10.0,
  "tax": 2.0
}
```

**cURL:**
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{
    "id": "123",
    "price": 10.0,
    "tax": 2.0
  }'
```

**Resposta esperada:**
```json
{
  "id": "123",
  "price": 10.0,
  "tax": 2.0,
  "final_price": 12.0
}
```

### Listar todas as Orders

**HTTP:**
```http
GET http://localhost:8000/order
```

**cURL:**
```bash
curl -X GET http://localhost:8000/order
```

**Resposta esperada:**
```json
[
  {
    "id": "123",
    "price": 10.0,
    "tax": 2.0,
    "final_price": 12.0
  },
  {
    "id": "456",
    "price": 20.0,
    "tax": 3.0,
    "final_price": 23.0
  }
]
```

## Serviços gRPC

### CreateOrder
```protobuf
rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
```

### ListOrders
```protobuf
rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
```

Para testar os serviços gRPC, você pode usar ferramentas como `grpcurl` ou criar um cliente gRPC.

Exemplo com grpcurl:
```bash
# Listar serviços disponíveis
grpcurl -plaintext localhost:50051 list

# Criar uma order
grpcurl -plaintext -d '{"id": "123", "price": 10.0, "tax": 2.0}' localhost:50051 pb.OrderService/CreateOrder

# Listar orders
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

## Queries GraphQL

### Criar uma Order (Mutation)

**GraphQL Playground:**
```graphql
mutation {
  createOrder(input: {
    id: "123"
    Price: 10.0
    Tax: 2.0
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

**cURL:**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createOrder(input: { id: \"123\", Price: 10.0, Tax: 2.0 }) { id Price Tax FinalPrice } }"
  }'
```

**Resposta esperada:**
```json
{
  "data": {
    "createOrder": {
      "id": "123",
      "Price": 10.0,
      "Tax": 2.0,
      "FinalPrice": 12.0
    }
  }
}
```

### Listar todas as Orders (Query)

**GraphQL Playground:**
```graphql
query {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

**cURL:**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { orders { id Price Tax FinalPrice } }"
  }'
```

**Resposta esperada:**
```json
{
  "data": {
    "orders": [
      {
        "id": "123",
        "Price": 10.0,
        "Tax": 2.0,
        "FinalPrice": 12.0
      },
      {
        "id": "456",
        "Price": 20.0,
        "Tax": 3.0,
        "FinalPrice": 23.0
      }
    ]
  }
}
```

Você pode testar as queries GraphQL acessando o GraphQL Playground em `http://localhost:8080`.

## Arquivo api.http

O arquivo `api.http` contém exemplos de requisições REST que podem ser usadas com extensões como REST Client do VS Code ou IntelliJ HTTP Client.

## Estrutura do Projeto

```
desafio3/
├── cmd/ordersystem/          # Entry point da aplicação
├── internal/
│   ├── entity/               # Entidades de domínio
│   ├── usecase/              # Casos de uso (regras de negócio)
│   ├── event/                # Eventos de domínio
│   └── infra/                # Implementações de infraestrutura
│       ├── database/         # Repositório SQL
│       ├── web/              # Handlers REST
│       ├── grpc/             # Serviços gRPC
│       └── graph/            # Schema e resolvers GraphQL
├── pkg/events/               # Dispatcher de eventos
├── configs/                  # Configurações
├── sql-scripts/              # Scripts SQL de inicialização
├── Dockerfile                # Dockerfile da aplicação
├── docker-compose.yaml       # Configuração Docker Compose
└── api.http                  # Exemplos de requisições REST
```

## Migrações

As migrações do banco de dados são executadas automaticamente quando o container MySQL é iniciado pela primeira vez através do script `sql-scripts/init_table.sql`, que cria a tabela `orders` com os seguintes campos:

- `id` (VARCHAR(255), PRIMARY KEY)
- `price` (FLOAT, NOT NULL)
- `tax` (FLOAT, NOT NULL)
- `final_price` (FLOAT, NOT NULL)

## Testes

Para executar os testes:
```bash
go test ./...
```

## Regenerar código gerado

### gRPC
```bash
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
```

### GraphQL
```bash
go run github.com/99designs/gqlgen generate
```

### Wire (Dependency Injection)
```bash
go run github.com/google/wire/cmd/wire ./cmd/ordersystem
```
