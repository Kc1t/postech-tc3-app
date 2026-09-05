# ADR-0007: Comunicação entre API Gateway e aplicação

- **Status:** Accepted
- **Data:** 2026-09
- **Decisão por:** Reservoir Devs

## Contexto

A Fase 3 exige um API Gateway na frente da aplicação, com rotas sensíveis protegidas por autenticação via CPF. A aplicação roda em pods no EKS; o gateway é um recurso gerenciado fora do cluster. Era preciso decidir **como** o gateway alcança os pods e **onde** a decisão de autorização acontece.

Duas perguntas independentes:

1. Qual o caminho de rede entre gateway e aplicação?
2. Quem valida o token — o gateway, a aplicação, ou os dois?

## Decisão

**Rede:** integração `HTTP_PROXY` do API Gateway v2 apontando para o **NLB público** provisionado pelo `Service` do tipo `LoadBalancer` da aplicação.

**Autorização:** validação em **duas camadas**. Um authorizer Lambda do tipo `REQUEST` valida a assinatura HS256 no gateway, e o middleware `Auth` da aplicação valida o mesmo token de novo. As duas camadas compartilham o `JWT_SECRET`.

O authorizer roda com `authorizer_result_ttl_in_seconds = 0`. O cache padrão de 300 s faria com que desativar um cliente só surtisse efeito cinco minutos depois — inaceitável para uma decisão de acesso.

## Alternativas consideradas

1. **VPC Link + NLB interno.** Fecharia o acesso direto à aplicação, deixando o gateway como única porta de entrada. Descartada pelo custo em tempo: exige NLB interno provisionado fora do ciclo de vida do `Service` do Kubernetes, criando dependência circular entre o repositório de infraestrutura e o da aplicação. Continua sendo o passo natural de evolução.

2. **JWT authorizer nativo do API Gateway.** Seria a opção mais barata operacionalmente — zero código, zero Lambda. **Tecnicamente inviável aqui:** o JWT authorizer de HTTP API exige um emissor OIDC com endpoint JWKS público para buscar as chaves. Nosso token é HS256, assinado com segredo compartilhado, sem JWKS. Migrar para RS256 com JWKS hospedado seria possível, mas exigiria publicar e rotacionar chaves públicas — desproporcional para o escopo.

3. **Só o gateway valida.** Deixaria a aplicação aceitar qualquer requisição que chegue ao NLB, que é público. Descartada: transformaria a exposição de rede em falha de autenticação.

4. **Só a aplicação valida.** Funcionaria, mas desperdiçaria o gateway como ponto de controle e faria a aplicação gastar recursos processando requisições que poderiam ser barradas antes. Também descumpre o espírito do requisito, que pede o gateway protegendo as rotas.

5. **Kong no cluster.** É o gateway ensinado com mais profundidade na fase. Descartado por custo operacional: exigiria pods, banco próprio para configuração e um Ingress adicional, sem ganho sobre a integração nativa do API Gateway com Lambda. Registrado na RFC-0003.

## Consequências

### Positivas

- Defesa em profundidade: comprometer o gateway não basta para chamar a aplicação, e alcançar o NLB diretamente não basta para passar do middleware.
- Zero infraestrutura adicional para operar — o gateway é totalmente gerenciado.
- O authorizer devolve `subject`, `role` e `document` no contexto da autorização, disponíveis para log e auditoria no gateway.
- As variáveis de integração são opcionais no Terraform, o que permite subir o gateway antes da Lambda e da aplicação existirem.

### Negativas

- **A aplicação continua alcançável direto pelo NLB**, contornando o gateway. Mitigado pela segunda camada de validação, mas significa que rate limiting e access log do gateway podem ser burlados.
- O segredo HS256 vive em três lugares (aplicação, emissor, authorizer). Rotacioná-lo exige coordenação entre dois repositórios.
- TTL de cache zero significa uma invocação de Lambda por requisição protegida. Com o free tier de 1 milhão de invocações/mês, é irrelevante no volume do projeto; em produção real seria um custo a reavaliar.

## Referências

- AWS — [Controlling access to HTTP APIs with Lambda authorizers](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-lambda-authorizer.html)
- Implementação: [`postech-tc3-infra-k8s/api_gateway.tf`](https://github.com/Kc1t/postech-tc3-infra-k8s/blob/main/api_gateway.tf) e [`postech-tc3-lambda-auth/cmd/authorizer`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/cmd/authorizer)
