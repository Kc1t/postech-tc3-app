# RFC-0003: Estratégia de autenticação por CPF

- **Status:** Aceita
- **Data:** 2026-09
- **Autor:** Reservoir Devs
- **Decisões relacionadas:** ADR-0005 (JWT com refresh token), ADR-0007 (comunicação no gateway)

## Resumo

Propõe autenticação de clientes por **CPF**, com uma Lambda própria emitindo JWT HS256, em vez de Amazon Cognito.

## Motivação

O enunciado é específico:

> Criar uma Function Serverless para: validar o CPF do cliente; consultar a existência e o status do cliente na base de dados; gerar e devolver um token (JWT) válido para consumo das APIs protegidas.

Ao mesmo tempo, a aula "Realizando Autenticação e serviços de identificação" ensina **Amazon Cognito** como o caminho AWS para autenticar uma API, com JWT authorizer no API Gateway. Era preciso decidir explicitamente entre seguir a ferramenta da aula e cumprir o requisito como está escrito.

Há também um contexto herdado: a aplicação já tem autenticação por e-mail e senha com JWT e refresh token rotativo (ADR-0005), usada pelo painel administrativo da oficina. A autenticação por CPF **não substitui** essa — é um segundo público.

## Dois públicos, dois fluxos

| Público | Quem é | Fluxo | Rotas |
|---|---|---|---|
| **Operação** | Atendente, mecânico, administrador | e-mail + senha → JWT, `role: admin` ou `client` | `POST /api/v1/auth/login`, aberta no gateway |
| **Cliente** | Dono do veículo | CPF → JWT, `role: client` | `POST /auth`, aberta no gateway |

Ambos produzem um JWT HS256 assinado com o **mesmo** `JWT_SECRET`, e ambos são validados pelo mesmo middleware `Auth`.

## Proposta

**`POST /auth`** com corpo `{"cpf": "529.982.247-25"}`:

1. Normaliza e valida o CPF pelos dígitos verificadores — rejeita tamanho errado, dígito inválido e sequências repetidas.
2. Consulta `requesters` por `document`.
3. Verifica `status`.
4. Emite o JWT.

| Situação | HTTP |
|---|---|
| CPF válido, cliente ativo | 200 com `access_token`, `token_type`, `expires_at` |
| Payload malformado | 400 |
| CPF inválido | 400 |
| Cliente inexistente | 404 |
| Cliente inativo | 403 |
| Falha interna | 500 |

**Claims emitidas:** `sub` (id do solicitante), `role: client`, `email`, `name`, `document`, `iat`, `exp`.

`role` é obrigatório: o middleware `Auth` da aplicação rejeita com 401 qualquer token sem `sub` e `role` preenchidos. Omiti-lo faria todo token emitido pela Lambda ser recusado pela aplicação.

**Escopo do token:** não basta apresentar um JWT válido — o documento consultado precisa ser o do próprio token. `GET /service-orders/requester?document=X` com um token emitido para o CPF Y responde 403. Sem essa checagem, qualquer cliente autenticado leria as ordens de serviço de qualquer outro.

**Validade:** 15 minutos, sem refresh token. O fluxo do cliente é curto — consultar a OS, aprovar ou recusar o orçamento. Reautenticar com o CPF é barato e evita gerenciar revogação para esse público.

## Alternativa principal: Amazon Cognito

É o caminho ensinado na aula, e tem vantagens concretas: MFA, verificação de e-mail e telefone, integração federada, JWT authorizer nativo no gateway sem escrever código, e rotação de chaves gerenciada.

**Por que não foi adotado:**

1. **Cognito não modela CPF como credencial.** O modelo é usuário e senha, ou federação. Autenticar *só* com CPF exigiria tratá-lo como username com senha fixa ou vazia — o que é uma distorção do produto, não um uso dele.

2. **O requisito não pede um provedor de identidade.** Pede validar o CPF, consultar existência e status **na base de dados da aplicação** e devolver um JWT. A fonte de verdade do cliente é a tabela `requesters`, não um user pool. Com Cognito, os dados existiriam em dois lugares e precisariam ser sincronizados.

3. **Duplicaria a identidade já existente.** A aplicação tem `users` com papéis e bloqueio por tentativas falhas. Introduzir Cognito criaria um segundo sistema de identidade, com dois formatos de token e dois caminhos de validação.

4. **CPF sozinho não é segredo.** Esse é o ponto mais importante, e vale para qualquer implementação: CPF é um identificador público, não uma credencial. Autenticar apenas com CPF é fraco por natureza. O enunciado pede exatamente isso, então foi implementado assim — mas o risco está registrado abaixo, e Cognito não o resolveria: continuaria sendo autenticação por um dado público.

## Outras alternativas avaliadas

| Alternativa | Por que não |
|---|---|
| **JWT authorizer nativo do gateway** | Exige emissor OIDC com JWKS público. Nosso token é HS256 com segredo compartilhado. Detalhado no ADR-0007. |
| **CPF + código enviado por e-mail/SMS** | Resolveria a fragilidade de credencial. Descartado por prazo e por adicionar dependência de entrega de mensagem no caminho crítico da demonstração. |
| **RS256 com JWKS publicado** | Permitiria o JWT authorizer nativo e eliminaria o segredo compartilhado em três lugares. É a evolução natural; não cabe nos 11 dias. |
| **Reaproveitar `POST /api/v1/auth/login`** | Não atende: o requisito exige uma function serverless, e o login existente é e-mail e senha. |

## Riscos

| Risco | Severidade | Mitigação |
|---|---|---|
| **CPF é dado público — qualquer um com o CPF de um cliente obtém um token válido** | Alta | Registrado como limitação deliberada do escopo. O token dura 15 min e dá acesso apenas a `role: client`; rotas administrativas exigem `role: admin`, que a Lambda nunca emite. Evolução: segundo fator. |
| Enumeração de CPFs pelos códigos 404 e 403 | Média | Throttling no stage do gateway (100 req/s, burst 200). Uniformizar as respostas esconderia a distinção entre inexistente e inativo, que é útil ao cliente legítimo. |
| Segredo HS256 replicado em aplicação, emissor e authorizer | Média | Todos leem do mesmo GitHub Secret. Rotação exige deploy coordenado dos dois repositórios. |
| Cliente desativado continuar acessando até o token expirar | Baixa | Cache do authorizer em zero garante revalidação por requisição, mas o token em si segue válido por até 15 min. Aceito. |

## Questões em aberto

- Vale reduzir a validade para 5 minutos? Diminui a janela de um token vazado ao custo de mais reautenticações num fluxo já curto.
- A resposta 403 para cliente inativo deveria explicar o motivo? Hoje devolve apenas "cliente inativo", sem dizer por quê — decisão consciente para não vazar estado interno.
