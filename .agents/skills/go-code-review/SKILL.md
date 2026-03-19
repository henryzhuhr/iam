---
name: go-code-review
description: 审查 Go 代码中的正确性、并发、安全性、错误处理、测试、性能和可维护性问题，并结合工具检查确认可疑点。适用于 Codex 被要求 review Go diff、PR、文件、包或整个仓库，定位 bug、回归和设计风险，或借助 go test、go vet、gofmt、goimports、staticcheck、golangci-lint、ineffassign、revive、govulncheck 对发现进行验证时。
---

# Go 代码 Review

## 概述

以发现缺陷和回归风险为目标审查 Go 代码，而不是只做风格点评。把对 diff 和上下文代码的推理，与静态分析、lint、测试输出结合起来，让结论更具体、更可验证。

需要更完整的检查清单时，先读取 `references/go-checklist.md`，再按主题按需加载：

- 风格、命名、注释：`references/go-checklist-01-style.md`
- 错误处理、资源释放、panic：`references/go-checklist-02-error-resource.md`
- 测试、项目组织、函数设计：`references/go-checklist-03-project-function.md`
- 变量、控制流、结构体、slice、map、序列化：`references/go-checklist-04-language-data.md`
- channel、goroutine、context、atomic、Mutex：`references/go-checklist-05-concurrency-context.md`

需要补充工具结果解读时，再读取 `references/review-checklist.md`。

## 工作流

1. 明确 review 范围。
先判断用户要的是 diff review、单文件 review，还是包级别或仓库级别审计。优先选择足以支撑结论的最小范围。

2. 先补上下文，再下结论。
阅读改动代码、调用方、接口、测试和配置，确认数据流如何经过 handler、service、goroutine 和错误分支，再报告问题。

3. 执行工具验证。
从仓库根目录运行 `scripts/run-go-review-checks.sh`。如果范围较小，就传入受影响的目录或包；否则默认检查整个模块。工具输出是证据，不是替代思考。

4. 区分已证实问题和推断性问题。
如果工具直接证明了问题，要明确说明。若问题主要来自代码推理而工具未报错，也要如实标记为基于推理的发现，不要夸大确定性。

5. 先输出 findings。
按严重程度排序，并附上文件引用。重点关注 bug、回归、竞态风险、错误处理缺陷、契约不一致、测试缺口和运维风险，摘要保持简洁。

## Review 优先级

- 先检查正确性：nil 处理、错误传播、分支条件、返回值是否忽略、变量遮蔽、状态不变式是否被破坏。
- 再检查并发：goroutine 生命周期、channel 所有权、加锁方式、context 取消传播、共享状态上的数据竞争。
- 再看外部行为：API 兼容性、配置默认值、HTTP 状态码、序列化 tag、数据库边界和向后兼容性。
- 检查测试与可观测性：改动逻辑是否缺少测试覆盖、测试是否脆弱、失败是否被静默吞掉、日志和指标是否足够定位问题。
- 最后看可维护性：重复逻辑、误导性命名、死代码、不可能分支和不必要的抽象复杂度。

## 工具使用建议

- 把 `go test` 和 `go vet` 当作基础检查。
- 用 `staticcheck` 捕获 `go vet` 之外的正确性和 API 误用问题。
- 如果仓库安装了 `golangci-lint`，优先使用，因为项目规则和额外分析器通常已经配置在里面。
- `gofmt -l` 和 `goimports -l` 适合检查未格式化或 import 未整理的改动，但纯格式问题通常不应上升为高严重级别。
- `ineffassign`、`revive`、`govulncheck` 作为可选补充；如果本机没有这些工具，要在结果里说明它们未执行。

## 输出格式

- 先列 findings，并按严重程度排序。
- 每条问题写清影响和具体代码路径。
- 每条问题附上文件引用。
- 说明哪些检查已执行，哪些因为工具缺失或范围原因被跳过。
- 如果没有发现问题，要明确说明，并补充剩余风险或未完成的验证项。

## 命令

从仓库根目录运行内置检查脚本：

```bash
./.agents/skills/go-code-review/scripts/run-go-review-checks.sh
./.agents/skills/go-code-review/scripts/run-go-review-checks.sh ./internal/... ./app
```

如果用户要求 review 特定 diff，先读 diff，再只对受影响的目录或包执行检查，不要无差别扩大到无关模块。
