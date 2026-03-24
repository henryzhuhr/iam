# REQ-014 用户组管理

| 项目 | 内容 |
|------|------|
| **优先级** | P1 |
| **估时** | 4 人天 |
| **关联用户故事** | US-016、US-017 |

**背景：** 企业用户需要通过用户组来批量管理用户和权限，支持组织架构的层级关系，简化大规模用户管理成本。

**目标：**

- 支持用户组的创建、查询、更新、删除
- 支持树形层级结构（父子用户组）
- 支持用户组成员管理
- 支持权限继承（子组继承父组权限）
- 支持批量用户操作

**功能描述：**

### 1. 用户组 CRUD

1. 创建用户组：指定名称、描述、父组（可选）
2. 查询用户组：支持列表查询、树形结构查询
3. 更新用户组：修改名称、描述、父组
4. 删除用户组：支持级联删除或移动到默认组

约束：
- 用户组名称在租户内唯一
- 不能将自己设置为父组（禁止循环引用）
- 最大层级深度限制（默认 10 级）
- 删除有子组的组时，需先处理子组

### 2. 树形层级结构

1. 支持无限层级（受最大深度限制）
2. 查询时支持返回完整路径
3. 支持按父组查询子组（递归/单层可选）
4. 支持获取节点的所有祖先/后代

示例层级：
```
总公司
├── 研发中心
│   ├── 前端组
│   ├── 后端组
│   └── 测试组
├── 产品中心
│   ├── 产品组
│   └── 设计组
└── 职能中心
    ├── 人力资源
    └── 财务组
```

### 3. 用户组成员管理

1. 添加成员：将用户添加到用户组
2. 移除成员：从用户组移除用户
3. 批量添加：支持 Excel 导入或选择多个用户
4. 成员列表：查询组内所有成员（含继承成员）
5. 成员去重：一个用户可在多个组，但组内不重复

### 4. 权限继承

1. 子组自动继承父组的所有权限
2. 用户权限 = 所属所有组的权限并集
3. 支持权限覆盖配置（允许/拒绝）
4. 权限变更时，继承关系实时生效

权限计算优先级：
```
用户直接权限 > 拒绝策略 > 允许策略 > 默认拒绝
```

### 5. 用户组角色

1. 每个用户组可设置一个组管理员
2. 组管理员可管理本组成员
3. 组管理员不能超越租户管理员权限

**用户组类型：**

| 类型 | 说明 | 可删除 |
|------|------|--------|
| 系统组 | 系统预置（如 All Users） | 否 |
| 普通组 | 租户自定义 | 是 |
| 动态组 | 基于规则自动成员 | 是 |

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| 循环引用 | 拒绝设置，提示「不能设置后代节点为父组」 |
| 超出最大层级 | 拒绝创建，提示「不能超过 N 级」 |
| 组名重复 | 拒绝创建，提示「组名已存在」 |
| 删除非空组 | 提示先移除成员或子组 |
| 用户已在组内 | 跳过或提示已存在 |

**API 接口：**

```
# 用户组管理
POST   /api/v1/user-groups              # 创建用户组
GET    /api/v1/user-groups              # 用户组列表
GET    /api/v1/user-groups/tree         # 树形结构查询
GET    /api/v1/user-groups/:id          # 用户组详情
PUT    /api/v1/user-groups/:id          # 更新用户组
DELETE /api/v1/user-groups/:id          # 删除用户组

# 成员管理
GET    /api/v1/user-groups/:id/members  # 获取组成员列表
POST   /api/v1/user-groups/:id/members  # 添加组成员
DELETE /api/v1/user-groups/:id/members/:userId  # 移除组成员
POST   /api/v1/user-groups/:id/members/batch  # 批量添加成员

# 权限管理
GET    /api/v1/user-groups/:id/permissions    # 获取组权限
POST   /api/v1/user-groups/:id/permissions    # 分配权限
```

**数据库设计：**

```sql
-- 用户组表
CREATE TABLE user_groups (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    parent_id BIGINT DEFAULT NULL,       -- 父组 ID
    level INT DEFAULT 1,                 -- 层级深度
    path VARCHAR(500),                   -- 完整路径如 /1/5/12/
    is_system BOOLEAN DEFAULT FALSE,     -- 是否系统组
    group_type VARCHAR(20) DEFAULT 'NORMAL', -- NORMAL/DYNAMIC
    sort_order INT DEFAULT 0,            -- 排序号
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_parent (parent_id),
    INDEX idx_path (path(100))
);

-- 用户组成员表
CREATE TABLE user_group_members (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,                   -- 操作人
    UNIQUE KEY uk_group_user (group_id, user_id),
    INDEX idx_user (user_id)
);

-- 用户组权限表
CREATE TABLE user_group_permissions (
    id BIGINT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_group_permission (group_id, permission_id),
    INDEX idx_group (group_id)
);

-- 用户组管理员表
CREATE TABLE user_group_admins (
    id BIGINT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_group_user (group_id, user_id)
);
```

**验收标准：**

- [ ] 用户组可正常创建、查询、更新、删除
- [ ] 树形层级结构正确维护
- [ ] 循环引用被正确阻止
- [ ] 成员可正常添加和移除
- [ ] 权限继承正确生效
- [ ] 用户权限计算为所有组权限并集
- [ ] 组管理员只能管理本组
