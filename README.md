# Api-billing

Provide an API to get info about tenant

## Mapa Mental do Sistema de Bilhetagem

![image](https://github.com/user-attachments/assets/0ed81418-77cd-473a-9336-7b17d5969f64)

## Estrutura do Projeto

```bash
├── internals
│   ├── api
│   │   ├── tenantHandler.go
│   ├── app
│   │   ├── app.go
│   ├── middleware
│   │   ├── middleware.go
│   ├── routes
│   │   ├── routes.go
│   ├── utils
│   │   ├── utils.go
│   ├── store
│   │   ├── database.go
│   │   ├── tenant_store.go
├── main.go
├── go.mod
├── go.sum
├── .env
├── .gitignore
├── README.md
```

### Variáveis de Ambiente

```bash
MONGO_URL=mongodb://localhost:27017
MONGO_DATABASE="magalu_billing"
MONGO_COLLECTION_AGGREGATION="magalu_billing_aggregation"
JWT_SIGNING_KEY="MCX1s22aiFnkN/X05iXZPZCWVz1DTzOWgjG7npuAX80="
```

## Air

Air é uma ferramenta de hot reload para aplicações Go. Ele monitora as alterações nos arquivos e reinicia automaticamente o servidor, facilitando o desenvolvimento.

### Instalação

#### Via `go install` (Recommended)

Com go 1.23 ou superior:

```bash
go install github.com/air-verse/air@latest
```

## Iniciando o Projeto

1. Clone o repositório:

```bash
git clone
cd api-billing
```

2. Instale as dependências:

```bash
go mod tidy
```

3. Crie um arquivo `.env` na raiz do projeto e adicione as variáveis de ambiente necessárias. Você pode usar o arquivo `.env.example` como referência.
4. Inicie o MongoDB:

```bash
docker compose up -d
```

5. Popular a base de dados com os dados de exemplo (opcional):

Ir na pasta `import` e executar o script `import.go`:

Ele irá criar uma collection de pulse com tenants e sku aleatórios.
Após isso ele realizará o processo de agregação e irá criar uma collection com o resultado.

5. Inicie o projeto:

```bash
air
```

6. Acesse a API em `http://localhost:8081`
7. Para testar a API, tem alguns exemplos usando curl no arquivo `notes.txt` na raiz do projeto.

## Criar um Token JWT

Para criar um token JWT, você pode usar o site [jwt.io](https://www.jwt.io).

Basta passar os campos que estão no struct `CustomClaims` no `middlewae.go`, a JWT_SIGNING_KEY e o algoritmo HS256.

## To Do List

- [x] Criar um endpoint para retornar o consolidade de consumo por tenant
- [x] Criar um endpoint para retornar o consumo de um SKU do tenant
- [x] Criar um endpoint para retornar uma lista de todos os sku do tenant
- [x] Adicionar paginação na rota de listagem de sku
- [x] Adicionar variável de ambiente para o MONGO_HOST e as collections
- [x] Criar um middleware para validar o token JWT e adicionar o username no context
