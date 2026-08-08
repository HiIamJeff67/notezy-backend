# CI/CD 與 Jenkins Pipeline

NOT-43 定義 GitHub Actions 與 Jenkins 的分工。兩者都只能呼叫根目錄
`Makefile`／`internal/cli` 已存在的命令，不在 pipeline 內複製另一套測試流程。

## GitHub Actions

`.github/workflows/ci.yml` 是每次 pull request、`main`／`refactor/**` push 以及
version tag 的品質閘門，包含：

- 每個 Go module 的格式檢查、`go vet`、unit test 與 `go test -race`。
- GraphQL 產物重新生成後的 clean-tree 檢查。
- 五個 Go runtime 與 Yjs Worker 的 production container build。
- version tag 通過全部閘門後，將 runtime images 發布至 GitHub Container Registry。

`.github/workflows/integration.yml` 只在手動觸發或排程執行 Docker-backed
integration tests。它啟用 `NOTEZY_RUN_INTEGRATION=1`，執行 PostgreSQL／Redis／Kafka
Testcontainers 與 Core／DurableJob broker flow；因此一般 pull request 不會因本機
沒有 Docker 而失敗。

`.github/workflows/staging.yml` 是需要 `staging` environment approval 的 promotion
workflow。它在標記為 `staging` 的 self-hosted runner 上，以指定 GHCR tag 執行
`docker-compose.prod.yaml`，透過 `infra/staging/deploy.sh` 做 immutable image
promotion，再由 `infra/staging/smoke.sh` 檢查每個 runtime 的
`/startedz` 與 `/healthz`。Compose 支援 `GATEWAY_IMAGE`、`CORE_IMAGE`、
`DURABLE_JOB_IMAGE`、`EMAIL_IMAGE`、`REALTIME_GATEWAY_IMAGE` 與
`YJS_WORKER_IMAGE`，因此 promotion 不會重新 build image。staging runner 必須在
`/etc/notezy/staging.env` 提供環境設定；該檔案不由 repository checkout，也不應
提交 secrets。部署後的 Compose logs 會以 14 天 retention 上傳為 artifact。

## Jenkins

根目錄 `Jenkinsfile` 是 self-hosted agent 的 delivery pipeline：

1. checkout、format、vet、unit、race 與 generated contract gate；
2. 可選的 production container build；
3. 透過 `DEPLOY_STAGING` 與 `IMAGE_TAG` 參數，在標記為 `notezy-staging` 的
   agent 上 promotion immutable image 並執行相同的 smoke script；
4. 透過 `RUN_INTEGRATION` 參數選擇 Docker-backed integration tests。

Jenkins staging agent 必須提供 Docker、Compose plugin、Git 與可存取 GHCR 的
Jenkins Credential。staging deployment 只允許從已通過 CI 的 immutable image tag
promotion；不允許在 staging agent 重新編譯 source。

Jenkins 不取代 GitHub Actions 的 PR gate，也不在本 issue 執行正式 production
rollout、migration rollback 或 disaster recovery。image publish、環境 secrets 與
部署權限應由 Jenkins credentials／GitHub OIDC 管理，不得提交到 repository。

## 本地對應命令

```sh
make ci-format
make ci-vet
make ci-unit
make ci-race
make ci-generated
make ci-containers
```

staging runner 的 delivery commands 為：

```sh
IMAGE_REGISTRY=ghcr.io/ORG/REPO IMAGE_TAG=TAG \
COMPOSE_ENV_FILE=/etc/notezy/staging.env make staging-deploy

make staging-smoke
```

整合測試則使用 `make test-integration` 與 `make test-integration-kafka`；需要
Docker 或既有 Kafka broker 時才執行。
