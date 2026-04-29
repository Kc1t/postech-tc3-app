# ADR-0005: Autenticação JWT com refresh token rotativo

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

O Tech Challenge exige **autenticação JWT para APIs administrativas**. O grupo precisa decidir:

- Estrutura do token (qual padrão de claims, tempo de expiração).
- Como permitir renovação sem forçar login frequente.
- Como revogar acesso de um usuário comprometido.
- Como tratar o portal público de consulta de OS (que **não** exige autenticação).

## Decisão

Implementar autenticação em duas camadas:

### 1. Access token (JWT)

- Tipo: `Bearer` no header `Authorization`.
- TTL: **15 minutos**.
- Algoritmo: `HS256` com chave secreta carregada de variável de ambiente.
- Claims: `sub` (user id), `email`, `role` (`admin` ou `client`), `iat`, `exp`.

### 2. Refresh token rotativo

- TTL: **7 dias**.
- Armazenado no banco (`refresh_tokens`) com referência ao usuário, IP de emissão e hash do token.
- A cada uso, o refresh token antigo é **revogado** e um novo é emitido.
- Endpoint público: `POST /api/v1/auth/refresh`.
- Operação atômica no repositório (`RotateToken`) — ver [ADR-0007](./0007-domain-errors-sentinels.md) sobre transações no repositório.

### 3. Logout explícito

`POST /api/v1/auth/logout` revoga o refresh token atual. Access tokens não são revogáveis (ficam válidos até `exp`); a janela de exposição é mitigada pelo TTL curto.

### 4. Account lockout

Após **5 tentativas de login falhadas** em sequência, a conta é bloqueada por **15 minutos** (campos `failed_attempts` e `locked_until` na entidade `User`). Esta lógica vive **no domínio**, não em middleware.

### 5. Portal público

Endpoints `GET /api/v1/service-orders/code/:code`, `GET /api/v1/service-orders/customer` e `PUT /api/v1/service-orders/code/:code/status` são **explicitamente públicos** (sem middleware de auth). Identificação se dá pelo **código da OS** (token de acesso implícito que o cliente recebeu).

## Alternativas consideradas

1. **Sessions com cookies HttpOnly.** Mais simples para um SPA do mesmo domínio, mas o TC pode evoluir para múltiplos clientes (mobile, parceiros). JWT funciona em todos os casos.
2. **Apenas access token de longa duração.** Solução comum em prototipagem, mas permite que um token vazado dure dias. Inaceitável.
3. **OAuth2 / OIDC com provider externo (Auth0, Keycloak).** Excelente para produção, porém adiciona dependência externa fora do escopo de um MVP.
4. **Refresh tokens não rotativos (mesmo refresh válido por todo TTL).** Simples mas perde a capacidade de detectar reuso (token roubado pode ser usado várias vezes). Optamos pelo rotativo.

## Consequências

### Positivas

- TTL curto do access token reduz janela de exploração.
- Rotação do refresh permite **detectar reuso** (se um refresh já revogado for apresentado, é sinal de roubo — mitigado revogando todos os tokens do usuário).
- Account lockout protege contra brute-force de senha sem depender de WAF/Cloudflare.
- Roles em claim simplificam autorização nos middlewares (`requireRole("admin")`).
- Portal público para clientes não exige cadastro — UX simplificada.

### Negativas

- Stateless do access token impede revogação imediata; aceito dado o TTL de 15 min.
- Refresh token rotativo exige cuidado com **race conditions**: se duas requisições do mesmo cliente apresentarem o mesmo refresh simultaneamente, uma falha. Mitigado pela operação atômica `RotateToken` (transação no repositório).
- Account lockout pode ser usado como vetor de **DoS** (atacante trava contas legítimas). Mitigação adicional (CAPTCHA, lockout por IP) fica para fases futuras.

## Notas de implementação

- Implementação do `TokenService` em `internal/adapters/outbound/jwt/service.go` (não em `pkg/`, pois depende de `entities.User`).
- Use cases recebem `ports.TokenService` (interface), nunca a biblioteca de JWT diretamente — preserva o domínio limpo (ver [ADR-0001](./0001-hexagonal-architecture.md)).
- `BCRYPT_COST=12` para hashing de senhas (configurável via env).

## Referências

- RFC 7519 (JSON Web Token).
- Auth0 — *Refresh Token Rotation* (https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation).
