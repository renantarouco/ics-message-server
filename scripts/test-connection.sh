#!/usr/bin/env sh
TOKEN=$(curl --request POST \
  --url http://localhost:7000/auth \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data nickname=test4 | jq -r ".token")

curl --include \
  --no-buffer \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer $TOKEN" \
  --header "Connection: Upgrade" \
  --header "Upgrade: websocket" \
  --header "Host: localhost:7000" \
  --header "Origin: http://example.com:80" \
  --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
  --header "Sec-WebSocket-Version: 13" \
  http://localhost:7000/ws
