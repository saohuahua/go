// 05_reader_writer.go —— io.Reader / io.Writer：Go 的「流」宇宙
//
// basics/interface/04 里见过 io.Writer（CountingWriter），这里接上真实场景
//
// Go 数据流动的核心抽象：不管数据从哪来到哪去，两头都是接口
//
// io.Reader：能被读 —— 文件、网络 body、字符串、缓冲区……
// io.Writer：能被写 —— 文件、响应、缓冲区、Builder……
//
// 前端对照：Web Stream 的 ReadableStream / WritableStream
// 不同点：Go 的最小契约只有一个方法，小到极致，所以万物都能实现
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func demoReaderWriter() {
	fmt.Println("\n========== 05 io.Reader / io.Writer：流式读写 ==========")

	// ---------- ① 万物皆 Reader + io.Copy 万能接头 ----------
	// strings.NewReader 把字符串包成 Reader（net_http 里给请求造 body 用过）
	// os.Stdout 也是个 Writer —— 标准输出
	// io.Copy：数据从 Reader 流进 Writer，一根管子接两头
	src := strings.NewReader("go\n后端\n流式处理\n")
	n, _ := io.Copy(os.Stdout, src) // demo 简化；真实代码要查 err
	fmt.Printf("① io.Copy 往 os.Stdout 搬了 %d 字节\n", n)

	//? 前端对照：ReadableStream pipe 到 WritableStream
	//? 文件上传的本质：io.Copy(目标文件, 上传流)；下载：io.Copy(响应, 文件)
	//? net_http 的 writeJSON 里 NewEncoder(w) 直接吃 ResponseWriter，就是这套

	// ---------- ② strings.Builder：拼字符串的正确姿势 ----------
	// 字符串不可变，+= 每拼一次分配一次新内存；Builder 内部维护 []byte 增长
	// 前端对照：老 JS 时代「数组 push + join」防字符串 += 的坑，同一个思路
	var b strings.Builder
	for i := 1; i <= 3; i++ {
		fmt.Fprintf(&b, "第%d行\n", i) // Builder 实现了 io.Writer，fmt 直接写
	}
	fmt.Printf("② Builder 拼出：%q\n", b.String())

	// ---------- ③ bufio.Scanner：一行一行流式读 ----------
	// ReadFile 是一口吞，文件 2GB 就吃 2GB 内存；Scanner 是吸管，一次含一行
	//? 为什么新建 src2？① 里 src 已被 Copy 读完，Reader 读到底就没了
	src2 := strings.NewReader("第一行\n第二行\n第三行没有换行符结尾")
	scanner := bufio.NewScanner(src2)
	for scanner.Scan() { // Scan 每读到一行返回 true（换行符被剥掉）
		fmt.Printf("③ 扫到一行：%q\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil { // 循环结束必须查一次 Err
		fmt.Println("③ Scanner 出错：", err)
	}

	//! Scanner 默认单行上限 64KB，超长行报 bufio.Scanner: token too long
	//! 大文件/长行要先扩：scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	// ---------- ④ 为什么 json.NewEncoder(w) 直接吃 Writer ----------
	// Marshal：先在内存里生成完整 JSON，再整体交出去
	// NewEncoder(w).Encode：边序列化边写进 w —— 大响应省一倍内存
	// Gin 的 c.JSON 内部走的就是 Encoder 思路
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(map[string]int{"a": 1, "b": 2})
	fmt.Printf("④ Encoder 直写 Writer：%q\n", buf.String())

	fmt.Println("🔑 数据流动 = Reader 进 Writer 出，io.Copy 是万能接头")
}
