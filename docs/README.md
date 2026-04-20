# 文档目录

`docs/` 目录用于维护 IAM 项目的需求、设计、测试方案和参考文档。

## 文档入口

| 文档分类 | 说明 | 维护角色 | 入口 |
|----------|------|----------|------|
| **PRD** | 产品需求文档，按主题和 `REQ` 粒度拆分维护 | 产品经理 | [PRD/README.md](./PRD/README.md) |
| **standards/** | 项目规范，包含项目结构、编码规范、Git 工作流、术语表、API/数据库设计规范 | 研发团队 | [standards/README.md](./standards/README.md) |
| **references/** | 参考资料，收录 OAuth、OIDC、JWT、MFA、RBAC、OWASP 等技术概念说明 | 研发团队 | [references/README.md](./references/README.md) |
| **issues/** | 项目级 Issue 跟踪目录，记录开发过程中的问题与待办 | 全体研发 | [../issues/README.md](../issues/README.md) |
| **TDD** | 技术设计文档，开发视角的架构与实现设计 | 研发团队 | [TDD/README.md](./TDD/README.md) |

查看各目录的 `README.md` 获取详细文档列表：

```bash
docs/
├── standards/      # 项目规范目录
├── PRD/            # 产品需求文档目录
├── TDD/            # 技术设计文档目录（包括前端、后端、测试方案）
└── references/     # 技术资料参考
```

项目级 Issue 目录位于仓库根目录：[`issues/README.md`](../issues/README.md)

## 需求开发标准化流程

主流程：`阅读索引 -> 更新 PRD -> 做影响分析 -> 更新 TDD 与测试设计 -> 记录 issues -> 开发与测试`

新增需求时，优先按下面的顺序推进，并同步更新对应文档：

1. 阅读索引并定位上下文：先阅读当前文档和相关子目录的 `README.md`，再按需阅读相关 `PRD`、`TDD`、`standards` 与 `references` 文档，不要求一次性通读全部文档。
2. 明确需求并维护 `PRD/`：确认需求目标、范围、优先级、验收标准和约束条件。已有对应 `REQ` 时优先增量更新；没有时再新增需求文档。
3. 做影响分析：检查新需求是否与已有需求、系统架构、接口约束或安全规则冲突；如有冲突或信息缺失，应先澄清再继续。
4. 更新 `TDD/` 并同步做测试设计：需求涉及实现方案、模块拆分、接口变更、数据结构时，必须同步更新对应设计文档；同时补充测试策略、测试范围、关键场景和必要的测试用例，确保设计与验证方案一起落地。测试体系相关内容可按需补充到 `docs/testing/`。
5. 记录开发跟踪项：如果需求在实现过程中还需要进一步拆解、排期、缺陷跟踪或风险记录，统一在仓库根目录 `issues/` 下维护。
6. 进入开发与测试：实现前再次核对 `PRD`、`TDD` 和测试设计是否完整，开发过程中如发生设计偏移，需同步回写文档，避免需求、设计、测试与代码脱节。

顶层 `docs/README.md` 只保留流程摘要和目录入口。更细粒度的规范，应维护在各自目录下的 `README.md` 或具体规范文档中。

## 推荐阅读顺序

### 入门路径

1. [PRD/README.md](./PRD/README.md) - 查看完整 PRD 结构、模块导航和需求入口
2. [standards/README.md](./standards/README.md) - 了解项目结构、编码约定和 issue 规范
3. [references/README.md](./references/README.md) - 查阅技术概念和协议标准

## 维护约定

- 新增功能需求时，优先在 `PRD/05-functional-requirements/` 下按 `REQ` 粒度维护。
- 偏接口、数据结构、实现约束的补充内容，优先更新对应 PRD 文档或系统架构文档。
- 项目结构、技术栈、编码规范等研发约定，统一维护在 `standards/` 目录。
- 技术概念、协议标准、设计方案等参考资料，收录到 `references/` 目录。
- Issue 跟踪规则统一维护在仓库根目录 [`issues/README.md`](../issues/README.md)。

### Issues 维护规范

- 整个项目统一使用名为 `issues/` 的目录记录 issue，不限定于 skill。
- issue 文件名使用三位递增编号开头，格式为 `NNN-short-kebab-case.md`，例如 `001-path-name-collision.md`。
- 新 issue 必须延续当前最大编号，不能跳号，也不要重用已有编号。
- 新增 issue 时，需要同步在对应层级的 `README.md` 中维护 index，方便按编号查阅。
