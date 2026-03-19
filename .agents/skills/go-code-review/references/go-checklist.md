# Go 代码审查清单索引

这组清单用于 code review，不是通用编码法典。重点是帮助 reviewer 发现会导致错误行为、资源泄漏、并发缺陷、兼容性回归和维护成本失控的问题，而不是把所有个人偏好都上升为“必须遵守”。

使用这组清单前，先把问题分成三类：

- 高风险问题：会直接影响正确性、可靠性、安全性、兼容性或资源生命周期。这类问题即使仓库没写规范，也值得优先指出。
- 仓库约定：`gofmt`、`goimports`、lint、目录组织、命名风格等。优先尊重现有代码、CI 和团队约定，而不是套用通用模板。
- 上下文相关建议：性能优化、抽象层次、注释密度、字段对齐、预分配容量等。这类问题要结合热点路径、公开 API 和维护成本判断，不要机械判定。

按需读取下列分文件，而不是一次性加载全部内容：

- `go-checklist-01-style.md`：格式、命名、注释和可维护性提示。
- `go-checklist-02-error-resource.md`：错误传播、日志、panic/recover、资源释放、内存保留。
- `go-checklist-03-project-function.md`：测试、包边界、函数契约、参数和返回值设计。
- `go-checklist-04-language-data.md`：变量作用域、控制流、结构体、slice、map、序列化、时间语义。
- `go-checklist-05-concurrency-context.md`：channel、goroutine、context、锁、atomic、随机数、unsafe、cgo。

使用建议：

- 做全仓 review 或关键路径 review 时，先读 `go-checklist-02-error-resource.md` 和 `go-checklist-05-concurrency-context.md`。
- 如果改动涉及 API、默认值、序列化、时间处理或边界条件，再读 `go-checklist-04-language-data.md`。
- 如果改动涉及包边界、测试和函数契约，再读 `go-checklist-03-project-function.md`。
- 只有在问题已经影响可读性、维护性或违反仓库约定时，再重点看 `go-checklist-01-style.md`。
- 需要解释工具结果时，再补充读取 `review-checklist.md`。

判断原则：

- 先问“这段代码会不会出错、泄漏、阻塞、竞态或破坏兼容性”。
- 再问“它是否违反了仓库已有约定、工具链或公共 API 契约”。
- 最后才问“这是不是只是另一种也合理的写法”。
