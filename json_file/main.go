// main.go —— json + 文件 IO：数据的进出（进 Gin 前的补课章）
//
// 本文件夹共 5 个演示文件 + 本入口文件
//
// 01_json_basics.go     —— Marshal / Unmarshal / struct tag / HTML 转义
// 02_json_zero_value.go —— 零值困境：omitempty / 指针 / any / float64 精度
// 03_time_format.go     —— 时间布局 / Parse / 时区 / JSON 里的时间
// 04_file_io.go         —— 读写文件 / 建目录 / 判断文件不存在
// 05_reader_writer.go   —— io.Reader/Writer / io.Copy / Builder / Scanner
//
// 运行方法
//
// cd json_file && go run .
// go run ./json_file    （在项目根目录执行）
//
// 学习方法
//
// 1 整篇跑一遍，输出和注释对着看
// 2 把 main 里只留一个 demoXxx()，自己改数据做实验
// 3 每节都配 JS/TS 对照，把前端经验迁移过来
package main

import "fmt"

func main() {
	// 依次演示，建议逐个放开调用，边看输出边理解
	demoJSONBasics()    // 01 Marshal / Unmarshal / tag
	demoJSONZeroValue() // 02 零值困境
	demoTimeFormat()    // 03 time 布局 / 时区
	demoFileIO()        // 04 文件读写
	demoReaderWriter()  // 05 流式 IO

	fmt.Println("\n🎉 json_file 完成：数据的进出都通了，参数绑定和配置管理有地基了")
}
