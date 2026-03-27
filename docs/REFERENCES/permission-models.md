# 权限模型对比与选型

> 最后更新：2026-03-25
> 适用场景：IAM 权限系统设计

## 1. 权限模型概述

权限模型是用于定义和管理用户如何访问系统资源的抽象框架。常见的权限模型包括：

| 模型 | 全称 | 核心思想 |
|------|------|----------|
| **RBAC** | Role-Based Access Control | 基于角色的访问控制 |
| **ABAC** | Attribute-Based Access Control | 基于属性的访问控制 |
| **ReBAC** | Relationship-Based Access Control | 基于关系的访问控制 |
| **ACL** | Access Control List | 访问控制列表 |

---

## 2. RBAC（基于角色的访问控制）⭐

IAM 系统当前选择的权限模型。

### 2.1 核心概念

```
用户 (User) → 角色 (Role) → 权限 (Permission) → 资源 (Resource)
```

| 概念 | 说明 | 示例 |
|------|------|------|
| **用户** | 系统的使用者 | 张三、李四 |
| **角色** | 权限的集合，代表一类职能 | 管理员、编辑、普通用户 |
| **权限** | 对资源的操作许可 | `user:read`、`user:write` |
| **资源** | 被保护的对象 | 用户 API、订单 API |

### 2.2 角色层级

RBAC 支持角色继承，减少权限重复配置：

```
超级管理员
    ├── 系统管理员
    │     ├── 用户管理员
    │     └── 配置管理员
    └── 审计员
```

- 子角色自动继承父角色的所有权限
- 可以同时属于多个角色

### 2.3 数据结构

```sql
-- 用户 - 角色关联
user_roles (user_id, role_id)

-- 角色 - 权限关联
role_permissions (role_id, permission_id)

-- 权限定义
permissions (id, resource, action)
```

### 2.4 权限检查

```go
// 伪代码
func HasPermission(user, resource, action) bool {
    roles := GetUserRoles(user)
    for _, role := range roles {
        permissions := GetRolePermissions(role)
        for _, perm := range permissions {
            if perm.Resource == resource && perm.Action == action {
                return true
            }
        }
    }
    return false
}
```

### 2.5 优缺点

| 优点 | 缺点 |
|------|------|
| 模型简单，易于理解 | 角色爆炸问题（角色过多） |
| 管理成本低 | 难以处理细粒度权限 |
| 支持角色继承 | 不支持动态条件判断 |

---

## 3. ABAC（基于属性的访问控制）

### 3.1 核心概念

ABAC 根据**属性**和**策略**动态计算权限：

```
权限 = f(用户属性，资源属性，环境属性，操作类型)
```

| 属性类型 | 示例 |
|----------|------|
| **用户属性** | 部门、职级、角色 |
| **资源属性** | 所有者、敏感级别、分类 |
| **环境属性** | 时间、地点、设备、IP |
| **操作类型** | read、write、delete |

### 3.2 策略示例

```json
{
  "effect": "allow",
  "subject": {
    "role": "manager",
    "department": "sales"
  },
  "resource": {
    "type": "report",
    "owner": "${subject.department}"
  },
  "action": ["read", "write"],
  "condition": {
    "time": "09:00-18:00",
    "ip_range": "10.0.0.0/8"
  }
}
```

**策略含义：** 销售部门的经理只能在工作时间内、从内网访问本部门的报告。

### 3.3 优缺点

| 优点 | 缺点 |
|------|------|
| 细粒度控制 | 策略复杂，难以管理 |
| 支持动态条件 | 性能开销大 |
| 灵活性强 | 学习和实现成本高 |

---

## 4. ReBAC（基于关系的访问控制）

### 4.1 核心概念

ReBAC 根据**实体之间的关系**定义权限：

```
用户 A 是 文档 X 的 编辑者 → 用户 A 可以编辑 文档 X
用户 A 是 项目 Y 的 成员 → 用户 A 可以访问 项目 Y 的资源
```

### 4.2 典型场景

| 关系 | 权限 |
|------|------|
| 文档所有者 | 完全控制 |
| 项目成员 | 查看项目内容 |
| 团队成员 | 编辑团队资源 |
| 好友关系 | 查看私人动态 |

### 4.3 与 RBAC 对比

| 维度 | RBAC | ReBAC |
|------|------|-------|
| 权限来源 | 角色 | 关系 |
| 适用场景 | 组织架构明确 | 社交网络、协作工具 |
| 管理方式 | 管理员分配角色 | 用户建立关系 |

---

## 5. ACL（访问控制列表）

### 5.1 核心概念

ACL 直接为每个资源维护一个访问控制列表：

```
资源 A 的 ACL:
  - 用户 1: read, write
  - 用户 2: read
  - 用户 3: read, write, delete
```

### 5.2 优缺点

| 优点 | 缺点 |
|------|------|
| 简单直接 | 资源多时列表过长 |
| 细粒度控制 | 难以批量管理 |
| 适合小系统 | 不适合多租户 SaaS |

---

## 6. 模型对比总结

| 维度 | RBAC | ABAC | ReBAC | ACL |
|------|------|------|-------|-----|
| 复杂度 | 低 | 高 | 中 | 低 |
| 灵活性 | 中 | 高 | 高 | 低 |
| 性能 | 好 | 一般 | 好 | 一般 |
| 易管理性 | 好 | 一般 | 好 | 差 |
| 适用场景 | 企业 SaaS | 高安全场景 | 社交/协作 | 小系统 |

---

## 7. IAM 系统选型

### 7.1 当前选择：**RBAC**

**理由：**
1. 符合企业组织架构（角色明确）
2. 管理简单，租户管理员可自助操作
3. 性能优秀，适合高并发场景
4. 社区成熟，有大量最佳实践

### 7.2 未来扩展

**RBAC + ABAC 混合模式：**

- 基础权限使用 RBAC
- 数据范围使用 ABAC（如只能查看本部门数据）
- 敏感操作增加环境属性校验（时间、IP）

**示例：**
```
用户角色 = 经理（RBAC）
         ↓
  可以访问报表系统
         ↓
  但只能查看本部门报表（ABAC：department 匹配）
         ↓
  且只能在工作时间访问（ABAC：time 条件）
```

---

## 8. 权限设计最佳实践

### 8.1 权限命名规范

```
<资源>:<操作>

示例:
  user:read      - 查看用户
  user:write     - 修改用户
  user:delete    - 删除用户
  order:read     - 查看订单
  order:approve  - 审批订单
```

### 8.2 角色设计原则

1. **角色数量可控**：建议 10-50 个角色
2. **角色职责单一**：一个角色代表一类职能
3. **避免角色重叠**：减少权限重复
4. **预置角色**：系统内置常用角色（Admin、Viewer）

### 8.3 权限分配流程

```mermaid
flowchart LR
    A[创建角色] --> B[配置权限]
    B --> C[分配给用户]
    C --> D[权限生效]
```

### 8.4 权限回收

- 删除角色时，通知受影响的用户
- 支持权限审计，查看谁有什么权限
- 支持临时权限（带过期时间）

---

## 9. 开源权限库参考

| 项目 | 语言 | 特点 |
|------|------|------|
| **Casbin** | Go | 支持 RBAC、ABAC，轻量级 |
| **Open Policy Agent** | Go | 通用策略引擎，功能强大 |
| **Permify** | Go | 细粒度授权服务 |
| **SpiceDB** | Go | Google Zanzibar 论文实现 |

---

## 10. 参考链接

- NIST RBAC 标准：https://csrc.nist.gov/projects/role-based-access-control
- Casbin 文档：https://casbin.org/
- Google Zanzibar 论文：https://research.google/pubs/pub48190/

---

## 11. 相关需求文档

- [REQ-005 角色管理功能](../05-functional-requirements/REQ-005-role-management.md)
- [REQ-006 权限分配功能](../05-functional-requirements/REQ-006-permission-assignment.md)
