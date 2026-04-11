# TASK: Autenticacao JWT Profissional

Responsavel: Kauam
Escopo: Auth completo com refresh token rotation, logout, roles e rate limiting.

---

## Visao geral

```
POST /auth/register   → cria usuario (publico)
POST /auth/login      → retorna access_token (15min) + refresh_token (7d)
POST /auth/refresh    → rotaciona tokens (refresh antigo morre)
POST /auth/logout     → invalida refresh token
```

Roles: `admin` e `client` (hardcoded, sem RBAC).

---

## Dependencias a instalar

```bash
go get golang.org/x/crypto/bcrypt       # hash de senha
go get github.com/google/uuid           # gerar ID do refresh token
```

Nao precisamos de lib de rate limiting — vamos fazer com middleware puro usando `sync.Map` + `time.Ticker` (idiomatico em Go, sem dep externa).

---

## FASE 1 — Domain + Persistencia

### 1.1 Entidade User

**Arquivo:** `internal/domain/user/entity.go`

```go
type Role string

const (
    RoleAdmin  Role = "admin"
    RoleClient Role = "client"
)

type User struct {
    id           string
    name         string
    email        string
    passwordHash string
    role         Role
    createdAt    time.Time
    updatedAt    time.Time
}
```

Seguir o padrao existente:
- Campos privados + getters
- `New(name, email, passwordHash string, role Role) *User`
- `Reconstitute(id, name, email, passwordHash string, role Role, createdAt, updatedAt time.Time) *User`
- `SetID(id string)`
- Sem setter pra `role` (imutavel apos criacao)
- Sem setter pra `email` (por enquanto — MVP)

### 1.2 Entidade RefreshToken

**Arquivo:** `internal/domain/user/refresh_token.go`

```go
type RefreshToken struct {
    id        string
    userID    string
    tokenHash string    // hash SHA-256 do token (nao guardar raw)
    expiresAt time.Time
    revoked   bool
    createdAt time.Time
}
```

- `New(userID, tokenHash string, expiresAt time.Time) *RefreshToken`
- `Revoke()` — seta `revoked = true`
- `IsExpired() bool`
- `IsValid() bool` — `!revoked && !IsExpired()`

### 1.3 Model GORM — User

**Arquivo:** `internal/adapters/outbound/postgresql/model/user.go`

```go
type User struct {
    ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name         string    `gorm:"not null"`
    Email        string    `gorm:"uniqueIndex;not null"`
    PasswordHash string    `gorm:"not null"`
    Role         string    `gorm:"not null;default:'client'"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

- `FromUser(u *user.User) *User`
- `(m *User) ToDomain() *user.User`

### 1.4 Model GORM — RefreshToken

**Arquivo:** `internal/adapters/outbound/postgresql/model/refresh_token.go`

```go
type RefreshToken struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID    string    `gorm:"index;not null"`
    TokenHash string    `gorm:"uniqueIndex;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    Revoked   bool      `gorm:"default:false"`
    CreatedAt time.Time
}
```

### 1.5 Ports (interfaces)

**Arquivo:** `internal/ports/auth.go`

```go
// Repositories
type UserRepository interface {
    Create(ctx context.Context, u *user.User) error
    FindByEmail(ctx context.Context, email string) (*user.User, error)
    FindByID(ctx context.Context, id string) (*user.User, error)
}

type RefreshTokenRepository interface {
    Create(ctx context.Context, rt *user.RefreshToken) error
    FindByTokenHash(ctx context.Context, hash string) (*user.RefreshToken, error)
    RevokeByUserID(ctx context.Context, userID string) error
    Revoke(ctx context.Context, id string) error
}

// Use Cases
type RegisterUseCase interface {
    Execute(ctx context.Context, u *user.User) error
}

type LoginUseCase interface {
    Execute(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
}

type RefreshTokenUseCase interface {
    Execute(ctx context.Context, rawRefreshToken string) (accessToken, newRefreshToken string, err error)
}

type LogoutUseCase interface {
    Execute(ctx context.Context, userID string) error
}
```

### 1.6 Repositories PostgreSQL

**Arquivo:** `internal/adapters/outbound/postgresql/user_repository.go`

Metodos:
- `Create` — insert + SetID
- `FindByEmail` — `WHERE email = ?`
- `FindByID` — `WHERE id = ?`

**Arquivo:** `internal/adapters/outbound/postgresql/refresh_token_repository.go`

Metodos:
- `Create` — insert
- `FindByTokenHash` — `WHERE token_hash = ? AND revoked = false`
- `RevokeByUserID` — `UPDATE SET revoked = true WHERE user_id = ?` (invalida todos os tokens do user)
- `Revoke` — `UPDATE SET revoked = true WHERE id = ?`

---

## FASE 2 — Use Cases (logica de negocio)

### 2.1 Register

**Arquivo:** `internal/application/usecase/auth/register.go`

```
Recebe: name, email, senha raw, role
1. Verificar se email ja existe (repo.FindByEmail)
2. Hash da senha com bcrypt (cost 12)
3. Criar entidade User
4. Salvar (repo.Create)
```

Dependencias: `UserRepository`

### 2.2 Login

**Arquivo:** `internal/application/usecase/auth/login.go`

```
Recebe: email, senha raw
1. Buscar user por email
2. Comparar senha com bcrypt
3. Gerar access token JWT (15min):
   - Claims: sub (userID), email, role, exp, iat
   - Signing: HMAC-SHA256 com JWT_SECRET
4. Gerar refresh token:
   - Token: UUID v4 random
   - Salvar hash SHA-256 no banco (nunca guardar o raw)
   - Expiracao: 7 dias
5. Retornar access_token + refresh_token (raw)
```

Dependencias: `UserRepository`, `RefreshTokenRepository`, `Config`

### 2.3 Refresh

**Arquivo:** `internal/application/usecase/auth/refresh.go`

```
Recebe: refresh_token raw
1. Gerar hash SHA-256 do token recebido
2. Buscar no banco por hash
3. Validar: existe? nao revogado? nao expirado?
4. Revogar o refresh token usado (rotation — single use)
5. Buscar user por ID
6. Gerar novo access token
7. Gerar novo refresh token (salvar hash no banco)
8. Retornar novo par
```

Dependencias: `UserRepository`, `RefreshTokenRepository`, `Config`

Isso e **refresh token rotation**: cada refresh token so funciona UMA vez. Se alguem reusar um token ja usado, significa que vazou — e voce pode invalidar todos os tokens do user como medida de seguranca.

### 2.4 Logout

**Arquivo:** `internal/application/usecase/auth/logout.go`

```
Recebe: userID (do access token via context)
1. Revogar TODOS os refresh tokens do user
```

Dependencias: `RefreshTokenRepository`

---

## FASE 3 — Handler HTTP + DTOs

### 3.1 Commands (DTOs)

**Arquivo:** `internal/adapters/inbound/http/commands/auth.go`

```go
type RegisterRequest struct {
    Name     string `json:"name"     binding:"required"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"`     // segundos ate expirar
    TokenType    string `json:"token_type"`     // "Bearer"
}

type MessageResponse struct {
    Message string `json:"message"`
}
```

### 3.2 Auth Handler

**Arquivo:** `internal/adapters/inbound/http/auth/handler.go`

```go
type AuthHandler struct {
    register ports.RegisterUseCase
    login    ports.LoginUseCase
    refresh  ports.RefreshTokenUseCase
    logout   ports.LogoutUseCase
}

func (h *AuthHandler) SetupRoutes(rg *gin.RouterGroup) {
    auth := rg.Group("/auth")
    auth.POST("/register", h.Register)
    auth.POST("/login", h.Login)
    auth.POST("/refresh", h.Refresh)
    auth.POST("/logout", middleware.Auth(secret), h.Logout)  // unica rota auth que precisa de JWT
}
```

**Arquivos separados** (seguir padrao existente):
- `internal/adapters/inbound/http/auth/register.go`
- `internal/adapters/inbound/http/auth/login.go`
- `internal/adapters/inbound/http/auth/refresh.go`
- `internal/adapters/inbound/http/auth/logout.go`

---

## FASE 4 — Middlewares

### 4.1 Refatorar auth.go existente

**Arquivo:** `cmd/api/middleware/auth.go`

Alem de validar o token, extrair dados uteis pro contexto:

```go
// Antes (hoje)
c.Set("claims", claims)

// Depois
c.Set("user_id", claims["sub"])
c.Set("user_role", claims["role"])
c.Set("user_email", claims["email"])
```

### 4.2 RequireRole (novo)

**Arquivo:** `cmd/api/middleware/role.go`

```go
func RequireRole(allowed ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists {
            c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
            return
        }
        for _, a := range allowed {
            if role == a {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(403, gin.H{"error": "insufficient permissions"})
    }
}
```

### 4.3 Rate Limiter (novo)

**Arquivo:** `cmd/api/middleware/rate_limit.go`

Rate limiter in-memory por IP, sem dependencia externa:

```go
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc
```

Internamente usa `sync.Map` com IP como chave e um contador + timestamp.
Retorna `429 Too Many Requests` quando excede.

Aplicar APENAS nas rotas de auth:
```go
auth.POST("/login", middleware.RateLimit(5, time.Minute), h.Login)
auth.POST("/register", middleware.RateLimit(3, time.Minute), h.Register)
```

---

## FASE 5 — Wiring (DI + Rotas)

### 5.1 Config

**Arquivo:** `config/config.go`

Adicionar:
```go
AccessTokenExpMin   int    // env ACCESS_TOKEN_EXP_MIN, default 15
RefreshTokenExpDays int    // env REFRESH_TOKEN_EXP_DAYS, default 7
BcryptCost          int    // env BCRYPT_COST, default 12
```

### 5.2 .env.example

Adicionar:
```
ACCESS_TOKEN_EXP_MIN=15
REFRESH_TOKEN_EXP_DAYS=7
BCRYPT_COST=12
```

### 5.3 Container (bootstrap/container.go)

Adicionar na sequencia existente:

```go
// setupDatabase — adicionar AutoMigrate:
pgmodel.User{}
pgmodel.RefreshToken{}

// setupRepositories — adicionar:
c.UserRepo = postgresql.NewUserRepository(c.db)
c.RefreshTokenRepo = postgresql.NewRefreshTokenRepository(c.db)

// setupUseCases — adicionar:
c.RegisterUseCase = auth.NewRegister(c.UserRepo)
c.LoginUseCase = auth.NewLogin(c.UserRepo, c.RefreshTokenRepo, c.Config)
c.RefreshTokenUseCase = auth.NewRefresh(c.UserRepo, c.RefreshTokenRepo, c.Config)
c.LogoutUseCase = auth.NewLogout(c.RefreshTokenRepo)

// setupHandlers — adicionar:
c.AuthHandler = authhandler.NewAuthHandler(c.RegisterUseCase, c.LoginUseCase, c.RefreshTokenUseCase, c.LogoutUseCase)
```

### 5.4 Routes (routes/routes.go)

```go
// ANTES do grupo protected:
c.AuthHandler.SetupRoutes(v1)  // rotas publicas em /api/v1/auth/*

// Grupo protected (ja existe):
protected := v1.Group("/")
protected.Use(middleware.Auth(c.Config.JWTSecret))

// Admin-only (novo):
admin := protected.Group("/")
admin.Use(middleware.RequireRole("admin"))

// Customer + Vehicle + Service + Part handlers no admin
// ServiceOrder read no protected (client pode consultar as dele)
```

Organizacao de rotas por role:

```
/api/v1/auth/*                          → publico (com rate limit)

/api/v1/customers/*                     → admin
/api/v1/vehicles/*                      → admin
/api/v1/services/*                      → admin
/api/v1/parts/*                         → admin
/api/v1/service-orders (POST/PUT/DELETE) → admin
/api/v1/service-orders (GET)            → admin + client (client filtra por customerID)
```

---

## FASE 6 — Testes

### 6.1 Unitarios (dominios criticos)

**Arquivos:**
- `internal/domain/user/entity_test.go` — criacao, reconstitute, getters
- `internal/domain/user/refresh_token_test.go` — Revoke, IsExpired, IsValid
- `internal/application/usecase/auth/login_test.go` — senha errada, user inexistente, sucesso
- `internal/application/usecase/auth/refresh_test.go` — token expirado, revogado, rotation
- `internal/application/usecase/auth/register_test.go` — email duplicado, sucesso
- `cmd/api/middleware/role_test.go` — role permitida, proibida, ausente
- `cmd/api/middleware/rate_limit_test.go` — dentro do limite, excedido

### 6.2 Integracao

**Arquivo:** `internal/adapters/outbound/postgresql/user_repository_test.go`
- Create + FindByEmail
- Email duplicado retorna erro
- FindByID inexistente

---

## Arvore de arquivos novos

```
config/config.go                                          (modificar)
cmd/api/bootstrap/container.go                            (modificar)
cmd/api/routes/routes.go                                  (modificar)
cmd/api/middleware/auth.go                                 (modificar)
cmd/api/middleware/role.go                                 (criar)
cmd/api/middleware/rate_limit.go                           (criar)
internal/domain/user/entity.go                            (criar)
internal/domain/user/refresh_token.go                     (criar)
internal/ports/auth.go                                    (criar)
internal/adapters/outbound/postgresql/model/user.go       (criar)
internal/adapters/outbound/postgresql/model/refresh_token.go (criar)
internal/adapters/outbound/postgresql/user_repository.go  (criar)
internal/adapters/outbound/postgresql/refresh_token_repository.go (criar)
internal/application/usecase/auth/register.go             (criar)
internal/application/usecase/auth/login.go                (criar)
internal/application/usecase/auth/refresh.go              (criar)
internal/application/usecase/auth/logout.go               (criar)
internal/adapters/inbound/http/commands/auth.go           (criar)
internal/adapters/inbound/http/auth/handler.go            (criar)
internal/adapters/inbound/http/auth/register.go           (criar)
internal/adapters/inbound/http/auth/login.go              (criar)
internal/adapters/inbound/http/auth/refresh.go            (criar)
internal/adapters/inbound/http/auth/logout.go             (criar)
.env.example                                              (modificar)
```

---

## Ordem de implementacao

```
FASE 1  →  Domain + Models + Repos + Ports     (fundacao)
FASE 2  →  Use Cases                            (logica)
FASE 3  →  Handler + DTOs                       (HTTP)
FASE 4  →  Middlewares (role + rate limit)       (seguranca)
FASE 5  →  Wiring no container + rotas          (conectar tudo)
FASE 6  →  Testes                               (cobertura)
```

Cada fase compila e pode ser testada isoladamente.
