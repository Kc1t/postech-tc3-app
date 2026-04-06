# TODO

## Refatorar Use Cases para um arquivo por operacao

Atualmente os use cases estao agrupados por entidade (ex: `customer_usecase.go` com Create, GetByID, GetAll, Update, Delete).

A ideia e ter **um arquivo por operacao**, seguindo o principio de responsabilidade unica (SRP).
Cada use case e uma struct com um unico metodo `Execute`.

### Estrutura proposta

```
internal/application/usecase/
  customer/
    create.go
    get_by_id.go
    get_by_document.go
    list.go
    update.go
    delete.go
  vehicle/
    create.go
    get_by_id.go
    list.go
    list_by_customer.go
    update.go
    delete.go
  service_order/
    create.go
    get_by_id.go
    list.go
    list_by_customer.go
    update_status.go
    update.go
    delete.go
  service/
    create.go
    get_by_id.go
    list.go
    update.go
    delete.go
  part/
    create.go
    get_by_id.go
    list.go
    update.go
    delete.go
    adjust_stock.go
```

### Exemplo de implementacao

```go
// internal/application/usecase/customer/create.go
package customeruc

type CreateCustomer struct {
    repo ports.CustomerRepository
}

func NewCreateCustomer(repo ports.CustomerRepository) *CreateCustomer {
    return &CreateCustomer{repo: repo}
}

func (uc *CreateCustomer) Execute(ctx context.Context, c *customer.Customer) error {
    // validar CPF/CNPJ
    // verificar duplicidade por documento
    return uc.repo.Create(ctx, c)
}
```

### Impacto nos ports

A interface `CustomerUseCase` em `ports/usecases.go` pode ser mantida como agregador
(o container injeta cada use case individualmente), ou quebrada em interfaces atomicas:

```go
// ports/usecases.go
type CreateCustomerUseCase interface {
    Execute(ctx context.Context, c *customer.Customer) error
}

type GetCustomerUseCase interface {
    Execute(ctx context.Context, id string) (*customer.Customer, error)
}
// ...
```

O handler recebe apenas as interfaces que precisa — sem depender de tudo.

```go
type CustomerHandler struct {
    create    ports.CreateCustomerUseCase
    getByID   ports.GetCustomerUseCase
    listAll   ports.ListCustomersUseCase
    update    ports.UpdateCustomerUseCase
    delete    ports.DeleteCustomerUseCase
}
```

### Vantagens

- Cada arquivo tem uma unica razao para mudar
- Testes unitarios mais simples e focados
- Facil identificar o que cada use case precisa (deps minimas)
- Novos devs encontram a logica pelo nome do arquivo

### Dependencias a resolver antes

- [ ] Definir convencao de nome do package (`customeruc`, `customer`, etc.)
- [ ] Atualizar `ports/usecases.go` com interfaces atomicas ou manter agregada
- [ ] Atualizar `bootstrap/container.go` para injetar cada UC separadamente
- [ ] Atualizar handlers para receber interfaces atomicas

---

## Unificar entidades de dominio em /entities

Atualmente cada entidade tem sua propria pasta dentro de `internal/domain/`:

```
internal/domain/
  customer/entity.go      (package customer)
  vehicle/entity.go       (package vehicle)
  service_order/entity.go  (package serviceorder)
  service/entity.go       (package service)
  part/entity.go          (package part)
```

Isso cria um package por entidade, o que gera imports verbosos e dificulta
referencias cruzadas entre entidades (ex: ServiceOrder referenciando Customer).

### Estrutura proposta

```
internal/entities/
  customer.go       (package entities)
  vehicle.go        (package entities)
  service_order.go  (package entities)
  service.go        (package entities)
  part.go           (package entities)
```

Todas as entidades no mesmo package `entities`, eliminando imports cruzados
e simplificando referencias entre elas.

### Exemplo de impacto

```go
// antes — ServiceOrder precisava importar packages separados
import (
    "github.com/fiap/postech-tc1/internal/domain/customer"
    "github.com/fiap/postech-tc1/internal/domain/vehicle"
)

// depois — tudo em entities
import "github.com/fiap/postech-tc1/internal/entities"

so := entities.ServiceOrder{...}
c  := entities.Customer{...}
```

### Impacto nos ports e adapters

- `ports/repositories.go` e `ports/usecases.go` passam a importar apenas `entities`
- Repositories e models MongoDB atualizam os imports
- HTTP models atualizam os imports

### Checklist

- [ ] Mover e renomear arquivos para `internal/domain/entities`
- [ ] Atualizar package de `customer`, `vehicle`, etc. para `entities`
- [ ] Atualizar todos os imports nos ports, adapters, usecases e handlers
- [ ] Remover pasta `internal/domain/` apos migracao


## Criar Commands/DTOs e Erros de Dominio

### Commands / DTOs

Atualmente os DTOs de request/response estao em `internal/adapters/inbound/http/model/`,
acoplados ao adapter HTTP. A ideia e mover as structs de entrada para a camada de dominio.

#### Nomenclatura a definir

| Opcao    | Pasta                        | Uso                           |
|----------|------------------------------|-------------------------------|
| Commands | `internal/domain/commands/`  | Enfatiza intencao (CQRS-like) |
| DTOs     | `internal/domain/dto/`       | Mais generico, familiar       |

#### Estrutura proposta

```
internal/domain/
  commands/
    customer_commands.go    — CreateCustomerCommand, UpdateCustomerCommand
    vehicle_commands.go
    serviceorder_commands.go
    service_commands.go
    part_commands.go
  responses/
    customer_response.go    — CustomerResponse, CustomerListResponse
    vehicle_response.go
    serviceorder_response.go
    service_response.go
    part_response.go
```

#### Exemplo

```go
// domain/commands/customer_commands.go
type CreateCustomerCommand struct {
    Name     string
    Document string
    Email    string
    Phone    string
}

// handler faz bind JSON → command, sem saber de dominio
func (h *CustomerHandler) Create(c *gin.Context) {
    var cmd commands.CreateCustomerCommand
    c.ShouldBindJSON(&cmd)
    h.create.Execute(ctx, cmd)
}

// usecase opera sobre o command, sem saber de HTTP
func (uc *CreateCustomer) Execute(ctx context.Context, cmd commands.CreateCustomerCommand) error {
    entity := entities.NewCustomer(cmd.Name, cmd.Document, cmd.Email, cmd.Phone)
    return uc.repo.Create(ctx, entity)
}
```

### Erros de Dominio

```
internal/domain/
  errors/
    errors.go
```

```go
var (
    ErrNotFound            = errors.New("not found")
    ErrAlreadyExists       = errors.New("already exists")
    ErrInvalidDocument     = errors.New("invalid CPF/CNPJ")
    ErrInvalidPlate        = errors.New("invalid plate")
    ErrInvalidStatus       = errors.New("invalid status transition")
    ErrInsufficientStock   = errors.New("insufficient stock")
    ErrOrderNotCancellable = errors.New("order cannot be cancelled at current status")
)
```

Handler mapeia para HTTP via funcao centralizada:

```go
func handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domainerrors.ErrNotFound):
        c.JSON(http.StatusNotFound, ...)
    case errors.Is(err, domainerrors.ErrAlreadyExists):
        c.JSON(http.StatusConflict, ...)
    case errors.Is(err, domainerrors.ErrInvalidDocument):
        c.JSON(http.StatusUnprocessableEntity, ...)
    default:
        c.JSON(http.StatusInternalServerError, ...)
    }
}
```

### Checklist

- [ ] Decidir entre `commands/` ou `dto/`
- [ ] Criar structs de command por entidade
- [ ] Criar structs de response por entidade
- [ ] Criar `internal/domain/errors/errors.go` com erros sentinela
- [ ] Atualizar use cases para receber commands em vez de entidades cruas
- [ ] Adicionar `handleError` centralizado nos handlers
- [ ] Remover ou simplificar `internal/adapters/inbound/http/model/` apos migracao

---

## Refatorar Handlers para um arquivo por operacao

Atualmente os handlers estao agrupados por entidade (ex: `customer_handler.go` com Create, FindAll, FindByID, Update, Delete).

A ideia e ter **um arquivo por operacao**, seguindo o mesmo padrao proposto para os use cases (SRP).
O `SetupRoutes` fica num arquivo dedicado por entidade.

### Estrutura proposta

```
internal/adapters/inbound/http/
  customer/
    handler.go        — struct CustomerHandler + NewCustomerHandler + SetupRoutes
    create.go
    find_by_id.go
    find_by_document.go
    list.go
    update.go
    delete.go
  vehicle/
    handler.go
    create.go
    find_by_id.go
    list.go
    list_by_customer.go
    update.go
    delete.go
  service_order/
    handler.go
    create.go
    find_by_id.go
    list.go
    list_by_customer.go
    update_status.go
    update.go
    delete.go
  service/
    handler.go
    create.go
    find_by_id.go
    list.go
    update.go
    delete.go
  part/
    handler.go
    create.go
    find_by_id.go
    list.go
    update.go
    delete.go
    adjust_stock.go
```

### Exemplo de implementacao

```go
// internal/adapters/inbound/http/customer/handler.go
package customerhandler

import "github.com/gin-gonic/gin"

type CustomerHandler struct {
    create        ports.CreateCustomerUseCase
    getByID       ports.GetCustomerUseCase
    getByDocument ports.GetCustomerByDocumentUseCase
    listAll       ports.ListCustomersUseCase
    update        ports.UpdateCustomerUseCase
    delete        ports.DeleteCustomerUseCase
}

func NewCustomerHandler(
    create ports.CreateCustomerUseCase,
    getByID ports.GetCustomerUseCase,
    getByDocument ports.GetCustomerByDocumentUseCase,
    listAll ports.ListCustomersUseCase,
    update ports.UpdateCustomerUseCase,
    delete ports.DeleteCustomerUseCase,
) *CustomerHandler { ... }

func (h *CustomerHandler) SetupRoutes(rg *gin.RouterGroup) {
    g := rg.Group("/customers")
    g.POST("",            h.Create)
    g.GET("",             h.FindAll)
    g.GET("/:id",         h.FindByID)
    g.GET("/doc/:doc",    h.FindByDocument)
    g.PUT("/:id",         h.Update)
    g.DELETE("/:id",      h.Delete)
}

// internal/adapters/inbound/http/customer/create.go
func (h *CustomerHandler) Create(c *gin.Context) {
    // bind JSON, chamar h.create.Execute(...)
}
```

### SetupRoutes centralizado

```go
// routes/routes.go — so orquestra
func Setup(router *gin.Engine, c *bootstrap.Container) {
    v1 := router.Group("/api/v1")
    protected := v1.Group("/").Use(middleware.Auth(...))

    c.CustomerHandler.SetupRoutes(protected)
    c.VehicleHandler.SetupRoutes(protected)
    c.ServiceOrderHandler.SetupRoutes(protected)
    c.ServiceHandler.SetupRoutes(protected)
    c.PartHandler.SetupRoutes(protected)
}
```

### Vantagens

- Cada arquivo tem uma unica razao para mudar (mesmo padrao dos use cases)
- Adicionar um novo handler nao requer alterar `routes.go`
- Facil localizar o handler pelo nome do arquivo

### Checklist

- [ ] Criar subpasta por entidade em `internal/adapters/inbound/http/`
- [ ] Mover struct + construtor + SetupRoutes para `handler.go` de cada entidade
- [ ] Separar cada metodo HTTP em seu proprio arquivo
- [ ] Atualizar imports no `container.go` e `routes.go`
- [ ] Remover arquivos `*_handler.go` da pasta raiz de http apos migracao

---

## Outros TODOs

- [ ] Implementar handlers (atualmente retornam 501)
- [ ] Validacao de CPF/CNPJ no dominio
- [ ] Validacao de placa no dominio
- [ ] Maquina de estados para status da OS
- [ ] Autenticacao JWT (endpoint de login)
- [ ] Paginacao nos endpoints de listagem
- [ ] Testes unitarios (cobertura minima 80%)
- [ ] Testes de integracao
- [ ] Scan de vulnerabilidades (requisito do tech challenge)
- [ ] Documentacao DDD (Event Storming no Miro)


---


### DÉBITOS TÉCNICOS


## Estrutura de Injecao de Dependencia (pkg/ioc)

Atualmente o `pkg/ioc` esta vazio (reservado). A ideia e implementar um container
proprio e minimalista que sirva de base para o wiring da aplicacao.

### Opcoes em avaliacao

#### Opcao A — Manter constructor injection explicito (atual)
O `container.go` faz o wiring manual em `setup*`. Simples, sem magia, facil de testar.
Nenhuma mudanca necessaria no `pkg/ioc`.

#### Opcao B — Fluent Builder sobre golobby/container
Adicionar `github.com/golobby/container/v3` e expor apenas um Builder fluente:

```go
// pkg/ioc/builder.go
ioc.NewBuilder().
    Singleton(func() ports.CustomerRepository { ... }).
    Singleton(func() ports.VehicleRepository { ... })

// struct injection via tag
type customerUseCase struct {
    Repo ports.CustomerRepository `container:"type"`
}
ioc.Fill(uc) // golobby resolve por tipo
```

#### Opcao C — IoC proprio com generics (sem dependencia externa)
Implementar `Register[T]`, `Resolve[T]` e `Inject` usando `sync.Map` + `reflect`:

```go
ioc.Register[ports.CustomerRepository](func() ports.CustomerRepository { ... })
repo := ioc.Resolve[ports.CustomerRepository]()

// ou fluente
ioc.NewBuilder().
    Singleton(func() ports.CustomerRepository { ... }).
    Singleton(func() ports.VehicleRepository { ... })

ioc.Inject(uc) // resolve campos com tag `container:"type"` por reflect
```

### Decisoes pendentes

- [ ] Escolher entre opcao A, B ou C
- [ ] Se B ou C: atualizar `container.go` para usar o novo pkg/ioc
- [ ] Se B ou C: adicionar `Reset()` para limpar o registry nos testes
- [ ] Avaliar se `pkg/ioc` deve expor `Singleton`, `Resolve`, `Fill` diretamente
      ou apenas o `Builder` fluente

### Convencao de tags

Independente da opcao escolhida, manter consistencia nas tags:

| Layer       | Tag                      |
|-------------|--------------------------|
| Repository  | `container:"repository"` |
| Use Case    | `container:"usecase"`    |
| Handler     | `container:"handler"`    |

---


