# Postman - HelpTech OS

Esta pasta já vem pronta para você **importar no Postman e testar o projeto sem frontend**.

## Arquivos

- `helptech-os.postman_collection.json`: collection completa com requests e testes automatizados.
- `helptech-os.local.postman_environment.json`: ambiente local com URLs e variáveis usadas pela collection.

## Como importar

1. Abra o Postman.
2. Clique em **Import**.
3. Importe os dois arquivos desta pasta:
   - `helptech-os.postman_collection.json`
   - `helptech-os.local.postman_environment.json`
4. No canto superior direito do Postman, selecione o environment **HelpTech OS Local**.
5. Suba o projeto com Docker:

```bash
cd helptech-os-go-grpc
docker compose up --build
```

## Ordem recomendada para executar

1. Pasta **00 - Health checks**: confirma que gateway e serviços estão no ar.
2. Pasta **01 - Auth**:
   - Execute **Register**.
   - Execute **Login - salva token**.
   - O teste do login salva automaticamente o JWT nas variáveis `token` da collection e do environment.
3. Pasta **02 - Fluxo completo OS**:
   - Criar cliente
   - Listar clientes
   - Buscar cliente
   - Atualizar cliente
   - Criar ordem
   - Listar ordens
   - Buscar ordem
   - Atualizar ordem
   - Criar pagamento
   - Listar pagamentos
   - Buscar pagamento
4. Pasta **03 - Limpeza opcional** somente se quiser apagar cliente/ordem criados.

## O que os testes verificam

A collection contém scripts na aba **Tests** para validar:

- status HTTP esperado (`200`, `201`, `401` ou `409` quando registrar usuário já existente);
- resposta JSON de health check com `status: ok` e nome correto do serviço;
- token JWT retornado no login;
- salvamento automático de `token`, `cliente_id`, `ordem_id` e `pagamento_id`;
- listagens retornando arrays;
- buscas por ID e atualizações básicas.

## Observação importante

Se você executar o **Register** mais de uma vez com o mesmo e-mail, é normal retornar conflito/usuário já existente. A collection aceita `201` ou `409` nesse request para facilitar o reuso durante estudos.
