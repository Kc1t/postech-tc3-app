#!/usr/bin/env bash
set -euo pipefail

DEFS="${1:-$(dirname "$0")/../docs/newrelic/alerts.json}"
API="${NEW_RELIC_API_URL:-https://api.newrelic.com/graphql}"

for var in NEW_RELIC_API_KEY NEW_RELIC_ACCOUNT_ID; do
  if [[ -z "${!var:-}" ]]; then
    cat >&2 <<EOF
Falta a variavel $var.

A chave e do tipo User (prefixo NRAK), criada em
https://one.newrelic.com/api-keys — a license key nao serve aqui.

  export NEW_RELIC_API_KEY=NRAK-...
  export NEW_RELIC_ACCOUNT_ID=1234567

Se uma execucao anterior falhou no meio, a politica ja existe. Passe o id
dela em NEW_RELIC_POLICY_ID para nao criar uma duplicada.

Uso: $0 [caminho-do-alerts.json]
EOF
    exit 1
  fi
done

if [[ ! -f "$DEFS" ]]; then
  echo "Arquivo de definicoes nao encontrado: $DEFS" >&2
  exit 1
fi

nerdgraph() {
  local response
  response="$(curl -sS -X POST "$API" \
    -H "Api-Key: $NEW_RELIC_API_KEY" \
    -H 'Content-Type: application/json' \
    --data-binary @-)"

  if echo "$response" | jq -e '.errors' >/dev/null 2>&1; then
    echo "NerdGraph recusou a chamada:" >&2
    echo "$response" | jq -r '.errors[].message' >&2
    exit 1
  fi
  echo "$response"
}

POLICY_NAME="$(jq -r '.policy.name' "$DEFS")"
PREFERENCE="$(jq -r '.policy.incidentPreference' "$DEFS")"

if [[ -n "${NEW_RELIC_POLICY_ID:-}" ]]; then
  policy_id="$NEW_RELIC_POLICY_ID"
  echo "Reusando politica existente (id $policy_id)"
else
  policy_id="$(jq -n \
    --arg account "$NEW_RELIC_ACCOUNT_ID" --arg name "$POLICY_NAME" --arg pref "$PREFERENCE" \
    '{query:"mutation($account:Int!,$name:String!,$pref:AlertsIncidentPreference!){alertsPolicyCreate(accountId:$account,policy:{name:$name,incidentPreference:$pref}){id}}",
      variables:{account:($account|tonumber),name:$name,pref:$pref}}' \
    | nerdgraph | jq -r '.data.alertsPolicyCreate.id')"

  echo "Politica criada: $POLICY_NAME (id $policy_id)"
fi

terms() {
  jq -c '[
    (if .critical then {priority:"CRITICAL", threshold:.critical.threshold, thresholdDuration:.critical.thresholdDuration, operator:.critical.operator, thresholdOccurrences:"ALL"} else empty end),
    (if .warning  then {priority:"WARNING",  threshold:.warning.threshold,  thresholdDuration:.warning.thresholdDuration,  operator:.warning.operator,  thresholdOccurrences:"ALL"} else empty end)
  ]' <<<"$1"
}

jq -c '.conditions[]' "$DEFS" | while read -r condition; do
  name="$(jq -r '.name' <<<"$condition")"

  jq -n \
    --arg account "$NEW_RELIC_ACCOUNT_ID" \
    --arg policy "$policy_id" \
    --arg name "$name" \
    --arg description "$(jq -r '.description' <<<"$condition")" \
    --arg nrql "$(jq -r '.nrql' <<<"$condition")" \
    --argjson terms "$(terms "$condition")" \
    '{query:"mutation($account:Int!,$policy:ID!,$condition:AlertsNrqlConditionStaticInput!){alertsNrqlConditionStaticCreate(accountId:$account,policyId:$policy,condition:$condition){id}}",
      variables:{
        account:($account|tonumber),
        policy:$policy,
        condition:{
          name:$name, description:$description, enabled:true,
          nrql:{query:$nrql},
          signal:{aggregationWindow:60, aggregationMethod:"EVENT_FLOW", aggregationDelay:120},
          terms:$terms,
          violationTimeLimitSeconds:86400
        }
      }}' \
    | nerdgraph >/dev/null

  echo "  condicao criada: $name"
done

echo
echo "Pronto. Verifique em https://one.newrelic.com/alerts — a politica nao tem"
echo "canal de notificacao; adicione um workflow se quiser e-mail ou Slack."
