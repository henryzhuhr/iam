# 文档目录

`docs/` 目录用于维护 IAM 项目的需求、设计、测试方案和规划类文档。

## 文档入口

| 文档分类 | 说明 | 维护角色 | 入口 |
|----------|------|----------|------|
| **PRD** | 产品需求文档，按主题和 `REQ` 粒度拆分维护 | 产品经理 | [PRD/README.md](./PRD/README.md) |
| **standards/** | 项目规范，包含项目结构、编码规范、Git 工作流、术语表、API/数据库设计规范 | 研发团队 | [standards/README.md](./standards/README.md) |
| **references/** | 参考资料，收录 OAuth、OIDC、JWT、MFA、RBAC、OWASP 等技术概念说明 | 研发团队 | [references/README.md](./references/README.md) |
| **issues/** | Issue 跟踪目录，记录项目开发过程中的所有 Issue | 全体研发 | [issues/README.md](./issues/README.md) |
| **TDD** | 技术设计文档，开发视角的架构与实现设计 | 研发团队 | **[待设计]** |

查看各目录的 `README.md` 获取详细文档列表：

```bash
docs/
├── standards/      # 项目规范目录
├── PRD/            # 产品需求文档目录
├── TDD/            # 技术设计文档目录 [待设计]
├── issues/         # Issue 跟踪目录
└── references/     # 技术资料参考
```

## 推荐阅读顺序

### 入门路径

1. [PRD/README.md](./PRD/README.md) - 查看完整 PRD 结构、模块导航和需求入口
2. [standards/README.md](./standards/README.md) - 了解项目结构、编码约定和 issue 规范
3. [references/README.md](./references/README.md) - 查阅技术概念和协议标准

### 待补充

- [TDD/README.md](./TDD/README.md) - 技术设计文档，开发视角的架构与实现设计

## 维护约定

- 新增功能需求时，优先在 `PRD/05-functional-requirements/` 下按 `REQ` 粒度维护。
- 偏接口、数据结构、实现约束的补充内容，优先更新对应 PRD 文档或系统架构文档。
- 项目结构、技术栈、编码规范、issue 规范等研发约定，统一维护在 `standards/` 目录。
- 技术概念、协议标准、设计方案等参考资料，收录到 `references/` 目录。

### Issues 维护规范

- 整个项目统一使用名为 `issues/` 的目录记录 issue，不限定于 skill。
- issue 文件名使用三位递增编号开头，格式为 `NNN-short-kebab-case.md`，例如 `001-path-name-collision.md`。
- 新 issue 必须延续当前最大编号，不能跳号，也不要重用已有编号。
- 新增 issue 时，需要同步在对应层级的 `README.md` 中维护 index，方便按编号查阅。
