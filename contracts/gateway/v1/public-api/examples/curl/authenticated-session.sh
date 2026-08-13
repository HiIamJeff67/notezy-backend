#!/usr/bin/env bash
set -euo pipefail

gateway_base_url="${GATEWAY_BASE_URL:-http://localhost/api/development/v1}"
account="${ACCOUNT:?set ACCOUNT}"
password="${PASSWORD:?set PASSWORD}"
cookie_jar="$(mktemp -t notezy-cookie-jar.XXXXXX)"
response_file="$(mktemp -t notezy-login-response.XXXXXX)"
trap 'rm -f "$cookie_jar" "$response_file"' EXIT
login_payload="$(jq -cn --arg account "$account" --arg password "$password" '{account:$account,password:$password}')"

curl --fail-with-body --silent --show-error \
  -c "$cookie_jar" \
  -H 'Content-Type: application/json' \
  -H 'User-Agent: NotezyCurlSession/1.0' \
  --data "$login_payload" \
  "$gateway_base_url/auth/login" > "$response_file"

csrf_token="$(jq -er '.data.csrfToken' "$response_file")"

curl --fail-with-body --silent --show-error \
  -b "$cookie_jar" -c "$cookie_jar" \
  -H 'User-Agent: NotezyCurlSession/1.0' \
  -H "X-CSRF-Token: $csrf_token" \
  "$gateway_base_url/users/me"

# Keep the cookie jar private and short-lived. Read replacement CSRF values from
# X-CSRF-Token or refreshableTokens.newCSRFToken after access-token rotation.
