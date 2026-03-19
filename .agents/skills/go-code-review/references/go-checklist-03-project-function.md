# Go 代码审查清单 03: 项目组织与函数设计

本文件聚焦测试、模块化、包导入，以及函数参数、返回值和组织方式。

## 03. 项目组织

### 03.1 单元测试

- 【强制】项目应提供单元测试。
- 【强制】并发相关代码测试时应考虑 `go test -race`。

### 03.2 模块化

- 【建议】项目应包含 `go.mod`，模块名、包名和目录结构保持清晰一致。

### 03.3 包导入

- 【强制】导入分组、别名和空白导入要合理，避免未使用导入和无意义别名。

## 03A. 函数

### 03A.1 参数

- 【建议】相同类型参数相邻。
- 【建议】优先值传递；需要修改对象、对象很大、或含同步原语时再考虑指针。
- 【建议】不要给 `map`、`slice`、`chan`、`interface` 传指针。
- 【建议】优先用切片而不是数组作为参数。

### 03A.2 返回值

- 【强制】不要返回多个仅用于流程控制的状态。
- 【强制】返回值个数不超过 3 个。

反例：

```go
isContinue, retCode := p.processUnity()
```

优秀示例：

```go
retCode := p.processUnity()
```

### 03A.3 实现

- 【建议】把 `map` 或 `slice` 存入结构体时，考虑浅拷贝和外部修改风险。

反例：

```go
func (d *Driver) SetTrips(trips []Trip) {
	d.trips = trips
}
```

优秀示例：

```go
func (d *Driver) SetTrips(trips []Trip) {
	d.trips = append([]Trip(nil), trips...)
}
```

### 03A.4 分组

- 【建议】构造函数靠前，导出函数在前，同一 receiver 靠拢，工具函数放末尾。

反例：

```go
func (s *something) Cost() {
	return calcCost(s.weights)
}

type something struct{ /* ... */ }

func calcCost(n []int) int { return 0 }

func (s *something) Stop() {}

func newSomething() *something { return &something{} }
```

优秀示例：

```go
type something struct{ /* ... */ }

func newSomething() *something { return &something{} }

func (s *something) Cost() {
	return calcCost(s.weights)
}

func (s *something) Stop() {}

func calcCost(n []int) int { return 0 }
```
