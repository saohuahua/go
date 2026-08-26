# json + 文件 IO：数据的进出

> 面向「前端（TS/JS）转 Go 后端」的学习代码
> Gin 的参数绑定（ShouldBindJSON）、配置管理（viper）、日志落盘、文件上传，地基全在这一章

## 怎么运行

```bash
cd json_file
go run .
```

> 04 会往系统临时目录写一个 `go_json_file_demo` 文件夹并在结束时删掉，不污染仓库

## 注释标记

| 标记 | 含义 |
| --- | --- |
| `//!` | 必踩的坑 / 必记重点 |
| `//?` | 思考题 |
| `//*` | 关键铺垫 / 前端对照 |

> Better Comments 标记只出现在函数体内

## 学习顺序

| 文件 | 内容 | Gin / 项目对照 |
| --- | --- | --- |
| [01_json_basics.go](01_json_basics.go) | Marshal/Unmarshal、json tag、HTML 转义 | `c.JSON` / `ShouldBindJSON` 的底层 |
| [02_json_zero_value.go](02_json_zero_value.go) | omitempty、指针、any、float64 精度 | DTO 设计、PATCH 局部更新 |
| [03_time_format.go](03_time_format.go) | 时间布局、Parse、时区、JSON 时间 | createdAt 字段、GORM 时间列 |
| [04_file_io.go](04_file_io.go) | 读写文件、建目录、判断不存在 | 配置读取、日志、上传落盘 |
| [05_reader_writer.go](05_reader_writer.go) | io.Reader/Writer、io.Copy、Builder、Scanner | 流式响应、文件上传下载 |

## 一句话总结

| 概念 | 一句话 | 前端对照 |
| --- | --- | --- |
| json tag | 一份 tag 同时管序列化和反序列化的字段名 | TS DTO 字段映射 |
| `json:"-"` | 序列化时忽略（密码等敏感字段） | DTO 里删字段 |
| omitempty | 零值字段不出现在 JSON 里 | TS 可选字段 `age?: number` |
| 指针字段 | nil→null、非nil→有值，区分「没填」和「填了0」 | `score?: number \| null` |
| map[string]any | 结构未知的 JSON 兜底 | `Record<string, any>` |
| float64 精度 | any 解码数字全是 float64，大整数丢精度 | JS Number 同款 IEEE 754 坑 |
| 时间布局 | 布局写参考时间 `2006-01-02 15:04:05` | dayjs 的 YYYY-MM-DD |
| time.Time→JSON | 自动 RFC3339 字符串（带时区） | `new Date()` 直接解析 |
| os.ReadFile | 一口读完整个文件（小文件） | Node `fs.readFile` |
| errors.Is(err, fs.ErrNotExist) | 判断「文件不存在」的标准姿势 | `err.code === 'ENOENT'` |
| io.Reader/Writer | 数据进出统一契约 | Readable/WritableStream |
| io.Copy | Reader→Writer 万能接头 | stream 的 pipe |
| bufio.Scanner | 一行一行流式读，大文件不爆内存 | Node readline |
| strings.Builder | 高效拼字符串 | 数组 push + join |

## 面试速背

- **json tag**：`json:"name"` 定字段名；`json:"-"` 忽略；`json:"name,omitempty"` 零值省略
- **Unmarshal 细节**：第二参数必须传指针；字段匹配大小写不敏感；多余字段默认忽略
- **零值困境**：0/""/false 表达不了「没填」→ 指针字段（nil→null）或 omitempty
- **any 的数字**：解码一律 float64 → 大整数丢精度（雪花ID）；解法：定义 struct 或 `Decoder.UseNumber`
- **HTML 转义**：Marshal 默认转义 `<` `>` `&`（防 XSS）；关闭要 `Encoder.SetEscapeHTML(false)`
- **时间布局**：参考时间 `2006-01-02 15:04:05`（1月2日15点4分5秒2006年），写错输出乱码
- **Parse**：默认按 UTC 解析无时区字符串；本地时区用 `time.ParseInLocation`
- **time.Time JSON**：RFC3339 字符串带时区偏移，前端 `new Date()` 直接可用
- **空导入**：`import _ "pkg"` 不用其函数，只为执行包的 init()（tzdata、Swagger docs 都靠它）
- **文件不存在**：`errors.Is(err, fs.ErrNotExist)`，不靠字符串匹配
- **io.Copy**：`func Copy(dst Writer, src Reader) (int64, error)`，上传下载的本质
- **Scanner 上限**：单行默认 64KB，超长行要 `scanner.Buffer` 扩，否则 token too long
- **Encoder vs Marshal**：NewEncoder(w) 边序列化边写，大对象省一次完整内存拷贝

## 最容易踩的坑

1. **把布局写成目标日期**：`Format("2026-08-27")` 不会输出当天日期，数字被当占位符解析成乱码
2. **用 any 接大数字 ID**：解码变 float64 丢精度；已知结构永远定义 struct，大 ID 用 string
3. **零值当「没填」**：更新接口 age=0 和「没传」分不清，要区分就用指针字段
4. **相对路径**：CWD 取决于启动位置，`go run ./json_file` 和 `cd json_file && go run .` 的相对路径指向不同地方
5. **大文件 ReadFile**：一口吞进内存；日志/大文件用 Scanner 或 io.Copy 流式处理

## 下一步

1. （可选）补 `struct_embed/`：结构体嵌入与组合 —— GORM 的 `gorm.Model` 嵌入、tag 系统讲
2. 进 Gin：把本章 DTO 写法直接用到 `ShouldBindJSON` / `c.JSON`
3. 项目期回看：04 的配置闭环 → viper；05 的 `io.Copy` → 文件上传/下载

> 建议顺序：把 02 的「指针字段」在 Gin 的 PATCH 接口里真用一次，才算真懂
