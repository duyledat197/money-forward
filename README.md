# Money Forward interview

## Requirements
The server api should serve for this endpoints:

* `GET /users/{id}`
* `GET /users/{id}/accounts`
* `GET /accounts/{id}`
  

All endpoints should return json payloads matching the following criteria.
`/users/{id}`: A detail of `users` including list `account_ids` that have been owned by the user that are equal to the passing `{id}` from the database.

Example:
```json
{
  "id": 1,
  "name": "Alice",
  "account_ids": [
    1,
    3,
    5,
    ...
  ]
}

``` 
`/users/{id}/accounts`: A list detail of `accounts` that have been owned by the user that are equal to the passing `{id}` from the database.

Example:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "A銀行",
    "balance": 20000
  },
  {
    "id": 3,
    "user_id": 1,
    "name": "C信用金庫",
    "balance": 120000
  },
  {
    "id": 5,
    "user_id": 1,
    "name": "E銀行",
    "balance": 5000
  },
  ...
  ]
``` 
`/accounts/{id}`: A detail of `accounts` that are equal to the passing `{id}` from the database.

Example:
```json
{
  "id": 2,
  "user_id": 2,
  "name": "Bカード",
  "balance": 200
}
``` 

## Features:

- [x] Finish three APIs in **Requirements**.
- [x] Adding mutation APIs for `accounts`, `users`.
- [x] Adding login flow.
- [x] Apply caching for optimize API.
- [x] Adding some middlewares hor http handle steps like `recovery`,`cors`,`authenticate`,`rbac`.


# Architecture: 
The Architecture using the [clean architecture](https://raw.githubusercontent.com/phungvandat/clean-architecture/dev/images/clean-arch.png) that control the system split to 3 layers.

![clean architecture](https://raw.githubusercontent.com/phungvandat/clean-architecture/dev/images/clean-arch.png)


# Built With

This section list any major frameworks/libraries used to bootstrap project.
* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
* ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) 


# Folder structure
```sh
.
├── LICENSE
├── Makefile
├── README.md
├── app-exe
├── cmd
│   ├── root.go
│   ├── srv.go
│   └── start.go
├── configs
│   ├── address.go
│   ├── configs.go
│   └── database.go
├── deployments
├── developments
│   ├── dev.env
│   └── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── deliveries
│   │   └── http
│   │       ├── account.go
│   │       ├── auth.go
│   │       └── user.go
│   ├── entities
│   │   ├── account.go
│   │   └── user.go
│   ├── models
│   │   ├── account.go
│   │   ├── auth.go
│   │   ├── common.go
│   │   └── user.go
│   ├── repositories
│   │   ├── account.go
│   │   └── user.go
│   └── services
│       ├── account.go
│       ├── auth.go
│       └── user.go
├── main.go
├── migrations
│   ├── 00001_migrate.up.sql
│   └── 00002_migrate.up.sql
└── pkg
    ├── cache
    │   └── cache.go
    ├── crypto_utils
    │   └── util.go
    ├── database
    │   ├── executor.go
    │   ├── type.go
    │   └── util.go
    ├── http_server
    │   ├── common.go
    │   ├── http.go
    │   ├── http_test.go
    │   ├── middleware.go
    │   ├── response.go
    │   ├── util.go
    │   ├── util_test.go
    │   └── xcontext
    │       ├── context.go
    │       └── ctx.go
    ├── id_utils
    │   ├── id.go
    │   └── snowflake.go
    ├── logger
    │   └── logger.go
    ├── lru
    │   └── cache.go
    ├── postgres_client
    │   ├── client.go
    │   └── tx.go
    ├── processor
    │   └── processor.go
    ├── reflect_utils
    │   ├── util.go
    │   └── util_test.go
    └── token_utils
        ├── authenticator.go
        ├── jwt.go
        └── paseto.go
```


# Prerequisites

- Make sure you have Go installed ([download](https://golang.org/dl/)). Version `1.21` or higher is required.
- Docker (version `20.10.22+`)

# Getting started
First of all, we should set up environments:

```sh
  copy ./developments/.env.example ./developments/.env
  copy ./developments/dev.env.example ./developments/dev.env
```

1. Starting database:

```sh
  make start-db
```

2. Migrate:

```sh
  make migrate
```

3. Simple start server:

```sh
  make start
```

# Action flows:

We already migrate a default super admin user by **admin** and **donkihote**.
```sh
  curl --location 'localhost:8080/auth/login' \
    --header 'Content-Type: application/json' \
    --data '{
        "user_name":"admin",
        "password": "donkihote"
    }'
```