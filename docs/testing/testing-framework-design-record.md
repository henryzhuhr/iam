# IAM Pytest API 测试框架设计文档

> 目标：把这套 `pytest` API 测试框架的设计背景、问题澄清、用户决策和落地约定记录成一份可持续迭代的文档。
> 状态：已实现 v1

---

## 1. 一页结论

 IAM 当前采用一套本地优先的 `pytest + httpx` 黑盒 HTTP API 测试框架。默认情况下，`pytest` 会在测试会话内自动拉起 Go 服务并做健康检查；如果需要命中已有环境，也支持显式切换到外部服务模式。

这一版先把“测试框架本身”搭稳，只覆盖当前已落地的 `/api/health`，同时为认证、租户上下文、测试数据工厂、契约测试和 CI 接入预留扩展点。这样做的目标是先降低新增接口测试的门槛，同时把“程序能否成功启动”也纳入自动化测试的一部分。

## 2. Grounding：设计前确认过的仓库事实

这些内容不是主观偏好，是当时在仓库里先确认过的现状：

- Go 服务当前真实落地的 HTTP 路由只有 `/api/health`
- Python 侧已有 `pyproject.toml` 和 `uv`，开发依赖里已有 `pytest`
- 仓库里当时没有现成的 `tests/` 目录，也没有任何 `*_test.go`
- 当前 `.github/workflows/` 里没有接口测试流水线，只有镜像构建工作流
- 当前 `health` 的 Swagger 示例与真实实现不一致：
  - Swagger 示例：`{"status":"healthy"}`
  - 当前运行时返回：`{"message":"ok"}`

这些事实直接影响了框架第一版的范围选择。如果接口数量很少、服务启动很轻、依赖编排又还没稳定，那么先做黑盒 API 测试骨架是成本最低、最不容易返工的方案。

## 3. Questions And Decisions

这一节记录当时真正问过你的问题、你做出的选择，以及这些选择对实现的约束。随着框架演进，下面的解释已经同步更新为当前实现版本。

### Q1. 这套 pytest API 测试框架，v1 你希望做到哪一层？

候选方向：

- 黑盒 HTTP
- 自启动服务
- 连同集成环境一起设计

最终决策：

- 选择 `黑盒 HTTP`

决策含义：

- 测试框架的核心入口是 HTTP 请求和响应断言
- 默认不走 UI 或浏览器端到端测试路径
- 除 Go 服务进程本身外，v1 仍不统一编排 MySQL、Redis、Kafka 等外部依赖生命周期
- 新增测试时优先复用 HTTP 客户端、fixture 和断言支撑层

### Q2. 对认证、测试数据、租户等业务夹具，v1 你希望做到什么程度？

候选方向：

- 只预留扩展点
- 带认证骨架
- 带数据工厂

最终决策：

- 选择 `只预留扩展点`

决策含义：

- `auth_session`、`request_context` 这些 fixture 先定义出来
- 先把扩展位留好，但不在 v1 强行把登录、造数、清理逻辑塞进框架
- 后续业务接口增多时，再按实际需要补 token 获取、租户上下文和测试数据工厂

### Q3. v1 的默认断言应该以哪个契约来源为准？

候选方向：

- 以当前实现为准
- 以 Swagger 为准
- 两者都要

最终决策：

- 选择 `以当前实现为准`

决策含义：

- 示例测试先验证真实运行结果，不用 Swagger 差异卡住第一版框架
- 文档与实现不一致的问题单独记录为契约漂移
- 后续如果要做 Swagger/OpenAPI 比对，放进 `contract` 类测试治理

### Q4. 测试执行时由谁负责启动 Go 服务？

候选方向：

- 手动提前启动服务
- 默认由 `pytest` 自动启动
- 同时支持两种模式

最终决策：

- 选择 `默认由 pytest 自动启动，同时保留外部服务模式作为显式覆盖`

决策含义：

- 服务能否正常启动本身就是测试的一部分
- 默认测试路径不依赖人工提前启动
- 如果需要命中特定环境，允许通过配置切换到已有服务模式
- Go 服务参数注入改为交给 pytest 配置项和环境变量管理

## 4. Why This Shape

这套方案最后落成现在这个形状，核心原因有四个：

1. 当前接口面很小。只有 `/api/health`，先做完整环境编排会明显重于业务价值。
2. 现有 Go 服务启动足够轻。对于当前版本，把它交给 `pytest` 在 session 级托管成本可控，且能提升自动化程度。
3. Python 测试基础已具备。仓库已有 `uv` 和 `pytest`，继续沿用能避免引入第二套测试运行方式。
4. 当前最稀缺的是统一约定，而不是更复杂的基础设施。先把目录、fixture、客户端、断言风格和运行方式固定下来，后续扩展成本最低。
5. 服务可启动性本身就是自动化测试的重要一环。把 Go 服务交给 pytest 托管，才能真正支持参数注入和无人值守执行。

## 5. v1 设计定稿

### 5.1 测试模型

- 测试类型：黑盒 HTTP API 测试
- 测试框架：`pytest`
- HTTP 客户端：`httpx`
- 默认基础地址：`http://127.0.0.1:8080/api`
- 服务生命周期：默认由 `pytest` 在 session 级自动启动和停止
- 覆盖模式：可通过显式配置切换到外部已启动服务

### 5.2 测试目录

```bash
tests/
├── hooks/
│   └── cli.py
├── fixtures/
│   └── runtime.py
├── framework/
│   ├── test_config.py
│   └── test_helpers.py
├── api/
│   └── health/
│       └── test_health.py
├── helpers/
│   ├── assertions.py
│   ├── client.py
│   ├── config.py
│   └── app.py
└── conftest.py
```

目录职责：

- `tests/hooks/`: 放 pytest hook，例如 `pytest_addoption`
- `tests/fixtures/`: 放共享 fixture，例如 `test_config`、`managed_go_app`、`api_client`
- `tests/framework/`: 放框架自测，覆盖配置解析、公共断言和客户端辅助逻辑
- `tests/api/`: 放业务接口测试，按业务域拆分
- `tests/helpers/config.py`: 解析 `base_url`、环境变量等运行配置
- `tests/helpers/client.py`: 统一请求客户端、认证会话、请求上下文
- `tests/helpers/assertions.py`: 放响应断言工具
- `tests/helpers/app.py`: 放 Go 服务生命周期托管逻辑
- `tests/conftest.py`: 只做 pytest 插件装配

### 5.3 固定的公共接口

pytest CLI 参数：

- `--base-url`
- `--use-existing-service`
- `--app-entry`
- `--app-config`
- `--go-run-extra-args`
- `--app-startup-timeout`
- `--env`

环境变量：

- `IAM_BASE_URL`
- `IAM_USE_EXISTING_SERVICE`
- `IAM_APP_ENTRY`
- `IAM_APP_CONFIG`
- `IAM_GO_RUN_EXTRA_ARGS`
- `IAM_APP_STARTUP_TIMEOUT`
- `IAM_ENV`
- `IAM_TOKEN`
- `IAM_TENANT_ID`

固定 fixture：

- `test_config`
- `auth_session`
- `request_context`
- `managed_go_app`
- `api_client`

固定 markers：

- `api`
- `smoke`
- `contract`

### 5.4 默认行为约定

- 所有 HTTP 请求统一走 `api_client`
- 默认测试真实接口行为，不默认做 Swagger 契约校验
- 所有测试文件使用 `test_*.py` 命名
- 新增接口测试按业务域放到 `tests/api/<domain>/`

## 6. 当前落地情况

当前已经实现的内容包括：

- `pyproject.toml` 中的 `pytest` 配置和 `httpx` 依赖
- `tests/hooks/cli.py` 中的 CLI 参数注册
- `tests/fixtures/runtime.py` 中的服务进程管理和全局 fixture
- `tests/helpers/app.py` 中的 Go 服务生命周期托管
- `tests/helpers/` 下的客户端、配置、断言支撑层
- `tests/api/health/test_health.py` 中的健康检查 smoke 用例
- `README.md` 中的最小运行说明

当前健康检查用例断言的是：

- `GET /health`
- 返回 `200`
- 返回体是 JSON
- 当前内容为 `{"message": "ok"}`

## 7. 运行方式

默认模式下直接运行测试，pytest 会自动拉起 Go 服务：

```bash
uv run pytest
uv run pytest -m smoke
uv run pytest --app-config etc/dev.yaml --base-url http://127.0.0.1:8080/api
```

如果已经有外部服务实例，需要显式切换模式：

```bash
IAM_USE_EXISTING_SERVICE=1 IAM_BASE_URL=http://127.0.0.1:8080/api uv run pytest
```

## 8. 新增一个接口测试的操作模板

这一节不是设计说明，而是后续维护者可以直接照着执行的操作模板。

### 8.1 先判断这个测试该放哪里

新增接口测试时，先做这三个判断：

1. 这个接口属于哪个业务域，例如 `auth`、`users`、`roles`。
2. 这个测试只是普通接口测试，还是需要加入 `smoke`。
3. 这个测试是否依赖认证、租户上下文或测试数据准备。

默认规则：

- 普通接口测试使用 `pytest.mark.api`
- 只有核心可用性检查才额外加 `pytest.mark.smoke`
- 如果测试需要鉴权、租户上下文或造数能力，先复用现有 fixture；确实缺能力时，再扩展 `auth_session`、`request_context` 或新增专门 helper，不要把临时逻辑散落在单个测试文件里
- 如果测试需要不同的 Go 启动参数，优先通过 `--app-config`、`--app-entry` 或 `--go-run-extra-args` 注入，而不是手动改运行方式

### 8.2 标准操作步骤

1. 在 `tests/api/<domain>/` 下创建测试文件，命名为 `test_<feature>.py`
2. 在文件顶部声明 `pytestmark`，至少包含 `pytest.mark.api`
3. 使用 `api_client` 发起请求，不直接在用例里创建裸 `httpx.Client`
4. 使用 `tests/helpers/assertions.py` 里的 `assert_json_response` 做 JSON 响应断言
5. 先断言状态码和基础返回结构，再断言业务字段
6. 直接运行当前文件，再按需要跑 `-m smoke` 或完整 `pytest`；默认不需要手动先起 Go 服务
7. 如果本次新增测试引入了新的公共约定，同步更新本文档

### 8.3 最小模板

```python
import pytest

from tests.helpers.assertions import assert_json_response

pytestmark = [pytest.mark.api]


def test_example(api_client) -> None:
    response = api_client.get("/example")

    assert response.status_code == 200

    payload = assert_json_response(response)
    assert payload["message"] == "ok"
```

替换时只做这几件事：

- 把 `/example` 换成真实接口路径
- 按接口语义调整 marker
- 把 `payload["message"] == "ok"` 换成真实业务断言

### 8.4 完成检查清单

提交前至少确认下面这些项：

- 文件位置正确，在 `tests/api/<domain>/` 下
- 文件名符合 `test_*.py`
- 所有请求都走 `api_client`
- JSON 响应统一通过 `assert_json_response` 处理
- 没有把认证、租户、造数等一次性逻辑硬编码进单个测试
- 已执行 `uv run ruff check .`，并确认 Ruff 结果为通过
- 已执行至少一种对应的 pytest 命令验证用例能跑通
- 如果测试方式或公共约定变了，文档已同步更新

## 9. 以后怎么迭代

建议按照下面顺序继续扩：

1. 在 `auth_session` 上补登录能力，让认证接口测试先成型
2. 增加 `tests/api/auth/`、`tests/api/users/`、`tests/api/roles/` 等业务域目录
3. 引入数据工厂或测试 helper，统一准备用户、租户、角色等测试数据
4. 单独扩展 `contract` 测试，显式治理 Swagger 与实现漂移
5. 等运行模型稳定后，再接 GitHub Actions 的 smoke/api 测试流水线

## 10. 文档维护规则

- 如果修改了测试目录结构、fixture 约定、运行命令或扩展策略，优先更新本文档
- 如果只是补具体接口的业务行为约束，优先更新 `docs/REQUIREMENTS.md`
- 如果后续测试文档增多，继续在 `docs/testing/` 下按语义化文件名扩展，而不是把所有内容堆回一个文件
