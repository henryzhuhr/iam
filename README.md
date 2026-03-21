# IAM

身份认证与访问管理 (Identity and Access Management)

## Run Service

```bash
go run app/main.go -f etc/dev.yaml
```

## API Tests

项目使用 `pytest` 编写黑盒 HTTP 接口测试。当前框架默认会在测试会话内自动拉起 Go 服务，再命中本地 `/api` 前缀；如果你已经有一个外部服务实例，也可以显式切换到“复用已有服务”模式。

```bash
uv run pytest
uv run pytest -m smoke
uv run pytest --app-config etc/dev.yaml --base-url http://127.0.0.1:8080/api
IAM_USE_EXISTING_SERVICE=1 IAM_BASE_URL=http://127.0.0.1:8080/api uv run pytest
```

可选环境变量：

- `IAM_BASE_URL`: 覆盖默认接口地址
- `IAM_USE_EXISTING_SERVICE`: 复用外部已启动服务，不由 pytest 管理 Go 进程
- `IAM_APP_ENTRY`: 覆盖默认 Go 入口文件
- `IAM_APP_CONFIG`: 覆盖默认配置文件
- `IAM_GO_RUN_EXTRA_ARGS`: 追加到 `go run` 后面的额外参数
- `IAM_APP_STARTUP_TIMEOUT`: 覆盖服务启动健康检查超时
- `IAM_TOKEN`: 预留给后续鉴权接口测试的 Bearer Token
- `IAM_TENANT_ID`: 预留给多租户请求头注入

目录约定：

- `tests/api/`: 接口测试用例
- `tests/support/`: 客户端、断言和公共 fixture 支撑代码
- `tests/conftest.py`: pytest 全局参数、服务进程管理和 fixture 入口

## Issues

- 项目内使用 `issues/` 目录记录 issue。
- issue 文件名格式为 `NNN-short-kebab-case.md`。
- 新增 issue 时，需要同步更新对应目录的 `README.md`，维护按编号排序的 index。
