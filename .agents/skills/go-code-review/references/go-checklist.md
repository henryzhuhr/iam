# Go 代码审查清单索引

按需读取下列分文件，而不是一次性加载全部内容：

- `go-checklist-01-style.md`：适用于命名、格式、注释、基础风格相关问题。
- `go-checklist-02-error-resource.md`：适用于错误处理、日志、panic/recover、资源关闭、内存泄漏相关问题。
- `go-checklist-03-project-function.md`：适用于项目组织、测试、包导入、函数参数和返回值设计。
- `go-checklist-04-language-data.md`：适用于变量声明、控制流、结构体、slice、map、序列化、time 等语言细节。
- `go-checklist-05-concurrency-context.md`：适用于 channel、goroutine、context、atomic、Mutex、随机数、unsafe、cgo。

使用建议：

- 做全仓 review 时，先读 `go-checklist-02-error-resource.md` 和 `go-checklist-05-concurrency-context.md`。
- 如果问题偏 API、数据结构、默认值、tag、边界条件，再读 `go-checklist-04-language-data.md`。
- 如果问题偏风格和可维护性，再读 `go-checklist-01-style.md` 与 `go-checklist-03-project-function.md`。
- 需要工具结果解释时，再补充读取 `review-checklist.md`。
