# Go 代码审查清单 01: 风格与注释

本文件聚焦格式、命名、注释和可维护性。这里的大多数条目都不是“只要不满足就一定是 bug”，而是帮助 reviewer 判断代码是否违背仓库约定、显著增加理解成本，或掩盖真实行为。

## 01. 风格规范

### 01.1 先看仓库约定，不要输出个人偏好

- 优先尊重仓库已经采用的 `gofmt`、`goimports`、lint、文件组织和命名习惯。
- 风格问题只有在影响阅读、违反 CI 或明显偏离同目录代码时，才值得在 review 里强调。

### 01.2 格式与工具链

- `gofmt` 是 Go 代码的基线；格式漂移通常是低严重度问题，但如果仓库依赖格式化检查或生成代码，它就会变成实际问题。
- 编码、换行符、文件头等要求更适合作为工具链约束；只有在会破坏编译、测试、CI 或协作时才需要上升到 review 发现。

### 01.3 命名

- 名字应体现语义、作用域和层次，避免误导性缩写、重复包名语义和“看起来像容器实现细节”的噪音命名。
- receiver 名称应短且一致，但这是包内约定，不必机械要求所有类型都统一成单字母。
- 导出的哨兵错误常见写法是 `ErrXxx`，导出的错误类型常见写法是 `XxxError`，但它们是惯例，不是绝对规则；优先看包内现有风格是否一致。

反例：

```go
filterHandlerMap
uidSlice
uidArray
```

更清晰的写法：

```go
opToHandler
uids
```

### 01.4 文件和函数规模

- 文件或函数过长本身不是问题，真正要看的是职责是否混杂、分支是否难以验证、局部改动是否需要理解过多无关逻辑。
- 如果长文件只是稳定的数据定义或规则表，未必值得报；如果一个函数同时做 IO、鉴权、转换、持久化和日志拼装，就值得指出拆分机会。

## 01A. 注释

### 01A.1 注释的目的

- 注释应该补充代码里看不出来的信息，例如约束、单位、默认值、并发假设、错误语义、兼容性承诺和外部协议。
- 不要为了“有注释”而重复代码表面意思，这类注释会快速失真。

### 01A.2 什么时候需要注释

- 公共包、导出 API、跨团队使用的类型通常应提供符合 godoc 的注释。
- 内部实现如果语义直接、名字清楚，可以不写声明性注释；真正需要注释的是非显而易见的决策点和坑点。
- `//` 和 `/* ... */` 都可以使用，关键是能否清楚表达意图并与仓库现有风格一致。

优秀示例：

```go
// HasPrefix reports whether name matches any prefix in prefixes.
// prefixes is expected to be small; callers in hot paths should prebuild an index.
func HasPrefix(name string, prefixes []string) bool {
	// ...
}
```

优秀示例：

```go
type User struct {
	Username string // Login name shown to other users.
	Email    string // Verified primary email.
}
```
