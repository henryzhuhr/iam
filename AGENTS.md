# Agent Conventions

## 基本原则

- 始终读取 [docs/README.md](./docs/README.md) 以了解文档结构（不同文档的入口）和维护原则。
- 文档是逐层展开的。每个目录下的 `README.md` 是该层级的总览和索引，指向更细粒度的文档。例如当前的 `AGENTS.md` 是项目级的文档索引，指向 `docs/` 目录下的各类文档入口 `README.md`，但是不要负责维护具体内容。
- 如果一个目录里有 `README.md`，优先阅读它以了解该目录的内容和结构。
- `AGENTS.md` 仅保留文档索引和最基本原则；项目规范、代码规范、目录约定等内容统一维护在 `docs/` 中。
- 始终读取 `docs/README.md` 以确定是否需要查看更细粒度的文档，而不是直接在 `AGENTS.md` 中查找具体内容。

## 代码提交

当前 Git 历史采用 Conventional Commit 风格，建议优先使用 `feat:`、`fix:`、`docs:`、`refactor:` 前缀，后面附带中文的提交信息。

例如：

```bash
# <emoji> <type>: 中文的提交信息
✨ feat: 添加什么功能
🐛 fix: 修复了什么问题
```
