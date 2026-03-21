# 文档目录

`docs/` 目录用于维护 IAM 项目的需求、设计、测试方案和规划类文档。

## 文档入口

| 文档 | 说明 | 入口 |
|------|------|------|
| PRD | 产品需求文档(product requirement document)，按主题和 `REQ` 粒度拆分维护 | [PRD/README.md](./PRD/README.md) |
| REQUIREMENTS | 需求分析文档(requirements analysis document)，包含功能清单、API 设计、数据库设计等 | [REQUIREMENTS.md](./REQUIREMENTS.md) |
| TESTING | 测试文档入口，索引项目测试框架设计、用例编写指南、架构说明、测试指标与结果说明 | [testing/README.md](./testing/README.md) |

## 推荐阅读顺序

1. 先阅读 [PRD/README.md](./PRD/README.md)，了解产品背景、需求总览和各模块入口。
2. 再阅读 [REQUIREMENTS.md](./REQUIREMENTS.md)，查看接口草案、数据库设计和更偏实现侧的分析内容。
3. 如果要补充或维护接口测试，阅读 [testing/README.md](./testing/README.md)。

## 当前结构

```bash
docs/
├── README.md
├── PRD # 产品需求文档目录
├── testing/
│   ├── README.md
│   └── testing-framework-design-record.md
└── REQUIREMENTS.md
```

## 维护约定

- `PRD/` 目录作为产品需求主入口，不再保留单文件 `PRD.md`。
- 新增功能需求时，优先在 `PRD/05-functional-requirements/` 下按 `REQ` 粒度维护。
- 偏接口、数据结构、实现约束的补充内容，优先更新 `REQUIREMENTS.md`。
- 偏测试框架、测试目录约定、运行方式、测试用例编写、测试架构、测试指标和测试结果说明的内容，优先更新 `docs/testing/` 目录下对应文档。
