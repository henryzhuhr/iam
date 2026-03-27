# Agent Conventions

## 基本原则

- 当你在查看项目文档的时候，始终读取 [docs/README.md](./docs/README.md) 以了解文档结构（不同文档的入口）和维护原则。
- 如果一个目录里有 `README.md`，优先阅读它以了解该目录的内容和结构。
- `AGENTS.md` 仅保留文档索引和最基本原则；项目规范、代码规范、目录约定等内容统一维护在 `docs/` 中。

## 文档索引

- 项目文档总入口：[docs/README.md](./docs/README.md)
- 项目规范与代码规范：[docs/project-conventions.md](./docs/project-conventions.md)
- 产品需求文档入口：[docs/PRD/README.md](./docs/PRD/README.md)
- 技术参考资料入口：[docs/references/README.md](./docs/references/README.md)

## Issues

### 目录和文件命名规范

- 整个项目统一使用名为 `issues/` 的目录记录 issue，不限定于 skill。
- issue 文件名使用三位递增编号开头，格式为 `NNN-short-kebab-case.md`，例如 `001-path-name-collision.md`。
- 新 issue 必须延续当前最大编号，不能跳号，也不要重用已有编号。
- 新增 issue 时，需要同步在对应层级的 `README.md` 中维护 index，方便按编号查阅。
