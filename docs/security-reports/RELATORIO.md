# Relatório de Análise de Vulnerabilidades

Workshop API — Tech Challenge Fase 1
Data da execução: 2026-04-27

## 1. Metodologia

Foram executados três scanners complementares, todos via container Docker para garantir reprodutibilidade:

| Ferramenta | Finalidade | Imagem |
|---|---|---|
| `gosec` | SAST estático para código Go | `securego/gosec:latest` |
| `govulncheck` | Vulnerabilidades em dependências com análise de call-graph | `golang:1.25` + `golang.org/x/vuln@v1.3.0` |
| `trivy` | Scan de `go.mod`, Dockerfile (misconfig) e segredos | `aquasec/trivy:latest` |

### Comandos executados

```bash
docker run --rm -v "$(pwd):/app" -w /app securego/gosec:latest -fmt=text -severity=low ./...
docker run --rm -v "$(pwd):/app" -w /app golang:1.25 sh -c "go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./..."
docker run --rm -v "$(pwd):/app" aquasec/trivy:latest fs --scanners vuln,secret,misconfig /app
```

Saídas brutas das execuções estão em `gosec.txt`, `govulncheck.txt` e `trivy.txt` neste mesmo diretório.

## 2. Resumo Consolidado

| Severidade | Qtd. | Origem |
|---|---|---|
| Crítica | 1 | Dependência (`pgx/v5`) |
| Alta | 1 | Dockerfile (sem `USER` não-root) |
| Média | 2 | Dependência (`golang.org/x/crypto`) |
| Baixa | 3 | Código (`gosec`), Dockerfile (sem `HEALTHCHECK`), Dependência |
| Segredos expostos | 0 | — |

## 3. Detalhamento

### 3.1 gosec (SAST) — 1 issue BAIXA

| ID | Arquivo | Descrição |
|---|---|---|
| G706 (CWE-117) | `cmd/api/bootstrap/container.go:284` | Possível log injection em `log.Printf("bootstrap: admin user seeded (%s)", email)`. A string vem de `ADMIN_EMAIL` (variável de ambiente controlada pelo operador), não de input externo. Risco aceito; para defesa em profundidade pode-se sanitizar quebras de linha antes do log. |

### 3.2 govulncheck — 0 vulnerabilidades exploráveis

```
Your code is affected by 0 vulnerabilities.
This scan also found 4 vulnerabilities in packages you import and 3
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
```

A análise de call-graph confirma que nenhuma CVE de dependência é alcançável a partir dos fluxos da aplicação.

### 3.3 trivy — dependências (`go.mod`)

| Biblioteca | CVE / Advisory | Severidade | Versão atual | Versão corrigida | Descrição |
|---|---|---|---|---|---|
| `github.com/jackc/pgx/v5` | CVE-2026-33816 | **CRÍTICA** | v5.6.0 | v5.9.0 | Memory-safety vulnerability |
| `github.com/jackc/pgx/v5` | GHSA-j88v-2chj-qfwx | BAIXA | v5.6.0 | v5.9.2 | SQL Injection via placeholder confusion com strings dollar-quoted |
| `golang.org/x/crypto` | CVE-2025-47914 | MÉDIA | v0.41.0 | v0.45.0 | DoS em SSH agent server por mensagens malformadas |
| `golang.org/x/crypto` | CVE-2025-58181 | MÉDIA | v0.41.0 | v0.45.0 | DoS em SSH GSSAPI por consumo de memória ilimitado |

Embora a severidade nominal seja CRÍTICA/MÉDIA, o `govulncheck` (call-graph) confirma que o código da aplicação **não invoca as funções vulneráveis**:

- As CVEs de `x/crypto` afetam SSH agent/GSSAPI, não usados pela aplicação (só bcrypt/JWT do mesmo módulo);
- A CVE crítica de `pgx` afeta caminhos de código que o GORM/aplicação não exercita.

Mesmo assim, recomenda-se atualizar como boa prática.

### 3.4 trivy — Dockerfile

| ID | Severidade | Descrição |
|---|---|---|
| DS-0002 | **ALTA** | Falta `USER` não-root — container roda como `root`, aumentando risco em caso de container escape. |
| DS-0026 | BAIXA | Falta `HEALTHCHECK` no Dockerfile. |

### 3.5 Segredos expostos

Nenhum segredo encontrado no repositório.

## 4. Mitigações Recomendadas

### 4.1 Atualizar dependências

```bash
go get github.com/jackc/pgx/v5@v5.9.2
go get golang.org/x/crypto@v0.45.0
go mod tidy
```

### 4.2 Hardening do Dockerfile

```dockerfile
RUN adduser -D -u 10001 appuser
USER appuser

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget -qO- http://localhost:8080/health || exit 1
```

### 4.3 Defesa em profundidade no log de bootstrap

Sanitizar quebras de linha em `cmd/api/bootstrap/container.go:284` antes de logar `email`.

## 5. Conclusão

Nível de risco residual **baixo**:

- 0 vulnerabilidades exploráveis no código (confirmado por call-graph com `govulncheck`);
- CVEs encontradas em dependências, embora algumas com severidade nominal alta, não são alcançáveis pelos fluxos da aplicação;
- Única issue do código é BAIXA e em contexto controlado (seed administrativo);
- Achados do Dockerfile são misconfigurações de hardening, não vulnerabilidades ativas;
- Nenhum segredo exposto.
