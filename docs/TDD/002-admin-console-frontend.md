# 002 前端控制台技术设计

## 1. 背景

IAM 项目当前为纯后端 API 服务（Go + go-zero），无任何前端代码。需要构建 Web 端控制台，覆盖两类用户：

- **用户门户**：登录、注册、密码重置、个人信息管理
- **管理控制台**：租户管理、用户管理、角色权限、Token 管理、应用管理、内部客户端管理

第一阶段目标：完成所有 P0 需求（REQ-001/002/003/004/005/006/007/012/017/018）的前端页面。

## 2. 技术选型

| 项目 | 选型 | 说明 |
|------|------|------|
| 框架 | Vue 3 (Composition API) + Vite + TypeScript | 类型安全、构建快速 |
| 路由 | Vue Router 4 | 官方路由 |
| 状态管理 | Pinia | Vue 3 官方推荐 |
| UI 组件库 | Element Plus | 管理后台生态成熟 |
| HTTP 客户端 | Axios | 拦截器、类型支持 |
| 构建工具 | Vite | 热更新、快速构建 |
| 测试框架 | Vitest + Vue Test Utils（测试）；Playwright（E2E）；agent-browser（开发期交互式验证） |

## 3. 项目结构

> 前端构建文件（`package.json`、`vite.config.ts`、`tsconfig.json` 等）放在项目根目录，`npm install` 无需切换目录。
> 前端源码放在 `web/` 目录下。

```
iam/
├── web/                          # 前端源码目录
│   ├── public/                   # 静态资源
│   ├── index.html                # HTML 入口
│   └── src/
│       ├── api/                  # API 请求封装
│       │   ├── request.ts        # Axios 实例 + 拦截器
│       │   ├── auth.ts           # 登录、注册、密码重置
│       │   ├── user.ts           # 用户 CRUD
│       │   ├── role.ts           # 角色管理
│       │   ├── permission.ts     # 权限分配
│       │   ├── tenant.ts         # 租户管理
│       │   ├── token.ts          # Token 管理
│       │   ├── app.ts            # 应用管理
│       │   └── client.ts         # 内部服务认证
│       ├── components/           # 公共组件
│       │   └── Layout/           # 侧边栏 + 顶栏 + 内容区
│       ├── router/
│       │   ├── index.ts          # 路由定义
│       │   └── guards.ts         # 路由守卫
│       ├── stores/
│       │   └── auth.ts           # 认证状态
│       ├── views/
│       │   ├── auth/             # 认证页（无 Layout）
│       │   │   ├── Login.vue
│       │   │   ├── Register.vue
│       │   │   └── ResetPassword.vue
│       │   └── admin/            # 管理控制台（含 Layout）
│       │       ├── Dashboard.vue
│       │       ├── Tenant/       # 租户管理
│       │       ├── User/         # 用户管理
│       │       ├── Role/         # 角色管理
│       │       ├── Permission/   # 权限分配
│       │       ├── Token/        # Token 管理
│       │       ├── App/          # 应用管理
│       │       └── Client/       # 客户端管理
│       ├── App.vue
│       └── main.ts
├── node_modules/                 # 前端依赖
├── package.json                  # 前端构建配置（新增）
├── vite.config.ts                # Vite 配置（新增）
├── tsconfig.json                 # TypeScript 配置（新增）
├── tsconfig.node.json            # Vite 的 Node 配置（新增）
├── .eslintrc.cjs                 # ESLint（可选）
├── .prettierrc                   # Prettier（可选）
├── .gitignore                    # 更新，添加 node_modules/
├── internal/                     # Go 后端代码（不变）
├── app/                          # Go 应用入口（不变）
├── docker-compose.yml            # 需要增加前端服务
└── ...
```

## 4. 路由设计

```
/                           → 重定向到 /login 或 /dashboard
/login                      → 登录页（REQ-001）
/register                   → 注册页（REQ-002）
/reset-password             → 密码重置（REQ-003）

/admin                      → 管理控制台 Layout
  /admin/dashboard          → 仪表盘
  /admin/tenants            → 租户管理（REQ-007）
  /admin/users              → 用户管理（REQ-004）
  /admin/roles              → 角色管理（REQ-005）
  /admin/permissions        → 权限分配（REQ-006）
  /admin/tokens             → Token 管理（REQ-012）
  /admin/apps               → 应用管理（REQ-017）
  /admin/clients            → 内部服务认证（REQ-018）
```

## 5. API 适配

- **Axios 拦截器**：自动注入 `Authorization: Bearer <token>`
- **401 处理**：Token 过期自动跳转登录页
- **统一响应格式**：后端 `{"code": 0, "message": "success", "data": {...}}` 格式，详见 `05-api-design.md`
- **开发代理**：Vite proxy → `http://localhost:8888`（本地开发）或 `http://iam:8888`（Docker Compose 环境）
- **生产部署**：Go 静态文件服务或 Nginx serve 构建产物

## 6. 实施步骤

1. **初始化项目**：`web/` 目录下创建 Vue 3 + TypeScript 项目，安装依赖
2. **基础架构**：Axios 封装、Pinia 认证 store、路由守卫、Layout 组件
3. **用户门户**：登录、注册、密码重置页面
4. **管理控制台**：7 个模块的 CRUD 页面
5. **集成部署**：更新 docker-compose、配置生产构建

## 7. 前端测试策略

前端测试分三层，与后端测试策略（`001-iam-system-architecture-design.md` 第 6 节）对齐。

### 7.1 测试分层

| 层级 | 范围 | 工具 | 覆盖率目标 |
|------|------|------|------------|
| **单元测试** | 工具函数、类型守卫、格式化函数、表单校验逻辑 | Vitest | > 80% |
| **组件测试** | Vue 组件独立渲染、交互行为、表单提交、错误展示 | Vitest + Vue Test Utils | 核心组件 100% |
| **E2E 测试** | 完整浏览器流程：登录 → 导航 → CRUD → 登出 | Playwright | 关键路径 100% |
| **开发验证** | 开发过程中交互式验证页面渲染和交互 | agent-browser | 每页必验 |

### 7.2 单元测试

测试纯逻辑，不渲染组件。

**测试范围：**

- 工具函数（`src/utils/`）：日期格式化、权限码解析、状态映射
- Axios 封装（`src/api/request.ts`）：拦截器逻辑、Token 注入、错误处理
- Store（`src/stores/`）：Pinia store 的状态变更、登录/登出动作

**测试文件命名：** `src/utils/__tests__/xxx.test.ts`

**示例（拦截器测试）：**

```typescript
// src/api/__tests__/request.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'

describe('request interceptor', () => {
  it('injects Authorization header when token exists', () => {
    localStorage.setItem('token', 'test-token')
    // ... assert config.headers.Authorization === 'Bearer test-token'
  })

  it('skips Authorization for login/register endpoints', () => {
    // ... assert /auth/login requests are not intercepted
  })
})

describe('response interceptor', () => {
  it('returns data when code is 0', () => {
    // ... assert { code: 0, data: { id: 1 } } => { id: 1 }
  })

  it('throws error with message when code is non-zero', () => {
    // ... assert { code: 20001, message: '用户不存在' } throws
  })

  it('redirects to /login on 401', () => {
    // ... mock router.push('/login')
  })
})
```

### 7.3 组件测试

使用 Vue Test Utils 挂载单个组件，验证渲染和行为。

**测试范围：**

- 登录表单：输入校验、提交、错误提示展示
- 数据表格：分页、搜索、空状态
- CRUD 对话框：创建/编辑表单、校验、提交
- 侧边栏：菜单展开/收起、当前路由高亮

**测试文件命名：** `src/components/__tests__/Xxx.test.ts`

**示例（登录表单测试）：**

```typescript
// src/views/auth/__tests__/Login.test.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Login from '../Login.vue'

describe('Login', () => {
  it('renders login form with email and password fields', () => {
    const wrapper = mount(Login)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('shows error when email is empty', async () => {
    const wrapper = mount(Login)
    await wrapper.find('form').trigger('submit')
    expect(wrapper.text()).toContain('请输入邮箱')
  })

  it('calls login API on valid submit', async () => {
    const mockLogin = vi.fn().mockResolvedValue({ code: 0, data: { token: 'x' } })
    // ... mount with mocked api
    await wrapper.find('form').trigger('submit')
    expect(mockLogin).toHaveBeenCalledWith({ email: 'a@b.com', password: 'pwd' })
  })

  it('shows error message when login fails', async () => {
    // ... mock API returning { code: 10001, message: '用户名或密码错误' }
    expect(wrapper.text()).toContain('用户名或密码错误')
  })
})
```

### 7.4 E2E 测试（正式测试：Playwright）

使用 Playwright 在真实浏览器中运行，覆盖关键用户路径，承担**正式端到端测试**职责。

**测试范围：**

| 场景 | 流程 |
|------|------|
| 登录闭环 | 打开首页 → 输入凭证 → 登录成功 → 跳转 Dashboard |
| 登录失败 | 输入错误密码 → 显示错误 → 不跳转 |
| 租户 CRUD | 登录 → 进入租户管理 → 创建租户 → 编辑 → 删除 → 验证 |
| 用户 CRUD | 登录 → 进入用户管理 → 创建用户 → 分配角色 → 搜索 → 禁用 |
| Token 过期 | 登录后等待 Token 过期 → 操作 → 自动跳转登录页 |
| 权限控制 | 未登录访问 /admin → 重定向 /login |

**测试文件位置：** `tests/e2e/`

**视觉回归：** 使用 Playwright `expect(page).toHaveScreenshot()` 做截图对比。

**示例（登录 E2E）：**

```typescript
// tests/e2e/login.spec.ts
import { test, expect } from '@playwright/test'

test('login success redirects to dashboard', async ({ page }) => {
  await page.goto('/login')
  await page.getByPlaceholder('邮箱').fill('admin@example.com')
  await page.getByPlaceholder('密码').fill('Admin@123')
  await page.getByRole('button', { name: '登录' }).click()
  await expect(page).toHaveURL(/\/dashboard/)
})

test('invalid credentials show error', async ({ page }) => {
  await page.goto('/login')
  await page.getByPlaceholder('邮箱').fill('admin@example.com')
  await page.getByPlaceholder('密码').fill('wrong')
  await page.getByRole('button', { name: '登录' }).click()
  await expect(page.getByText('用户名或密码错误')).toBeVisible()
})

test('unauthenticated access to admin redirects to login', async ({ page }) => {
  await page.goto('/admin/tenants')
  await expect(page).toHaveURL(/\/login/)
})
```

### 7.5 开发期交互式验证（agent-browser）

在开发过程中，**每完成一个页面或组件**后，必须使用 `agent-browser` 进行交互式验证，确认页面渲染和交互行为符合预期。这不是 CI 测试，而是**开发流程的强制步骤**。

**使用要求：**

| 要求 | 说明 |
|------|------|
| 有头模式 | 必须使用 `--headed`，让用户在开发过程中能实时看到浏览器操作 |
| 使用 Chrome | 通过 `--engine chrome` 指定 Chrome 引擎 |
| 不使用用户配置 | **禁止**使用 `--profile` 参数，避免用户本地 Chrome 的 session/token 残留影响验证结果 |
| 独立 session | 使用 `--session-name dev` 实现状态隔离和自动清理 |

**验证时机：**

每完成以下工作，必须跑一次 agent-browser 验证：
- 新建一个页面（登录页、租户列表页等）
- 修改了核心交互（表单提交、弹窗、路由跳转）
- 调整了布局/样式（确认视觉表现正确）

**典型使用流程：**

```bash
# 1. 打开页面（有头模式、Chrome、无用户配置）
agent-browser --headed --engine chrome --session-name iam-dev open http://localhost:5173/login

# 2. 等待页面加载完成，截图查看渲染效果
agent-browser wait --load networkidle
agent-browser screenshot dev-verify/login-rendered.png

# 3. 检查表单元素是否正确渲染
agent-browser snapshot -i
# 输出示例：
# @e1 [input type="email"] "邮箱"
# @e2 [input type="password"] "密码"
# @e3 [button] "登录"

# 4. 尝试交互：填写并提交
agent-browser fill @e1 "admin@example.com"
agent-browser fill @e2 "Admin@123"
agent-browser click @e3

# 5. 验证跳转结果
agent-browser wait --url "**/dashboard"
agent-browser screenshot dev-verify/dashboard-rendered.png

# 6. 开发完成后关闭，session 自动清理
agent-browser close
```

**便捷脚本（推荐添加到 `package.json`）：**

```json
{
  "scripts": {
    "dev:verify": "agent-browser --headed --engine chrome --session-name iam-dev"
  }
}
```

开发时运行：`npm run dev:verify open http://localhost:5173/login`

### 7.6 测试目录

```
web/
├── src/
│   ├── api/__tests__/
│   │   └── request.test.ts
│   ├── stores/__tests__/
│   │   └── auth.test.ts
│   ├── utils/__tests__/
│   │   ├── format.test.ts
│   │   └── validators.test.ts
│   └── views/auth/__tests__/
│       └── Login.test.ts
├── tests/
│   └── e2e/                    # Playwright E2E 测试
│       ├── login.spec.ts
│       ├── tenant.spec.ts
│       └── user.spec.ts
└── dev-verify/                 # agent-browser 开发验证截图（不提交到 git）
    ├── login-rendered.png
    └── dashboard-rendered.png
```

### 7.7 运行命令

| 命令 | 说明 |
|------|------|
| `npm run test` | 运行所有单元测试 + 组件测试 |
| `npm run test:unit` | 仅单元测试 |
| `npm run test:component` | 仅组件测试 |
| `npm run test:e2e` | 运行 Playwright E2E 测试 |
| `npm run test:e2e -- --ui` | Playwright UI 模式（可视化调试） |
| `npm run test -- --coverage` | 生成覆盖率报告 |
| `npm run dev:verify` | 启动 agent-browser 开发期验证（有头、Chrome、无用户配置） |

### 7.8 CI 集成

在现有 CI 流水线（`001-iam-system-architecture-design.md` 第 6 节）中增加前端测试步骤：

```
push/PR
  ├─ 1. lint:         golangci-lint + eslint
  ├─ 2. build:        go build ./... && npm run build
  ├─ 3. unit:         go test ./... -race -cover + npm run test:unit
  ├─ 4. integration:  docker compose up deps → go test -tags=integration
  ├─ 5. api:          启动服务 → 跑用例
  └─ 6. e2e:          启动前后端 → Playwright 跑 E2E 测试 + 截图对比
```

> agent-browser 仅用于开发期交互式验证，不进入 CI 流程。

## 8. 验证方式

1. 启动后端：`go run app/main.go -f etc/dev.yaml`
2. 启动前端：`npm run dev`（构建文件在根目录，无需切换目录）
3. 浏览器访问 `http://localhost:5173`
4. 构建测试：`npm run build`
5. 运行测试：`npm run test`（单元 + 组件）、`npm run test:e2e`（端到端）
