#!/usr/bin/env bash
set -euo pipefail

gateway_base_url="${GATEWAY_BASE_URL:-http://localhost/api/development/v1}"
cookie_jar="${COOKIE_JAR:-./notegic-cookies.txt}"
csrf_token="${CSRF_TOKEN:-}"
user_agent="${USER_AGENT:-NotegicCurlExample/1.0}"
account="${ACCOUNT:-}"
email="${EMAIL:-}"
password="${PASSWORD:-}"
id="${AVATAR_ID:-1}"
userPublicId="${USERPUBLICID:-00000000-0000-4000-8000-000000000001}"
stationId="${STATIONID:-00000000-0000-4000-8000-000000000001}"
routineId="${ROUTINEID:-00000000-0000-4000-8000-000000000001}"
routineTagId="${ROUTINETAGID:-00000000-0000-4000-8000-000000000001}"
routineTaskId="${ROUTINETASKID:-00000000-0000-4000-8000-000000000001}"
rootShelfId="${ROOTSHELFID:-00000000-0000-4000-8000-000000000001}"
subShelfId="${SUBSHELFID:-00000000-0000-4000-8000-000000000001}"
prevSubShelfId="${PREVSUBSHELFID:-00000000-0000-4000-8000-000000000001}"
parentSubShelfId="${PARENTSUBSHELFID:-00000000-0000-4000-8000-000000000001}"
materialId="${MATERIALID:-00000000-0000-4000-8000-000000000001}"
blockPackId="${BLOCKPACKID:-00000000-0000-4000-8000-000000000001}"
blockId="${BLOCKID:-00000000-0000-4000-8000-000000000001}"
itemId="${ITEMID:-00000000-0000-4000-8000-000000000001}"

# Run individual functions deliberately. DELETE/reset functions are not invoked automatically.

deleteMe() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authCode":"123456"}' \
    "$gateway_base_url/auth/delete-me"
}

forgetPassword() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    --data '{"account":"example","authCode":"123456","newPassword":"Example-Password-123!"}' \
    "$gateway_base_url/auth/forget-password"
}

login() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    --data "{\"account\":\"$account\",\"password\":\"$password\"}" \
    "$gateway_base_url/auth/login"
}

loginViaGoogle() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    --data '{"authorizationCode":"example"}' \
    "$gateway_base_url/auth/login-via-google"
}

logout() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/auth/logout"
}

register() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    --data "{\"name\":\"$account\",\"email\":\"$email\",\"password\":\"$password\"}" \
    "$gateway_base_url/auth/register"
}

registerViaGoogle() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    --data '{"authorizationCode":"example"}' \
    "$gateway_base_url/auth/register-via-google"
}

resetEmail() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authCode":"123456","newEmail":"developer@example.com"}' \
    "$gateway_base_url/auth/reset-email"
}

resetMe() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authCode":"123456"}' \
    "$gateway_base_url/auth/reset-me"
}

sendAuthCode() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    --data '{"email":"developer@example.com"}' \
    "$gateway_base_url/auth/send-auth-code"
}

validateEmail() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authCode":"123456"}' \
    "$gateway_base_url/auth/validate-email"
}

deleteMyBlockPacksByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"blockPackIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/block-packs/batch"
}

createBlockPacks() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"createdBlockPacks":[{"headerBackgroundURL":"https://example.com","icon":"😀","id":"00000000-0000-4000-8000-000000000001","name":"example","parentSubShelfId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/block-packs/batch"
}

updateMyBlockPacksByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"updatedBlockPacks":[{"blockPackId":"00000000-0000-4000-8000-000000000001","setNull":{},"values":{"headerBackgroundURL":"https://example.com","icon":"😀","name":"example"}}]}' \
    "$gateway_base_url/block-packs/batch"
}

moveMyBlockPacksByParentSubShelfIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"movedBlockPacks":[{"blockPackIds":["00000000-0000-4000-8000-000000000001"],"destinationParentSubShelfId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/block-packs/batch/position"
}

restoreMyBlockPacksByIds() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"blockPackIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/block-packs/batch/restore"
}

moveMyBlockPacksByParentSubShelfId() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"blockPackIds":["00000000-0000-4000-8000-000000000001"],"destinationParentSubShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/block-packs/position"
}

getAllMyBlockPacksByRootShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/block-packs/root-shelf/${rootShelfId}?areDeleted=true"
}

getMyBlockPacksByParentSubShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/block-packs/sub-shelf/${parentSubShelfId}?areDeleted=true"
}

createBlockPack() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"headerBackgroundURL":"https://example.com","icon":"😀","id":"00000000-0000-4000-8000-000000000001","name":"example","parentSubShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/block-packs/sub-shelf/${parentSubShelfId}"
}

deleteMyBlockPackById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/block-packs/${blockPackId}"
}

getMyBlockPackById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/block-packs/${blockPackId}?isDeleted=true"
}

updateMyBlockPackById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"headerBackgroundURL":"https://example.com","icon":"😀","name":"example"}}' \
    "$gateway_base_url/block-packs/${blockPackId}"
}

getMyBlockPackAndItsParentById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/block-packs/${blockPackId}/parent"
}

moveMyBlockPackByParentSubShelfId() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"blockPackId":"00000000-0000-4000-8000-000000000001","destinationParentSubShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/block-packs/${blockPackId}/position"
}

restoreMyBlockPackById() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/block-packs/${blockPackId}/restore"
}

getMyBlocksByIds() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/blocks/batch?blockIds=00000000-0000-4000-8000-000000000001"
}

getMyBlocksByBlockPackId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/blocks/block-pack/${blockPackId}"
}

getMyBlockById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/blocks/${blockId}"
}

graphQLGet() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/graphql"
}

graphQLPost() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"query":"query ContractCheck { __typename }","variables":{}}' \
    "$gateway_base_url/graphql"
}

deleteMyMaterialsByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"materialIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/materials/batch"
}

moveMyMaterialsByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"destinationParentSubShelfId":"00000000-0000-4000-8000-000000000001","materialIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/materials/batch/parent"
}

restoreMyMaterialsByIds() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"materialIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/materials/batch/restore"
}

getAllMyMaterialsByRootShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/materials/root-shelf/${rootShelfId}?areDeleted=true"
}

getMyMaterialsByParentSubShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/materials/sub-shelf/${parentSubShelfId}?areDeleted=true"
}

createMyMaterial() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"name":"example","parentSubShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/materials/sub-shelf/${parentSubShelfId}"
}

deleteMyMaterialById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/materials/${materialId}"
}

getMyMaterialById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/materials/${materialId}?isDeleted=true"
}

updateMyMaterialById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"name":"example"}}' \
    "$gateway_base_url/materials/${materialId}"
}

saveMyMaterialById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"contentFile":[1]}' \
    "$gateway_base_url/materials/${materialId}/content"
}

getMyMaterialAndItsParentById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/materials/${materialId}/parent"
}

moveMyMaterialById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"destinationParentSubShelfId":"00000000-0000-4000-8000-000000000001","materialId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/materials/${materialId}/parent"
}

restoreMyMaterialById() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/materials/${materialId}/restore"
}

getMyAccount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/me/account"
}

updateMyAccount() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authCode":"123456","setNull":{},"values":{"backupEmail":"developer@example.com","countryCode":"example","phoneNumber":"example"}}' \
    "$gateway_base_url/me/account"
}

unbindGoogleAccount() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authCode":"123456"}' \
    "$gateway_base_url/me/account/google"
}

bindGoogleAccount() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"authorizationCode":"example"}' \
    "$gateway_base_url/me/account/google"
}

getMyInfo() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/me/info"
}

updateMyInfo() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"avatarURL":"https://example.com","birthDate":"2026-01-01T00:00:00Z","country":"example","coverBackgroundURL":"https://example.com","gender":"example","header":"example","introduction":"example"}}' \
    "$gateway_base_url/me/info"
}

getMySetting() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/me/settings"
}

updateMySetting() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"density":"Balanced","language":"English","lineWrap":true,"privatePreviews":false,"quietMode":true,"quietModeEndMinute":480,"quietModeStartMinute":1320,"quickInsert":true,"reduceMotion":false,"routineNudges":true,"startSurface":"Dashboard","syncNotifications":true}}' \
    "$gateway_base_url/me/settings"
}

delete() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/notifications"
}

search() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/notifications"
}

markRead() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/notifications/read"
}

countUnread() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/notifications/unread-count"
}

createMyBlockPackChannelTicket() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"blockPackId":"00000000-0000-4000-8000-000000000001","permission":"read"}' \
    "$gateway_base_url/realtime/channel/block-pack/ticket"
}

createMyRealtimeConnectionTicket() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/realtime/connection/ticket"
}

createRootShelf() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"id":"00000000-0000-4000-8000-000000000001","name":"example"}' \
    "$gateway_base_url/root-shelves"
}

deleteMyRootShelvesByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"rootShelfIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/root-shelves/batch"
}

createRootShelves() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"insertedRootShelves":[{"id":"00000000-0000-4000-8000-000000000001","name":"example"}]}' \
    "$gateway_base_url/root-shelves/batch"
}

updateMyRootShelvesByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"updatedRootShelves":[{"rootShelfId":"00000000-0000-4000-8000-000000000001","setNull":{},"values":{"name":"example"}}]}' \
    "$gateway_base_url/root-shelves/batch"
}

restoreMyRootShelvesByIds() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"rootShelfIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/root-shelves/batch/restore"
}

leaveMyRootShelves() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"rootShelves":[{"rootShelfId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/root-shelves/memberships/me"
}

deleteMyRootShelfById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"rootShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/root-shelves/${rootShelfId}"
}

getMyRootShelfById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/root-shelves/${rootShelfId}?isDeleted=true"
}

updateMyRootShelfById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"name":"example"}}' \
    "$gateway_base_url/root-shelves/${rootShelfId}"
}

leaveMyRootShelf() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/root-shelves/${rootShelfId}/memberships/me"
}

transferMyRootShelfOwnership() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"targetUserPublicId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/root-shelves/${rootShelfId}/ownership"
}

deleteMyRootShelfPermissions() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"userPublicIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions"
}

upsertMyRootShelfPermissions() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"permissions":[{"permission":"Read","userPublicId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions"
}

deleteMyRootShelfPermission() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions/${userPublicId}"
}

getMyRootShelfPermission() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions/${userPublicId}"
}

updateMyRootShelfPermission() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions/${userPublicId}"
}

createMyRootShelfPermission() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"permission":"Read"}' \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions/${userPublicId}"
}

upsertMyRootShelfPermission() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/root-shelves/${rootShelfId}/permissions/${userPublicId}"
}

restoreMyRootShelfById() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"rootShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/root-shelves/${rootShelfId}/restore"
}

getAllMyRoutineTags() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tags?areDeleted=true"
}

createRoutineTag() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"color":"example","icon":"example","id":"00000000-0000-4000-8000-000000000001","name":"example"}' \
    "$gateway_base_url/routine-tags"
}

createRoutineTags() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"createdRoutineTags":[{"color":"example","icon":"example","id":"00000000-0000-4000-8000-000000000001","name":"example"}]}' \
    "$gateway_base_url/routine-tags/batch"
}

updateMyRoutineTagsByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"updatedRoutineTags":[{"routineTagId":"00000000-0000-4000-8000-000000000001","setNull":{},"values":{"color":"example","icon":"example","name":"example"}}]}' \
    "$gateway_base_url/routine-tags/batch"
}

hardDeleteMyRoutineTagsByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineTagIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/routine-tags/batch/permanently"
}

getMyRoutineTagById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tags/${routineTagId}?isDeleted=true"
}

updateMyRoutineTagById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"color":"example","icon":"example","name":"example"}}' \
    "$gateway_base_url/routine-tags/${routineTagId}"
}

hardDeleteMyRoutineTagById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/routine-tags/${routineTagId}/permanently"
}

getAllMyRoutineTaskRecordsByRoutineTaskId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-task-records/routine-task/${routineTaskId}?limit=1"
}

visualizeMyRoutineTaskRecordActualEndedAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-task-records/visualizations/actual-ended-at-count?permission=Read&queryRangeEndedAt=2026-01-01T00%3A00%3A00Z&queryRangeStartedAt=2026-01-01T00%3A00%3A00Z&routineTaskIds=00000000-0000-4000-8000-000000000001&timeHourUnit=1"
}

visualizeMyRoutineTaskRecordActualStartedAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-task-records/visualizations/actual-started-at-count?permission=Read&queryRangeEndedAt=2026-01-01T00%3A00%3A00Z&queryRangeStartedAt=2026-01-01T00%3A00%3A00Z&routineTaskIds=00000000-0000-4000-8000-000000000001&timeHourUnit=1"
}

visualizeMyRoutineTaskRecordPurposeCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-task-records/visualizations/purpose-count?permission=Read&routineTaskIds=00000000-0000-4000-8000-000000000001"
}

visualizeMyRoutineTaskRecordScheduledAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-task-records/visualizations/scheduled-at-count?permission=Read&queryRangeEndedAt=2026-01-01T00%3A00%3A00Z&queryRangeStartedAt=2026-01-01T00%3A00%3A00Z&routineTaskIds=00000000-0000-4000-8000-000000000001&timeHourUnit=1"
}

visualizeMyRoutineTaskRecordStatusCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-task-records/visualizations/status-count?permission=Read&routineTaskIds=00000000-0000-4000-8000-000000000001"
}

getAllMyRoutineTasks() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks?areDeleted=true"
}

hardDeleteMyRoutineTasksByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineTaskIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/routine-tasks/batch/permanently"
}

createRoutineTaskByRoutineId() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"maxAttempts":1,"nextScheduledAt":"2026-01-01T00:00:00Z","payload":{},"period":"Daily","priority":1,"purpose":"CreateRootShelf","routineId":"00000000-0000-4000-8000-000000000001","title":"example"}' \
    "$gateway_base_url/routine-tasks/routine/${routineId}"
}

getAllMyRoutineTasksByRoutineIds() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/routines?areDeleted=true&routineIds=00000000-0000-4000-8000-000000000001"
}

visualizeMyRoutineTaskActualEndedAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/visualizations/actual-ended-at-count"
}

visualizeMyRoutineTaskActualStartedAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/visualizations/actual-started-at-count"
}

visualizeMyRoutineTaskPurposeCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/visualizations/purpose-count"
}

visualizeMyRoutineTaskScheduledAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/visualizations/scheduled-at-count?permission=Read&queryRangeEndedAt=2026-01-01T00%3A00%3A00Z&queryRangeStartedAt=2026-01-01T00%3A00%3A00Z&timeHourUnit=1"
}

visualizeMyRoutineTaskStatusCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/visualizations/status-count?permission=Read"
}

getMyRoutineTaskById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routine-tasks/${routineTaskId}?isDeleted=true"
}

updateMyRoutineTaskById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineTaskId":"00000000-0000-4000-8000-000000000001","setNull":{},"values":{"maxAttempts":1,"nextScheduledAt":"2026-01-01T00:00:00Z","payload":{},"period":"Daily","priority":1,"purpose":"CreateRootShelf","routineId":"00000000-0000-4000-8000-000000000001","title":"example"}}' \
    "$gateway_base_url/routine-tasks/${routineTaskId}"
}

hardDeleteMyRoutineTaskById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineTaskId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routine-tasks/${routineTaskId}/permanently"
}

resumeMyRoutineTaskById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineTaskId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routine-tasks/${routineTaskId}/suspension"
}

pauseMyRoutineTaskById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineTaskId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routine-tasks/${routineTaskId}/suspension"
}

getAllMyRoutinesByTimeRange() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines?areDeleted=true&from=2026-01-01T00%3A00%3A00Z&stationIds=00000000-0000-4000-8000-000000000001&to=2026-01-01T00%3A00%3A00Z"
}

deleteMyRoutinesByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/routines/batch"
}

createRoutinesByStationIds() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"createdRoutines":[{"description":"example","id":"00000000-0000-4000-8000-000000000001","isPinned":true,"period":"Daily","scheduledEndAt":"2026-01-01T00:00:00Z","scheduledStartAt":"2026-01-01T00:00:00Z","stationId":"00000000-0000-4000-8000-000000000001","status":"Scheduled","timezone":"example","title":"example"}]}' \
    "$gateway_base_url/routines/batch"
}

updateMyRoutinesByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"updatedRoutines":[{"routineId":"00000000-0000-4000-8000-000000000001","setNull":{},"values":{"description":"example","isPinned":true,"period":"Daily","scheduledEndAt":"2026-01-01T00:00:00Z","scheduledStartAt":"2026-01-01T00:00:00Z","stationId":"00000000-0000-4000-8000-000000000001","status":"Scheduled","timezone":"example","title":"example"}}]}' \
    "$gateway_base_url/routines/batch"
}

hardDeleteMyRoutinesByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/routines/batch/permanently"
}

restoreMyRoutinesByIds() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/routines/batch/restore"
}

linkRoutineItemsByIds() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"isUnlink":true,"linkedRoutinesAndItems":[{"itemId":"00000000-0000-4000-8000-000000000001","itemType":"BlockPack","routineId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/routines/items"
}

getMyRoutinesByStationId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines/station/${stationId}?areDeleted=true"
}

createRoutineByStationId() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"description":"example","id":"00000000-0000-4000-8000-000000000001","isPinned":true,"period":"Daily","scheduledEndAt":"2026-01-01T00:00:00Z","scheduledStartAt":"2026-01-01T00:00:00Z","stationId":"00000000-0000-4000-8000-000000000001","status":"Scheduled","timezone":"example","title":"example"}' \
    "$gateway_base_url/routines/station/${stationId}"
}

linkRoutineTagsByIds() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"isUnlink":true,"linkedRoutinesAndTags":[{"routineId":"00000000-0000-4000-8000-000000000001","routineTagId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/routines/tags"
}

visualizeMyRoutinePeriodCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines/visualizations/period-count"
}

visualizeMyRoutineScheduledEndAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines/visualizations/scheduled-end-at-count"
}

visualizeMyRoutineScheduledStartAtCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines/visualizations/scheduled-start-at-count?permission=Read&queryRangeEndedAt=2026-01-01T00%3A00%3A00Z&queryRangeStartedAt=2026-01-01T00%3A00%3A00Z&timeHourUnit=1"
}

visualizeMyRoutineStatusCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines/visualizations/status-count?permission=Read"
}

deleteMyRoutineById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routines/${routineId}"
}

getMyRoutineById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/routines/${routineId}?isDeleted=true"
}

updateMyRoutineById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineId":"00000000-0000-4000-8000-000000000001","setNull":{},"values":{"description":"example","isPinned":true,"period":"Daily","scheduledEndAt":"2026-01-01T00:00:00Z","scheduledStartAt":"2026-01-01T00:00:00Z","stationId":"00000000-0000-4000-8000-000000000001","status":"Scheduled","timezone":"example","title":"example"}}' \
    "$gateway_base_url/routines/${routineId}"
}

linkRoutineItemById() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"isUnlink":true,"itemId":"00000000-0000-4000-8000-000000000001","itemType":"BlockPack","routineId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routines/${routineId}/items/${itemId}"
}

hardDeleteMyRoutineById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/routines/${routineId}/permanently"
}

restoreMyRoutineById() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"routineId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routines/${routineId}/restore"
}

linkRoutineTagById() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"isUnlink":true,"routineId":"00000000-0000-4000-8000-000000000001","routineTagId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/routines/${routineId}/tags/${routineTagId}"
}

getGlobalAvatar() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/static/global-images/avatars/${id}"
}

getAllMyStations() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/stations?areDeleted=true"
}

createStation() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"description":"example","headerBackgroundURL":"https://example.com","icon":"example","id":"00000000-0000-4000-8000-000000000001","name":"example"}' \
    "$gateway_base_url/stations"
}

deleteMyStationsByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stationIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/stations/batch"
}

createStations() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"createdStations":[{"description":"example","headerBackgroundURL":"https://example.com","icon":"example","id":"00000000-0000-4000-8000-000000000001","name":"example"}]}' \
    "$gateway_base_url/stations/batch"
}

updateMyStationsByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"updatedStations":[{"setNull":{},"stationId":"00000000-0000-4000-8000-000000000001","values":{"description":"example","headerBackgroundURL":"https://example.com","icon":"example","name":"example"}}]}' \
    "$gateway_base_url/stations/batch"
}

hardDeleteMyStationsByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stationIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/stations/batch/permanently"
}

restoreMyStationsByIds() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stationIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/stations/batch/restore"
}

leaveMyStations() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stations":[{"stationId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/stations/memberships/me"
}

visualizeMyTotalCount() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/stations/visualizations/total-count?permission=Read"
}

deleteMyStationById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stationId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/stations/${stationId}"
}

getMyStationById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/stations/${stationId}?isDeleted=true"
}

updateMyStationById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"description":"example","headerBackgroundURL":"https://example.com","icon":"example","name":"example"}}' \
    "$gateway_base_url/stations/${stationId}"
}

leaveMyStation() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/stations/${stationId}/memberships/me"
}

transferMyStationOwnership() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"targetUserPublicId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/stations/${stationId}/ownership"
}

hardDeleteMyStationById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stationId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/stations/${stationId}/permanently"
}

deleteMyStationPermissions() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"userPublicIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/stations/${stationId}/permissions"
}

upsertMyStationPermissions() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"permissions":[{"permission":"Read","userPublicId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/stations/${stationId}/permissions"
}

deleteMyStationPermission() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/stations/${stationId}/permissions/${userPublicId}"
}

getMyStationPermission() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/stations/${stationId}/permissions/${userPublicId}"
}

updateMyStationPermission() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/stations/${stationId}/permissions/${userPublicId}"
}

createMyStationPermission() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"permission":"Read"}' \
    "$gateway_base_url/stations/${stationId}/permissions/${userPublicId}"
}

upsertMyStationPermission() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/stations/${stationId}/permissions/${userPublicId}"
}

restoreMyStationById() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"stationId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/stations/${stationId}/restore"
}

deleteMySubShelvesByIds() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"subShelfIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/sub-shelves/batch"
}

createSubShelvesByRootShelfIds() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"createdSubShelves":[{"id":"00000000-0000-4000-8000-000000000001","name":"example","prevSubShelfId":"00000000-0000-4000-8000-000000000001","rootShelfId":"00000000-0000-4000-8000-000000000001"}]}' \
    "$gateway_base_url/sub-shelves/batch"
}

updateMySubShelvesByIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"updatedSubShelves":[{"setNull":{},"subShelfId":"00000000-0000-4000-8000-000000000001","values":{"name":"example"}}]}' \
    "$gateway_base_url/sub-shelves/batch"
}

moveMySubShelvesByRootShelfIds() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"movedSubShelves":[{"destinationRootShelfId":"00000000-0000-4000-8000-000000000001","destinationSubShelfId":"00000000-0000-4000-8000-000000000001","sourceRootShelfId":"00000000-0000-4000-8000-000000000001","sourceSubShelfIds":["00000000-0000-4000-8000-000000000001"]}]}' \
    "$gateway_base_url/sub-shelves/batch/position"
}

restoreMySubShelvesByIds() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"subShelfIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/sub-shelves/batch/restore"
}

moveMySubShelvesByRootShelfId() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"destinationRootShelfId":"00000000-0000-4000-8000-000000000001","destinationSubShelfId":"00000000-0000-4000-8000-000000000001","sourceRootShelfId":"00000000-0000-4000-8000-000000000001","sourceSubShelfIds":["00000000-0000-4000-8000-000000000001"]}' \
    "$gateway_base_url/sub-shelves/position"
}

getMySubShelvesByPrevSubShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/sub-shelves/prev-sub-shelf/${prevSubShelfId}?areDeleted=true"
}

getMySubShelvesAndItemsByPrevSubShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/sub-shelves/prev-sub-shelf/${prevSubShelfId}/items?areDeleted=true"
}

getAllMySubShelvesByRootShelfId() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/sub-shelves/root-shelf/${rootShelfId}?areDeleted=true"
}

createSubShelfByRootShelfId() {
  curl --fail-with-body --silent --show-error -X POST \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"id":"00000000-0000-4000-8000-000000000001","name":"example","prevSubShelfId":"00000000-0000-4000-8000-000000000001","rootShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/sub-shelves/root-shelf/${rootShelfId}"
}

deleteMySubShelfById() {
  curl --fail-with-body --silent --show-error -X DELETE \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/sub-shelves/${subShelfId}"
}

getMySubShelfById() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/sub-shelves/${subShelfId}?isDeleted=true"
}

updateMySubShelfById() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"name":"example"}}' \
    "$gateway_base_url/sub-shelves/${subShelfId}"
}

moveMySubShelfByRootShelfId() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"destinationRootShelfId":"00000000-0000-4000-8000-000000000001","destinationSubShelfId":"00000000-0000-4000-8000-000000000001","sourceRootShelfId":"00000000-0000-4000-8000-000000000001","sourceSubShelfId":"00000000-0000-4000-8000-000000000001"}' \
    "$gateway_base_url/sub-shelves/${subShelfId}/position"
}

restoreMySubShelfById() {
  curl --fail-with-body --silent --show-error -X PATCH \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    "$gateway_base_url/sub-shelves/${subShelfId}/restore"
}

getUserData() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/users/data"
}

getMe() {
  curl --fail-with-body --silent --show-error -X GET \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    "$gateway_base_url/users/me"
}

updateMe() {
  curl --fail-with-body --silent --show-error -X PUT \
    -b "$cookie_jar" -c "$cookie_jar" \
    -H "User-Agent: $user_agent" \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    --data '{"setNull":{},"values":{"displayName":"example","status":"example"}}' \
    "$gateway_base_url/users/me"
}
