// 03_time_format.go —— time 包：布局 / Parse / 时区 / JSON 里的时间
//
// 后端 API 十个有九个要返回 createdAt / updatedAt，先记住三个事实
//
// 1 布局不写 YYYY-MM-DD，写 Go 的参考时间 "2006-01-02 15:04:05"（Go 的生日）
// 2 time.Time 自带时区（JS Date 只是时间戳，显示时才套本地时区）
// 3 time.Time 转 JSON 自动变成 RFC3339 字符串，前端 new Date() 直接能吃
//
// 前端对照：dayjs 的 format("YYYY-MM-DD HH:mm:ss") ↔ time.Format("2006-01-02 15:04:05")
package main

import (
	"encoding/json"
	"fmt"
	"time"

	// 空导入：不用这个包的任何函数，只为执行它内部的 init()（把时区数据库嵌进程序）
	// 没有它，Windows 上 LoadLocation("Asia/Shanghai") 会报错（找不到系统的时区数据）
	// 以后会见的 _ "github.com/.../docs"（Swagger）也是同一个套路
	_ "time/tzdata"
)

func demoTimeFormat() {
	fmt.Println("\n========== 03 time：布局 / Parse / 时区 / JSON ==========")

	// ---------- ① 布局是「参考时间」，不是格式占位符 ----------
	now := time.Now()
	fmt.Println("① time.Now()：", now)
	fmt.Println("① 完整布局：", now.Format("2006-01-02 15:04:05"))
	fmt.Println("① 只要日期：", now.Format("2006-01-02"))
	fmt.Println("① 只要时间：", now.Format("15:04:05"))
	fmt.Println("① RFC3339 常量：", now.Format(time.RFC3339)) // API 里最常见

	//? 为什么是 2006-01-02 15:04:05？按美式日期念出来：
	//? 1月2日 15点4分5秒 2006年 → 01 02 15 04 05 2006，数字本身就是占位符
	//? 前端对照：dayjs 把占位符写成 YYYY/MM/DD；Go 干脆拿一个真实日期当占位符

	// ---------- ② 写错布局的下场（经典翻车现场） ----------
	// 想输出 "2026-08-27" 这种日期，顺手把布局写成 "2026-08-27"？看输出变成什么：
	fmt.Println("② 把布局写成日期本身 →", now.Format("2026-08-27"))

	//! 布局里的数字会被当占位符解析："2026" 里的 2 是「日」、02 也是「日」……
	//! 布局永远写参考时间 2006-01-02，想要什么格式就挑对应占位符来拼

	// ---------- ③ Parse：字符串 → time.Time ----------
	// 必须告诉它字符串长什么样（布局）；解析失败返回 error，而不是 NaN
	t, err := time.Parse("2006-01-02 15:04:05", "2026-08-27 20:30:00")
	if err != nil {
		fmt.Println("③ Parse 失败：", err)
		return
	}
	fmt.Println("③ Parse 出来：", t)

	//? Parse 默认把「无时区字符串」当 UTC 解析（上面输出带 +0000）
	//? 要按本地时区解析：time.ParseInLocation(布局, 字符串, loc)

	// ---------- ④ 时区：同一时刻，不同表述 ----------
	// time.Time = 时刻 + 时区；In(loc) 转换时区：时刻不变、表述变
	sh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("④ LoadLocation 失败：", err)
		return
	}
	ny, _ := time.LoadLocation("America/New_York")
	fmt.Println("④ 上海视角：", t.In(sh)) // 已经是 28 号凌晨
	fmt.Println("④ 纽约视角：", t.In(ny)) // 还是 27 号下午

	// 前端对照：JS Date 没有时区属性，只有「时间戳 + 显示时套本地格式化」两步
	// Go 的 Time 把时区绑在值上，跨时区转换是显式的一步，更不容易出错

	// ---------- ⑤ Duration：时间段（不是时间点） ----------
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("⑤ time.Since 耗时 = %v\n", time.Since(start).Round(time.Millisecond))
	d := 2*time.Hour + 30*time.Minute // Duration 字面量：底层是 int64 纳秒
	fmt.Printf("⑤ 2h30m = %v = %.0f 分钟\n", d, d.Minutes())

	// ---------- ⑥ time.Time 的 JSON 形态 ----------
	// 自动 Marshal 成 RFC3339 字符串（带时区偏移），前端 new Date(str) 直接解析
	b, _ := json.Marshal(map[string]time.Time{"createdAt": now})
	fmt.Printf("⑥ JSON：%s\n", b)

	//? 想自定义输出格式（比如 "2026-08-27 20:30:00"）？
	//? 给响应 DTO 定义 string 字段，手动 Format 后再返回，别在全局改序列化行为
	//? GORM 的时间列也是 time.Time，时区由连接串的 loc 参数决定（Week4 见）

	fmt.Println("🔑 布局背「2006-01-02 15:04:05」；时间自带时区；JSON 自动 RFC3339")
}
