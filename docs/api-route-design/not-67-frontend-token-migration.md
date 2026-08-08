# NOT-67 前端 Token 與 Gateway Response 遷移指南

本文是前端由舊版 Gateway response 與 `localStorage` access token 遷移到 typed Gateway contract 的施工文件。前端只需要依照本文件調整 API client、登入狀態與 CSRF 處理；Core 的 repository、service 與 internal delegation token 不屬於前端契約。

## 1. 遷移後的責任邊界

| 資料 | 維護者 | 前端是否可讀取 |
| --- | --- | --- |
| Access token | Gateway cookie，`HttpOnly` | 不可，也不應嘗試讀取 |
| Refresh token | Gateway cookie，`HttpOnly` | 不可，也不應嘗試讀取 |
| CSRF token | Gateway response data 或 refresh metadata；前端保留目前 session 的非敏感值 | 可以，僅用於 CSRF header |
| User public id、名稱、Email | `ClientResponse.data` | 可以 |
| Core delegation token | Gateway/Core internal envelope | 不可 |

前端不再自行保存、讀取或組裝 access token。瀏覽器會依照 Cookie 屬性自動保存與送出 access/refresh cookies，前端只負責帶上 `credentials`。

## 2. Public request 與 response 格式

對外 API 的請求由 `contracts/gateway/v1.ClientRequest[D]` 定義。它在 JSON 上會直接展開 DTO，因此前端送出的 body 仍然是該 route 的 DTO，不包含 `version`、`operation`、`metadata` 或 token。

一般成功回應使用 `contracts/gateway/v1.ClientResponse[D]`：

```json
{
  "success": true,
  "data": {},
  "exception": null
}
```

登入或註冊成功時，`data` 是不含 credential 的公開資料：

```json
{
  "success": true,
  "data": {
    "publicId": "user-public-id",
    "name": "jeff",
    "displayName": "Jeff",
    "email": "jeff@example.com",
    "csrfToken": "csrf-token",
    "createdAt": "2026-08-07T12:00:00Z",
    "updatedAt": "2026-08-07T12:00:00Z"
  },
  "exception": null
}
```

`data` 絕對不會包含 `accessToken` 或 `refreshToken`。這兩個值只會透過 Gateway 的 `Set-Cookie` 回應標頭傳遞。

失敗回應固定使用相同 envelope：

```json
{
  "success": false,
  "data": null,
  "exception": {
    "reason": "NotFound",
    "domain": "RootShelf",
    "operation": "Get",
    "message": "Root shelf was not found.",
    "retryable": false
  }
}
```

請使用 `success` 判斷結果，再依 `exception.reason`、`exception.domain` 與 `exception.operation` 顯示或處理錯誤；不要依賴舊版 controller 自行猜測的 `gin.H` 欄位。

## 3. Cookie 與 CSRF 呼叫方式

所有會使用登入狀態的瀏覽器請求都必須允許瀏覽器附帶 cookie：

```ts
await fetch(url, {
  method: "POST",
  credentials: "include",
  headers: {
    "Content-Type": "application/json",
    ...(csrfToken ? { "X-CSRF-Token": csrfToken } : {}),
  },
  body: JSON.stringify(dto),
});
```

若使用 Axios，API client 的單一 instance 設定 `withCredentials: true`：

```ts
const api = axios.create({
  withCredentials: true,
});
```

CSRF token 不是 access credential。前端可以將它放在記憶體中的 auth state；若產品需要跨頁保留，可放在非 credential 的 session state，但不能把 access/refresh token 放進任何 Web Storage。

當 Gateway 因 access cookie 到期而完成 refresh 時：

1. Gateway 重新設定 access cookie。
2. Gateway 在 response 的 `refreshableTokens.newCSRFToken` 提供新的 CSRF token。
3. Gateway 也可附帶 `X-CSRF-Token` header；前端應以 JSON 的 `refreshableTokens.newCSRFToken` 為主要來源，避免依賴 header expose 設定。
4. 前端更新記憶體中的 CSRF token，後續 mutation 使用新值。

前端不需要，也不應該自行呼叫 refresh token endpoint 或自行旋轉 access token。若 API client 有 retry 機制，最多只對同一 request retry 一次，避免失敗時形成無限迴圈。

### SSR / server function

Server function 不能依賴瀏覽器的 `credentials` 自動轉送 cookie，必須把目前 incoming
request 的完整 `Cookie` header 原樣轉送給 Gateway；收到 Gateway 的 `Set-Cookie`
則沿用既有的 `forwardUpstreamSetCookies` bridge 回傳給瀏覽器。這個 bridge 是轉送
HTTP cookie，不是把 token 讀成 JavaScript 字串。

因此 server function 應保留：

```ts
const inboundCookie = getRequestHeader("cookie");
const response = await fetch(url, {
  headers: inboundCookie ? { Cookie: inboundCookie } : {},
});
forwardUpstreamSetCookies(response);
```

並移除 `AccessTokenCookieHandler.ensure(...)`、從 Web Storage 組裝
`Authorization` header，以及任何把 `newAccessToken` 寫回 response state 的程式。

## 4. 移除 localStorage access token

請在前端 API client 與 auth store 完成下列變更：

1. 移除 `localStorage.setItem`、`getItem`、`removeItem` 對 access token 的所有呼叫。
2. 移除所有 `Authorization: Bearer ${localStorage...}` 的組裝邏輯。
3. 啟動新版本時可執行一次性清理，刪除舊 key；不要把舊值轉送給後端。
4. 登入、註冊成功後只保存公開 `data` 與 CSRF token，access/refresh token 不進入 React state、Redux、localStorage、sessionStorage、URL、analytics event 或 log。
5. 登出時呼叫 Gateway logout route，接著清除前端的公開 auth state 與 CSRF token；Cookie 由 Gateway 透過 `Set-Cookie` 清除。

目前前端掃描到的舊邏輯包含 `shared/api/interfaces/auth.interface.ts` 中的
`accessToken` response 欄位、`shared/api/cookies/accessToken.cookie.ts`、
`shared/api/cookies/refreshToken.cookie.ts`，以及各 feature hook/server function
對 `response.refreshableTokens.newAccessToken` 的寫入。這些欄位與 cookie helper
都應移除或改成只處理 CSRF；不要把它們改成另一種 localStorage key。

這個遷移不會改變前端原本的 route URL 或 DTO body。改變的是 credential 的保存位置與 response envelope，因此既有表單、列表與 GraphQL data mapping 不需要因 token 而加入特殊分支。

## 5. API client 建議介面

建議所有 HTTP 呼叫集中經過一個 client：

```ts
type ClientResponse<T> = {
  success: boolean;
  data: T;
  exception: {
    reason: string;
    domain: string;
    operation: string;
    message: string;
    retryable: boolean;
  } | null;
  refreshableTokens?: {
    newCSRFToken?: string;
  };
};

async function request<T>(input: RequestInfo, init?: RequestInit): Promise<T> {
  const response = await fetch(input, {
    ...init,
    credentials: "include",
  });
  const envelope = (await response.json()) as ClientResponse<T>;

  if (envelope.refreshableTokens?.newCSRFToken) {
    authStore.setCSRFToken(envelope.refreshableTokens.newCSRFToken);
  }
  if (!envelope.success) {
    throw toClientError(envelope.exception);
  }
  return envelope.data;
}
```

實際型別應由後端 `contracts/gateway/v1` 產生的 TypeScript contract 為準，不要在各 feature folder 重新宣告 response envelope。

## 6. 前端驗收清單

- [ ] 瀏覽器 request 帶有 `credentials: "include"`，Axios instance 帶有 `withCredentials: true`。
- [ ] 登入、註冊成功後能從 `ClientResponse.data` 取得使用者公開資料與 CSRF token。
- [ ] response JSON、console、analytics、error tracking 都找不到 access/refresh token。
- [ ] localStorage/sessionStorage 不再保存 access token 或 refresh token。
- [ ] mutation 會送出目前 CSRF token；refresh metadata 出現時會更新它。
- [ ] 登出後 auth state 清除，下一次需要登入的操作仍能正確得到 typed exception。
- [ ] API client 能處理 `data: null` 的失敗 response，不因泛型資料型別而 crash。
- [ ] 不依賴 `X-Core-*` headers，也不解析 Core internal envelope。
