# Go Checklist 问题清单

本文档只列出当前 checklist 中我认为值得你判断和决策的争议点，不直接修改原 checklist。

## 使用方式

- 先看“问题类型”，判断这是不是你们团队有意设定的约束。
- 如果是团队显式约束，可以保留，但建议在正文里写明“这是团队约定，不是通用 Go 规则”。
- 如果不是团队显式约束，建议降级为“建议”或改写为“结合上下文判断”。

## `go-checklist.md`

- 文件位置：`.agents/skills/go-code-review/references/go-checklist.md`
- 问题：索引文档只做主题分发，没有说明哪些条目属于高风险正确性问题，哪些只是风格或团队偏好。
- 影响：review 使用者容易把子文档中的所有“强制”条目等同于通用 Go 硬规则。
- 建议决策：
  - 保持现状，接受这套 checklist 本身带有强规则倾向。
  - 在索引里补充分层说明，区分“通用风险”“仓库约定”“建议项”。

## `go-checklist-01-style.md`

- 文件位置：`.agents/skills/go-code-review/references/go-checklist-01-style.md`
- 问题：`01.2 文件长度` 把“单文件尽量不超过 800 行”写成通用建议，但这个阈值没有统一依据。
- 影响：容易把组织偏好当作质量结论。
- 建议决策：
  - 如果团队确实有规模阈值，明确写成团队约定。
  - 否则改成“文件过长时关注职责混杂和 review 成本”。

- 问题：`01.4 标识符命名` 中“错误变量以 err/Err 开头，错误类型以 Error 结尾”写成强制。
- 影响：这不是 Go 的硬规则，更多是惯例；不少项目会使用其他同样清晰的命名。
- 建议决策：
  - 保留为建议。
  - 或改成“优先与包内既有风格一致”。

- 问题：`01A.1 通用规则` 与 `01A.2 声明注释` 对导出名字、类型别名、全局变量等注释要求过于一刀切。
- 影响：容易鼓励为了补注释而补注释，增加噪音和失真风险。
- 建议决策：
  - 保留对公共 API 的高要求。
  - 降低对内部实现和显而易见声明的刚性要求。

## `go-checklist-02-error-resource.md`

- 文件位置：`.agents/skills/go-code-review/references/go-checklist-02-error-resource.md`
- 问题：`02.1 返回错误` 中“接收 context.Context 的函数通常也应返回 error”过于泛化。
- 影响：有些函数只是查询状态、等待信号、计算派生值，并不天然需要返回 `error`。
- 建议决策：
  - 改成“执行可能失败的 IO、RPC、存储、编解码时通常返回 error”。

- 问题：`02.5 类型断言` 把 `v, ok := ...` 写成强制。
- 影响：在不变量已经严格成立的分支里，直接断言是合理的，机械禁止会降低表达力。
- 建议决策：
  - 改成“动态类型不确定时优先使用 ok 形式”。

- 问题：`02.6 panic / recover` 中“panic 必须在当前 goroutine 内捕获”“自行启动的 goroutine 要考虑 recover”过于绝对。
- 影响：有些系统希望后台 panic 直接暴露并终止进程，而不是吞掉。
- 建议决策：
  - 改成“要明确该 goroutine 的 panic 策略：快速失败还是隔离失败”。

- 问题：`02A.2 内存泄漏` 里关于 subslice 和 map 收缩的说法更接近性能/生命周期建议，不是一般性的确定缺陷。
- 影响：容易在没有 profile、没有长生命周期前提时过度报问题。
- 建议决策：
  - 降级为“热点路径或长生命周期对象需关注”。

## `go-checklist-03-project-function.md`

- 文件位置：`.agents/skills/go-code-review/references/go-checklist-03-project-function.md`
- 问题：`03.1 单元测试` 把“项目应提供单元测试”写成强制。
- 影响：有些仓库主要依赖集成测试、端到端测试、golden test 或契约测试，单元测试不是唯一正确形态。
- 建议决策：
  - 改成“行为变化应有合适层级的测试覆盖”。

- 问题：`03A.1 参数` 中“不要给 map、slice、chan、interface 传指针”写成普适建议，但少数场景确实会用到。
- 影响：这会把少见但合理的设计一律视为反模式。
- 建议决策：
  - 改成“除非需要重绑定、区分 nil 状态或修改头部结构，否则一般不需要指针”。

- 问题：`03A.2 返回值` 中“返回值个数不超过 3 个”写成强制。
- 影响：这是典型的人为阈值，不是 Go 语义规则。
- 建议决策：
  - 改成“当返回值过多且语义不清时，考虑封装结构体或领域对象”。

- 问题：`03A.4 分组` 是文件排布偏好，不适合被误读成 review 的硬性结论。
- 影响：容易产出低价值风格评论。
- 建议决策：
  - 保留为低优先级建议。
  - 或只在导航成本明显升高时再提示。

## `go-checklist-04-language-data.md`

- 文件位置：`.agents/skills/go-code-review/references/go-checklist-04-language-data.md`
- 问题：`04.2 语句` 中“禁止变量遮蔽”写成强制。
- 影响：变量遮蔽本身不是必错；真正有问题的是会改变行为、误导读者或掩盖 bug 的遮蔽。
- 建议决策：
  - 改成“重点关注危险遮蔽，尤其是 err、命名返回值、循环变量”。

- 问题：`04.3 魔数` 中“除 0 和 1 外，避免直接使用魔数”过于绝对。
- 影响：会迫使局部、显而易见的字面量都被抽成常量，反而降低可读性。
- 建议决策：
  - 改成“重复出现或具有业务语义的字面量提取常量”。

- 问题：`04.4 for / range / switch / goto` 中“只需要 value 时把 key 写成 _”不准确。
- 影响：Go 里 `for _, v := range xs` 常见且合理，但 `for v := range xs` 在 slice/array/string 上语义变成索引，不是“只取 value”。
- 建议决策：
  - 明确分别说明 map、slice、array、string 的 range 语义。

- 问题：`04.4 for / range / switch / goto` 中“switch 必须有 default”写成强制。
- 影响：穷举 enum、显式依赖编译期检查或希望未来新增值暴露问题时，不写 default 反而更合适。
- 建议决策：
  - 改成“非法值和未来扩展值要有明确策略”。

- 问题：`04.4 for / range / switch / goto` 中“业务代码禁止使用 goto”写成强制。
- 影响：`goto` 很少需要，但并非绝对错误；某些错误清理路径会用它简化控制流。
- 建议决策：
  - 改成“除非能明显降低复杂度，否则避免使用 goto”。

- 问题：`04.5 结构体` 中“字段排序考虑内存对齐”写成强制。
- 影响：这更多是性能/内存优化议题，不是普适 correctness 规则。
- 建议决策：
  - 仅在大对象、热点路径或高实例量时强调。

- 问题：`04.6 slice` 中“零长度返回值优先 nil，空判断统一 len(s)==0”过于绝对。
- 影响：`nil` slice 和空 slice 在 API、JSON、数据库扫描、前端契约里语义可能不同。
- 建议决策：
  - 改成“以 API 语义和序列化约定为准，保持一致”。

- 问题：`04.7 map` 中“固定元素列表直接用字面量，其他场景用 make”“可预估容量时显式给容量”写成强制。
- 影响：这属于风格或性能建议，通常不应上升为硬规则。
- 建议决策：
  - 降级为建议。

- 问题：`04.7 map` 中“判断 key 是否存在要用 v, ok := m[k]”写成强制。
- 影响：当零值和不存在在业务上等价时，没有必要引入 `ok`。
- 建议决策：
  - 改成“只有在需要区分不存在与零值时使用 ok”。

## `go-checklist-05-concurrency-context.md`

- 文件位置：`.agents/skills/go-code-review/references/go-checklist-05-concurrency-context.md`
- 问题：`05.1 channel` 中“channel 通常应无缓冲或容量为 1，其他容量要解释”本身就不正确。
- 影响：缓冲大小取决于协议、背压、批处理和吞吐模型，不存在这种通用 Go 规则。
- 建议决策：
  - 必须改写。

- 问题：`05.2 goroutine` 中“goroutine 必须可等待退出或取消，不能泄漏”写成强制，但“是否等待”要看语义。
- 影响：有些后台 goroutine 就是伴随进程存活，不需要被上层 join。
- 建议决策：
  - 改成“goroutine 需要清楚的生命周期和退出条件”。

- 问题：`05.2 goroutine` 中“协程池要限制最大并发数”过于泛化。
- 影响：是否限流取决于任务来源、资源消耗和业务模型，不是所有并发执行都需要池化。
- 建议决策：
  - 改成“无界并发需要重点审查”。

- 问题：`05.3 context` 中“context.Context 始终为第一个参数”“不要传 nil context”混合了通用惯例和合理例外。
- 影响：新 API 通常如此，但实现既有接口、兼容第三方签名时不能机械套用。
- 建议决策：
  - 改成“新设计 API 通常放首参；既有接口以接口契约为准”。

- 问题：`05.5 atomic / rand / unsafe / cgo` 中“禁止随意使用 unsafe”表述太模糊。
- 影响：问题不在于用了 `unsafe`，而在于是否缺少边界、不变量和测试。
- 建议决策：
  - 改成“unsafe 需要明确封装边界、不变量说明和额外验证”。

## `review-checklist.md`

- 文件位置：`.agents/skills/go-code-review/references/review-checklist.md`
- 问题：这里的“常见 Go 特有问题”里有几条与主 checklist 一样偏绝对化，例如“在循环里使用 defer”“在 goroutine 或闭包里捕获循环变量”。
- 影响：容易被读成“只要出现就是问题”，而不是“在特定条件下会成为问题”。
- 建议决策：
  - 与主 checklist 统一口径，改成带前提条件的表述。

## 优先级建议

- 高优先级：`go-checklist-05-concurrency-context.md`
- 高优先级：`go-checklist-04-language-data.md`
- 中优先级：`go-checklist-03-project-function.md`
- 中优先级：`go-checklist-02-error-resource.md`
- 低优先级：`go-checklist-01-style.md`
- 低优先级：`go-checklist.md`
- 低优先级：`review-checklist.md`
