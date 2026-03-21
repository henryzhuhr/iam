# 文档目录

`docs/` 目录用于维护 IAM 项目的需求、设计和规划类文档。

## 文档入口

| 文档 | 说明 | 入口 |
|------|------|------|
| PRD 一页摘要 | 几分钟内理解需求范围、优先级和主要风险 | [PRD/00-executive-summary.md](./PRD/00-executive-summary.md) |
| PRD | 产品需求文档(product requirement document)，按主题和 `REQ` 粒度拆分维护 | [PRD/README.md](./PRD/README.md) |
| REQUIREMENTS | 需求分析文档(requirements analysis document)，包含功能清单、API 设计、数据库设计等 | [REQUIREMENTS.md](./REQUIREMENTS.md) |

## 推荐阅读顺序

1. 先阅读 [PRD/00-executive-summary.md](./PRD/00-executive-summary.md)，快速了解产品范围、优先级和主要风险。
2. 再阅读 [PRD/README.md](./PRD/README.md)，查看完整 PRD 结构、模块导航和需求入口。
3. 最后阅读 [REQUIREMENTS.md](./REQUIREMENTS.md)，查看接口草案、数据库设计和更偏实现侧的分析内容。

## 当前结构

```bash
docs/
├── README.md
├── PRD # 产品需求文档目录
└── REQUIREMENTS.md
```

## 维护约定

- `PRD/` 目录作为产品需求主入口，不再保留单文件 `PRD.md`。
- 新增功能需求时，优先在 `PRD/05-functional-requirements/` 下按 `REQ` 粒度维护。
- 偏接口、数据结构、实现约束的补充内容，优先更新 `REQUIREMENTS.md`。
