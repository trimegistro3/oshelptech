# HelpTech OS Go gRPC

Projeto de estudo para um Sistema de Ordem de Serviço usando Go, microserviços, Docker Compose, PostgreSQL, GORM, JWT e Postman.

## Containers

O `docker-compose.yml` sobe 7 containers: `postgres`, `pgadmin`, `auth-service`, `cliente-service`, `ordem-service`, `financeiro-service` e `api-gateway`.

## Ordem de subida

O PostgreSQL sobe primeiro e tem healthcheck com `pg_isready`. Os serviços Go dependem do PostgreSQL saudável. O API Gateway depende dos quatro serviços Go saudáveis.

## Configuração

```bash
cp .env.example .env
```

## Rodar

```bash
docker compose up --build
```

## Parar

```bash
docker compose down
```

## Apagar banco e volumes

```bash
docker compose down -v
```

## pgAdmin

Acesse `http://localhost:5050` com email `admin@helptech.com` e senha `admin123`.

## API Gateway

Base URL: `http://localhost:8080`.

## Login e JWT

1. Faça `POST /auth/register`.
2. Faça `POST /auth/login`.
3. Copie o campo `token` retornado.
4. No Postman, envie `Authorization: Bearer SEU_TOKEN` nas rotas protegidas.

## Exemplos JSON

### Register

```json
{ "nome": "Allan Alex", "email": "allan@helptech.com", "password": "123456" }
```

### Login

```json
{ "email": "allan@helptech.com", "password": "123456" }
```

### Cliente

```json
{ "nome": "João Silva", "telefone": "31999999999", "email": "joao@email.com", "cpf_cnpj": "12345678900" }
```

### Ordem

```json
{ "cliente_id": 1, "equipamento": "Notebook Dell Inspiron", "marca": "Dell", "modelo": "Inspiron 15", "defeito": "Não liga", "diagnostico": "Aguardando análise", "status": "aberta", "valor": 390.00 }
```

### Pagamento

```json
{ "ordem_id": 1, "valor": 390.00, "forma_pagamento": "pix", "status": "pago" }
```

## AutoMigrate

O AutoMigrate é uma função do GORM que cria ou atualiza tabelas automaticamente conforme as structs Go de cada serviço.

## GORM

GORM é um ORM para Go. Ele permite salvar, buscar, atualizar e excluir dados do PostgreSQL usando structs Go.

## Docker Compose

Docker Compose permite declarar vários containers em um arquivo e subir tudo com um único comando.

## API Gateway

O API Gateway é a porta de entrada da aplicação. Ele recebe chamadas externas, valida JWT nas rotas protegidas e encaminha HTTP para os serviços internos.

## Postman

Importe `docs/postman/helptech-os.postman_collection.json` no Postman para testar as rotas.
