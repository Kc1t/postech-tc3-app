# Diagramas

Todos em Mermaid embutido em Markdown: renderizam direto no GitHub, ficam versionados junto com o código e aparecem no diff quando a arquitetura muda.

## Estrutura

| Documento | O que responde |
|---|---|
| [`componentes.md`](./componentes.md) | Como as peças se encaixam na nuvem — APIs, banco, monitoramento e fronteiras de repositório |
| [`roteamento.md`](./roteamento.md) | Por onde uma requisição passa, quem a barra e onde cada código de erro nasce |
| [`pipelines.md`](./pipelines.md) | O que cada um dos quatro pipelines faz, e quais secrets cada repositório precisa |

## Fluxos

| Documento | O que responde |
|---|---|
| [`sequencia-autenticacao.md`](./sequencia-autenticacao.md) | Como o CPF vira JWT e como esse token é aceito numa rota protegida |
| [`sequencia-ordem-servico.md`](./sequencia-ordem-servico.md) | Da abertura da OS até a aprovação do orçamento pelo cliente, com a máquina de estados |

## Dados

| Documento | O que responde |
|---|---|
| [`../MODELAGEM_DE_DADOS.md`](../MODELAGEM_DE_DADOS.md) | Diagrama ER, justificativa formal do banco e explicação dos relacionamentos |
| [`../database-model.dbml`](../database-model.dbml) | O mesmo modelo em DBML, para <https://dbdiagram.io/d> |
| [`../documentation-diagram.drawio`](../documentation-diagram.drawio) | Diagrama das Fases 1 e 2, mantido como referência histórica |

## Por onde começar

Quem está chegando agora: `componentes.md` para a visão geral, depois `roteamento.md` para entender as camadas de proteção, e `sequencia-autenticacao.md` para o fluxo que é o coração da Fase 3.
