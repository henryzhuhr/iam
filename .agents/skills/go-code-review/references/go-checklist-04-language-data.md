# Go 代码审查清单 04: 语言细节与数据结构

本文件聚焦变量声明、控制流、结构体、slice、map、序列化和 time 等语言细节。

## 04. 语言细节与数据结构

### 04.1 变量声明

- 【建议】变量就近声明并缩小作用域。
- 【建议】局部变量优先用 `:=`，零值局部变量优先用 `var`。
- 【建议】相似声明使用分组。

反例：

```go
const a = 1
const b = 2

var x = 1
var y = 2
```

优秀示例：

```go
const (
	a = 1
	b = 2
)

var (
	x = 1
	y = 2
)
```

反例：

```go
func foo() {
	var s = "foo"
}
```

优秀示例：

```go
func foo() {
	s := "foo"
}
```

反例：

```go
err := file.Chmod(0664)
if err != nil {
	return err
}
```

优秀示例：

```go
if err := file.Chmod(0664); err != nil {
	return err
}
```

### 04.2 语句

- 【强制】禁止变量遮蔽。
- 【强制】指针解引用前确认非 nil。
- 【强制】外部输入参与数值运算时，注意溢出、截断和符号错误。

### 04.3 魔数

- 【强制】除 `0` 和 `1` 外，避免直接使用魔数。
- 【强制】重复出现的字符串应提取常量。

### 04.4 for / range / switch / goto

- 【建议】不要持有循环变量地址。
- 【强制】只需要 `range` 的 key 时丢弃 value，只需要 value 时把 key 写成 `_`。
- 【强制】遍历大结构体切片优先按下标访问。
- 【强制】`switch` 必须有 `default`。
- 【强制】业务代码禁止使用 `goto`。

反例：

```go
for _, v := range ints {
	fmt.Println(&v)
}
```

优秀示例：

```go
for _, v := range ints {
	v := v
	fmt.Println(&v)
}
```

优秀示例：

```go
for key := range m {
	if key.expired() {
		delete(m, key)
	}
}
```

优秀示例：

```go
sum := 0
for _, v := range array {
	sum += v
}
```

反例：

```go
for _, item := range largeStructs {
	item.Age = 10
}
```

优秀示例：

```go
for i := range largeStructs {
	largeStructs[i].Age = 10
}
```

### 04.5 结构体

- 【强制】字段排序考虑内存对齐。
- 【强制】结构体拷贝要考虑浅拷贝问题。
- 【强制】结构体初始化使用带字段名的字面量。
- 【建议】结构体指针初始化优先 `&T{}`。

反例：

```go
user := User{"alice", 18}
```

优秀示例：

```go
user := User{
	Name: "alice",
	Age:  18,
}
```

### 04.6 slice

- 【强制】`slice` 应合理初始化，必要时预估容量。
- 【强制】索引前确保不越界。
- 【强制】`copy` 前确认目标长度足够。
- 【建议】空切片优先 `var s []T`。
- 【建议】零长度返回值优先 `nil`，空判断统一 `len(s) == 0`。

反例：

```go
nums := []int{}
```

优秀示例：

```go
nums := make([]int, 0, n)
```

反例：

```go
src := []int{0, 1, 2}
var dst []int
copy(dst, src)
```

优秀示例：

```go
src := []int{0, 1, 2}
dst := make([]int, len(src))
copy(dst, src)
```

反例：

```go
if x == "" {
	return []int{}
}

func isEmpty(s []string) bool {
	return s == nil
}
```

优秀示例：

```go
if x == "" {
	return nil
}

func isEmpty(s []string) bool {
	return len(s) == 0
}
```

### 04.7 map

- 【强制】固定元素列表直接用字面量，其他场景用 `make`。
- 【强制】可预估容量时显式给容量。
- 【强制】判断 key 是否存在要用 `v, ok := m[k]`。
- 【强制】禁止向 nil map 写入。

反例：

```go
m := make(map[T1]T2, 3)
m[k1] = v1
m[k2] = v2
m[k3] = v3
```

优秀示例：

```go
m := map[T1]T2{
	k1: v1,
	k2: v2,
	k3: v3,
}
```

### 04.8 序列化

- 【强制】序列化字段、tag、默认值语义必须准确，避免“看起来生效、实际无效”。

优秀示例：

```go
type Response struct {
	Message string `json:"message"`
}
```

### 04.9 time

- 【强制】时间运算、时区和超时处理要明确。

优秀示例：

```go
ctx, cancel := context.WithTimeout(parent, 100*time.Millisecond)
defer cancel()
```
