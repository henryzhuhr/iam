# Go 代码审查清单 01: 风格与注释

本文件聚焦风格规范、命名和注释，包含适合 review 场景的代表性反例和优秀示例。

## 01. 风格规范

### 01.1 遵从惯例

- 【建议】遵从 Go 社区惯例，不要引入反直觉写法。

### 01.2 文件长度

- 【建议】单文件尽量不超过 800 行；超长文件通常意味着职责混杂，review 成本高。

### 01.3 缩进、括号、空格

- 【强制】统一使用 `gofmt`。
- 【强制】tab 缩进，左大括号不换行。

优秀示例：

```go
var i int = 1 + 2
v := []float64{1.0, 2.0, 3.0}[i-i]
fmt.Printf("%f\n", v+1)
```

### 01.4 标识符命名

- 【强制】标识符遵循 `MixedCaps/mixedCaps`。
- 【强制】导出标识符大写开头，非导出标识符小写开头。
- 【强制】错误变量以 `err/Err` 开头，错误类型以 `Error` 结尾。
- 【强制】receiver 命名应短、稳定，不用 `this`、`self`。
- 【强制】标识符不要重复包名语义。

反例：

```go
filterHandlerMap
uidSlice
uidArray
uidSliceSlice
```

优秀示例：

```go
opToHandler
uids
uids
classesUids
```

优秀示例：

```go
type ExitError struct {
	// ...
}

func (ri *ResearchInfo) Load() {}
func (w *ReportWriter) Write() {}
```

### 01.5 编码格式

- 【强制】代码文件必须为 UTF-8。
- 【强制】换行符使用 Unix 风格。

## 01A. 注释

### 01A.1 通用规则

- 【建议】单行注释与内容之间留一个空格。
- 【建议】导出名字要有注释。
- 【建议】优先使用 `//`，`/* ... */` 只用于包级文档。
- 【建议】注释符合 godoc。

### 01A.2 声明注释

- 【建议】包、函数、方法、结构体、接口、全局变量、类型别名都应有清晰注释。

优秀示例：

```go
// NewAttrModel，属性数据层操作类的工厂方法。
// 返回值：属性操作类指针。
func NewAttrModel(ctx *common.Context) *AttrModel {
	// ...
}

// HasPrefix 返回 true，如果 name 包含指定 prefix。
func HasPrefix(name string, prefixes []string) bool {
	// ...
}
```

优秀示例：

```go
// User，用户实例，定义了用户的基础信息。
type User struct {
	Username string // 用户名
	Email    string // 邮箱
}
```
