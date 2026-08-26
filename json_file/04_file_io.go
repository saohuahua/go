// 04_file_io.go —— 文件 IO：读 / 写 / 目录 / 「文件不存在」
//
// 后端碰文件的场景：读配置、写日志、上传落盘
// 前端对照：Node 的 fs.readFile / fs.writeFile（浏览器没有 fs，这是纯后端新技能）
//
// 小文件（配置、小 JSON）用 os.ReadFile / os.WriteFile 一口搞定
// 大文件交给 05 的 Reader / Writer 流式处理
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// config 演示用配置结构（真实项目里这一坨来自 yaml/env，由 viper 解析）
type config struct {
	AppName string `json:"appName"`
	Port    int    `json:"port"`
	Debug   bool   `json:"debug"`
}

func demoFileIO() {
	fmt.Println("\n========== 04 文件 IO：写 / 读 / 存在性判断 ==========")

	// 演示目录放系统临时目录：不污染仓库，从哪个目录启动都不怕
	//? 为什么不用相对路径 "config.json"？相对路径 = CWD + 文件名
	//? CWD 取决于从哪里启动：根目录 go run ./json_file 和 cd json_file 后 go run . 不同
	//? 「配置/日志找不到」的经典元凶就是这个
	dir := filepath.Join(os.TempDir(), "go_json_file_demo")
	path := filepath.Join(dir, "config.json")

	// ---------- ① 写文件：MkdirAll + MarshalIndent + WriteFile ----------
	// MkdirAll 连父目录一起建，已存在也不报错（前端对照：fs.mkdir(dir, {recursive:true})）
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Println("① MkdirAll 失败：", err)
		return
	}

	// struct → 好看的 JSON → 落盘：这就是配置管理（viper 干的事）的迷你版
	cfg := config{AppName: "go-study", Port: 8080, Debug: false}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("① Marshal 失败：", err)
		return
	}

	// 0o644：Linux 权限位（属主读写、其他人只读）；Windows 忽略它，但习惯要写对
	// 前端对照：Node fs.writeFile(path, data) —— Go 多一个权限参数
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Println("① WriteFile 失败：", err)
		return
	}
	fmt.Printf("① 写入 %s（%d 字节）\n", path, len(data))

	// ---------- ② 读文件：ReadFile 一口读完（小文件专用） ----------
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("② ReadFile 失败：", err)
		return
	}
	var loaded config
	if err := json.Unmarshal(raw, &loaded); err != nil {
		fmt.Println("② Unmarshal 失败：", err)
		return
	}
	fmt.Printf("② 读回并解码：%+v\n", loaded)

	//? 缩进不影响解析 —— MarshalIndent 是给人看的，机器照样读

	// ---------- ③ 「文件不存在」是可预见的错误，不是异常 ----------
	_, err = os.ReadFile(filepath.Join(dir, "no_such.json"))
	fmt.Println("③ 读不存在的文件，err =", err)

	// 呼应 basics/error：fs.ErrNotExist 是标准库的哨兵错误
	// 用 errors.Is 精确判断，别用字符串匹配（错误文案可能变，哨兵不会）
	fmt.Println("③ errors.Is(err, fs.ErrNotExist) =", errors.Is(err, fs.ErrNotExist))

	// 前端对照：Node 里判断 err.code === 'ENOENT'，同一个意思

	// ---------- ④ os.Stat：不读内容，只看文件信息 ----------
	if info, err := os.Stat(path); err == nil {
		fmt.Printf("④ Stat：名字=%s 大小=%d字节 修改时间=%s\n",
			info.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
	}

	// ---------- ⑤ 清理现场 ----------
	if err := os.RemoveAll(dir); err != nil {
		fmt.Println("⑤ RemoveAll 失败：", err)
	}
	// 前端对照：fs.rm(dir, {recursive:true, force:true})

	fmt.Println("🔑 配置闭环：struct → Marshal → WriteFile → ReadFile → Unmarshal")
}
