# ADR-0007: Erros sentinela centralizados em `domainerrors`

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

Em arquitetura hexagonal (ver [ADR-0001](./0001-hexagonal-architecture.md)), o domínio precisa expressar falhas de negócio sem vazar detalhes de infraestrutura. Por exemplo:

- O repositório recebe `gorm.ErrRecordNotFound` ao buscar um cliente inexistente. Esse erro **não pode** subir até o handler como está, senão o use case e o handler ficam acoplados ao GORM.
- Use cases precisam sinalizar invariantes violadas (estoque insuficiente, transição de status inválida, documento inválido) sem inventar `errors.New()` espalhados pelo código.
- Handlers HTTP precisam mapear esses erros para status codes apropriados (`404`, `409`, `422`).

Sem uma estratégia clara, o código tende a:

- Repetir literais de erro em vários lugares (`errors.New("not found")` em três use cases diferentes).
- Vazar `gorm.ErrRecordNotFound` para o handler (acoplamento).
- Mapear erros para status HTTP de forma inconsistente.

## Decisão

Centralizar **todos os erros de domínio** em um pacote único: `internal/domain/errors/` (pacote Go: `domainerrors`).

### Padrão

Erros sentinela exportados:

```go
package domainerrors

import "errors"

var (
    ErrNotFound                 = errors.New("não encontrado")
    ErrConflict                 = errors.New("conflito de dados")
    ErrInvalidDocument          = errors.New("documento inválido (CPF/CNPJ)")
    ErrInvalidPlate             = errors.New("placa de veículo inválida")
    ErrInvalidStatusTransition  = errors.New("transição de status inválida")
    ErrInsufficientStock        = errors.New("estoque insuficiente")
    ErrInvalidCredentials       = errors.New("credenciais inválidas")
    ErrAccountLocked            = errors.New("conta bloqueada por excesso de tentativas")
    ErrUnauthorized             = errors.New("não autorizado")
    // ... outros conforme necessário
)
```

### Regras

1. **Use cases retornam apenas erros de `domainerrors`** (ou erros wrappados com `fmt.Errorf("...: %w", domainerrors.ErrXxx)` para preservar contexto).
2. **Repositórios traduzem** erros de infraestrutura para erros de domínio:
   ```go
   if errors.Is(err, gorm.ErrRecordNotFound) {
       return nil, domainerrors.ErrNotFound
   }
   ```
3. **Handlers fazem switch com `errors.Is`** para mapear para HTTP:
   ```go
   switch {
   case errors.Is(err, domainerrors.ErrNotFound):
       c.JSON(http.StatusNotFound, ...)
   case errors.Is(err, domainerrors.ErrInsufficientStock):
       c.JSON(http.StatusUnprocessableEntity, ...)
   ...
   }
   ```
4. **Nada de `errors.New("...")` solto em use cases ou handlers.** Se for necessário um erro novo, ele vai para `domainerrors`.

## Alternativas consideradas

1. **Tipos de erro estruturados (structs com método `Error()`).** Permitem carregar dados estruturados, mas adicionam boilerplate. Para o escopo atual, sentinelas com `errors.Is` são suficientes.
2. **Pacote `errors` por bounded context** (ex.: `customer/errors`, `vehicle/errors`). Mais granular, mas leva à duplicação (`ErrNotFound` em cada). A alternativa centralizada simplifica imports.
3. **Retornar HTTP status diretamente do use case.** Acopla domínio à camada de transporte — viola a regra de dependência da arquitetura hexagonal.
4. **Uso de `panic`/`recover`.** Não-idiomático em Go para fluxos esperados de negócio.

## Consequências

### Positivas

- Mensagens de erro consistentes em todo o sistema.
- Mapeamento HTTP centralizável em uma função utilitária (`adapters/inbound/http/errors.go`).
- Use cases podem ser testados com `errors.Is(err, domainerrors.ErrXxx)` sem comparar strings.
- Adicionar um novo erro é uma mudança em um único arquivo.
- Cumpre a regra de não vazar erros de infraestrutura (ex.: `gorm.ErrRecordNotFound`) para fora do adapter de persistência.

### Negativas

- Erros sentinela não carregam contexto adicional por padrão (ex.: "qual ID não foi encontrado"). Mitigado por `fmt.Errorf("customer %s: %w", id, domainerrors.ErrNotFound)` quando necessário.
- Para evolução com mensagens i18n, será necessário substituir os sentinelas por structs com código de erro — refactor previsto se a aplicação crescer.

## Referências

- Dave Cheney, *Don't just check errors, handle them gracefully* (2016).
- `errors.Is` / `errors.As` — https://pkg.go.dev/errors
- Repositório: `internal/domain/errors/errors.go`, `internal/adapters/inbound/http/errors.go`.
