# Go 代码审查清单 02: 错误处理与资源管理

本文件聚焦错误传播、日志、panic/recover、资源关闭和内存泄漏。

## 02. 错误处理

### 02.1 返回错误

- 【强制】错误统一作为最后一个返回值返回。
- 【强制】遵循 fail fast。
- 【建议】接收 `context.Context` 的函数通常也应返回 `error`。

### 02.2 日志

- 【强制】不要只打日志不返回错误，除非当前层就是最终消费错误的边界。
- 【建议】避免同一错误多层重复记录。

### 02.3 处理方式

- 【强制】不要吞掉错误。
- 【强制】出现失败后不要盲目继续执行。

### 02.4 错误字符串

- 【强制】错误字符串简洁稳定，避免首字母大写和末尾标点。

反例：

```go
errors.New("File Not Found.")
```

优秀示例：

```go
errors.New("file not found")
```

### 02.5 类型断言

- 【强制】类型断言优先使用 `v, ok := ...`。

反例：

```go
s := i.(string)
```

优秀示例：

```go
s, ok := i.(string)
if !ok {
	return errors.New("unexpected type")
}
```

### 02.6 panic / recover

- 【强制】`panic` 必须在当前 goroutine 内捕获。
- 【强制】自行启动的 goroutine 要考虑 `recover`。

优秀示例：

```go
go func() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
	}()

	doWork()
}()
```

## 02A. 资源管理

### 02A.1 主动关闭资源

- 【强制】文件、连接、响应体、锁等资源必须释放。
- 【强制】`defer` 关闭资源时，要判断关闭错误是否重要。
- 【建议】避免在长循环里无控制地 `defer`。

反例：

```go
for _, name := range files {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
}
```

优秀示例：

```go
for _, name := range files {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	if err := handleFile(f); err != nil {
		_ = f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
}
```

### 02A.2 内存泄漏

- 【建议】从大切片或长字符串上截取小片段时，关注底层内存未释放问题。
- 【建议】大流量场景下，map 删除元素后不会自动收缩，必要时重建。

反例：

```go
func head(data []byte) []byte {
	return data[:16]
}
```

优秀示例：

```go
func head(data []byte) []byte {
	out := make([]byte, 16)
	copy(out, data[:16])
	return out
}
```
