# Chirpy API 

Uma API RESTful simulando o backend de uma plataforma social estilo Twitter, desenvolvida em Go como projeto de portfólio. O foco é implementar autenticação robusta (JWT + Refresh Tokens), operações CRUD completas e boas práticas de organização de código.

## Funcionalidades

*   **Autenticação Completa**: Registro de usuários, login com JWT e sistema de refresh tokens para renovação segura de sessão.
*   **Gerenciamento de "Chirps" (Posts)**: Criar, listar (com filtro por autor e ordenação), buscar e deletar posts.
*   **Upgrade de Usuário**: Webhook para atualizar um usuário para o status "premium" (autenticado).
*   **Banco de Dados PostgreSQL**: Persistência de dados com migrações versionadas.
*   **Segurança**: Senhas hasheadas com bcrypt, tokens assinados e validação rigorosa.
*   **Arquitetura Limpa**: Separação de responsabilidades com handlers, middleware e lógica de banco.

## Stack 

*   **Linguagem**: [Go](https://golang.org/)
*   **Banco de Dados**: [PostgreSQL](https://www.postgresql.org/)
*   **Gerenciamento de Migrações**: [Goose](https://github.com/pressly/goose)
*   **Geração de Queries Type-Safe**: [SQLC](https://sqlc.dev/)
*   **Autenticação**: JWT (lib `golang-jwt/jwt`) e [argon2id](https://github.com/alexedwards/argon2id)
*   **Dependências Adicionais**: `github.com/google/uuid`, `github.com/joho/godotenv`

## Como Executar o Projeto Localmente

### Pré-requisitos

*   Go (versão 1.21 ou superior)
*   PostgreSQL instalado e rodando localmente (ou um contêiner Docker manual)
*   Git

### Passo a passo

1.  **Clone o repositório**
    ```bash
    git clone https://github.com/luis-octavius/chirpy.git
    cd chirpy
    ```

2.  **Configure o banco de dados PostgreSQL**
    Crie um banco de dados chamado `chirpy` (ou o nome que preferir). Exemplo no psql:
    ```sql
    CREATE DATABASE chirpy;
    ```

3.  **Configure as variáveis de ambiente**
    Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo (substitua os valores):
    ```env
    DB_URL=postgres://seu_usuario:sua_senha@localhost:5432/chirpy?sslmode=disable
    JWT_SECRET=uma_chave_secreta_muito_segura_e_aleatoria
    POLKA_KEY=chave_para_o_webhook_de_upgrade (opcional por enquanto)
    PLATFORM=dev
    ```

4.  **Execute as migrações do banco de dados**
    ```bash
    cd sql/schema
    goose postgres "$DB_URL" up
    cd ../..
    ```

5.  **Inicie a aplicação**
    ```bash
    go run .
    ```
    A API estará disponível em `http://localhost:8080`.

## Documentação da API

A seguir, os principais endpoints disponíveis.

### Usuários

| Método | Endpoint       | Descrição                          | Corpo da Requisição (JSON)                          | Autenticação |
|--------|----------------|------------------------------------|-----------------------------------------------------|--------------|
| POST   | `/api/users`   | Cria um novo usuário               | `{"email": "user@example.com", "password": "123"}` | Não          |
| PUT    | `/api/users`   | Atualiza email/senha do usuário logado | `{"email": "novo@email.com", "password": "456"}`   | Sim (JWT)    |

### Autenticação (Tokens)

| Método | Endpoint         | Descrição                                         | Corpo da Requisição (JSON) | Autenticação              |
|--------|------------------|---------------------------------------------------|----------------------------|---------------------------|
| POST   | `/api/login`     | Faz login e retorna access_token e refresh_token  | `{"email": "...", "password": "..."}` | Não                       |
| POST   | `/api/refresh`   | Gera um novo access_token usando o refresh_token  | (vazio)                    | Sim (Refresh Token no header `Authorization: Bearer <token>`) |
| POST   | `/api/revoke`    | Revoga o refresh_token usado na requisição        | (vazio)                    | Sim (Refresh Token no header) |

### Chirps (Posts)

| Método | Endpoint                  | Descrição                                                     | Parâmetros de Query | Autenticação |
|--------|---------------------------|---------------------------------------------------------------|---------------------|--------------|
| POST   | `/api/chirps`             | Cria um novo chirp (o autor é o usuário do token)             | -                   | Sim (JWT)    |
| GET    | `/api/chirps`             | Lista todos os chirps                                         | `?author_id=UUID` (opcional) `&sort=asc/desc` | Não          |
| GET    | `/api/chirps/{chirpID}`   | Busca um chirp específico pelo seu ID                         | -                   | Não          |
| DELETE | `/api/chirps/{chirpID}`   | Deleta um chirp (apenas o autor ou um admin podem deletar)    | -                   | Sim (JWT)    |

### Webhooks e Saúde

| Método | Endpoint                          | Descrição                                                     | Autenticação |
|--------|-----------------------------------|---------------------------------------------------------------|--------------|
| POST   | `/api/polka/webhooks`             | Webhook para upgrade de usuário (espera um evento específico) | Sim (API Key)|
| GET    | `/api/healthz`                    | Verifica se a API está rodando (útil para health checks)      | Não          |
| GET    | `/admin/metrics`                  | Retorna métricas básicas (número de visits, etc.)             | Não          |
| GET    | `/app/*`                          | Serve arquivos estáticos (assets HTML, etc.)                  | Não          |

*(Nota: O endpoint de admin e assets usam um contador de requests, visível na raiz `/`)*

## Estrutura do Projeto

```text
├── assets/ # Arquivos estáticos (CSS, JS para a página de admin)
├── internal/ # Código privado da aplicação
│ └── database/ # Lógica de acesso a dados gerada pelo SQLC
├── sql/
│ ├── queries/ # Arquivos .sql com as queries nomeadas para o SQLC
│ └── schema/ # Migrações SQL (arquivos .sql versionados)
├── handler_chirps.go # Handlers para endpoints de chirps
├── handler_healthz.go # Handler para health check
├── handler_tokens.go # Handlers para login, refresh, revoke
├── handler_users.go # Handlers para criação e atualização de usuários
├── middleware.go # Middlewares (autenticação, validação, contagem)
├── middleware_test.go # Testes para os middlewares
├── main.go # Ponto de entrada: configura rotas e servidor
├── sqlc.yaml # Configuração do SQLC
├── go.mod # Dependências do módulo Go
└── README.md # Este arquivo
```

## Testes

O projeto possui testes unitários, principalmente para os middlewares.

```bash
go test -v
```

Para rodar testes de um arquivo específico (ex: middleware):
```bash

go test -v ./... -run TestValidateMessage
```

## Próximos Passos

- [ ] Adicionar Docker e docker-compose.yml para facilitar a execução do ambiente (banco + app).
- [ ] Implementar testes de integração para os handlers.
- [ ] Criar um front-end simples para consumir a API (já existe um esboço em index.html e assets).
- [ ] Adicionar rate limiting para prevenir abusos.

## Contribuição

Este é um projeto de portfólio pessoal. Feedbacks, sugestões de melhoria ou relatos de bugs são muito bem-vindos através de Issues!

## Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

Desenvolvido por Luis Octávio como parte de seus estudos em desenvolvimento backend com Go.
