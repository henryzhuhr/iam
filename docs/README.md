# 文档目录

`docs/` 目录用于维护 IAM 项目的需求、设计和规划类文档。

## 文档入口

| 文档 | 说明 | 入口 |
|------|------|------|
| PRD 一页摘要 | 几分钟内理解需求范围、优先级和主要风险 | [PRD/00-executive-summary.md](./PRD/00-executive-summary.md) |
| PRD | 产品需求文档 (product requirement document)，按主题和 `REQ` 粒度拆分维护 | [PRD/README.md](./PRD/README.md) |
| TDD | 技术设计文档 (technical design document)，开发视角的架构与实现设计 | **[待设计]** |
| 项目规范 | 项目结构、技术栈、编码约定、Issue 规范等研发约定 | [project-conventions.md](./project-conventions.md) |
| REFERENCES | 参考资料，收录 OAuth、权限模型、JWT、MFA、多租户架构等技术概念说明 | [references/README.md](./references/README.md) |

## 推荐阅读顺序

1. 先阅读 [PRD/00-executive-summary.md](./PRD/00-executive-summary.md)，快速了解产品范围、优先级和主要风险。
2. 再阅读 [PRD/README.md](./PRD/README.md)，查看完整 PRD 结构、模块导航和需求入口。
3. 如需了解项目结构、编码约定和 issue 规范，请阅读 [project-conventions.md](./project-conventions.md)。
4. 技术概念查阅请参考 [references/README.md](./references/README.md)。

**待补充：**
- 技术设计文档 (TDD) - 开发视角的架构与实现设计

## 当前结构

```bash
docs/
├── README.md
├── project-conventions.md  # 项目规范
├── PRD/                    # 产品需求文档目录
│   ├── README.md
│   ├── 00-executive-summary.md
│   ├── 01-document-overview.md
│   ├── 02-product-overview.md
│   ├── 03-system-architecture.md
│   ├── 04-requirements-overview.md
│   ├── 05-functional-requirements/  # 功能需求详情 (按 REQ 拆分)
│   ├── 06-non-functional-requirements.md
│   ├── 07-business-flows.md
│   ├── 08-risks-and-dependencies.md
│   ├── 09-success-metrics.md
│   └── 10-appendix-user-stories.md
├── TDD/                    # 技术设计文档目录 [待设计]
│   └── README.md           # 待创建
└── references/             # 技术资料参考
    └── README.md
```

## 维护约定

- `PRD/` 目录作为产品需求主入口，不再保留单文件 `PRD.md`。
- 新增功能需求时，优先在 `PRD/05-functional-requirements/` 下按 `REQ` 粒度维护。
- 偏接口、数据结构、实现约束的补充内容，优先更新对应 PRD 文档或系统架构文档。
- 项目结构、技术栈、编码规范、issue 规范等研发约定，统一维护在 `project-conventions.md`。
- 技术概念、协议标准、设计方案等参考资料，收录到 `references/` 目录。
- **技术设计文档 (TDD)** - 开发视角的架构与实现设计，待创建。
