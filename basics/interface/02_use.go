// 02_use.go —— 接口的三种使用场景（日常开发最常用的写法，重点）
//
// 接口不是用来炫技的，它有三个真实落点，全是后端天天写的：
//
//	① 函数参数：函数只认契约（面向接口编程）
//	② 结构体字段：依赖注入 —— 把"要用的东西"从外面塞进来，随时能换
//	③ 函数返回值：工厂函数返回接口，调用方只依赖能力
//
// 一条实战铁律：接收接口，返回结构体（Accept interfaces, return structs）
package main

import "fmt"

// 本节演示用：一个"通知"契约
type Notifier interface {
	Notify(userID int, msg string) error
}

// 实现 1：邮件通知
type EmailNotifier struct{}

func (EmailNotifier) Notify(userID int, msg string) error {
	fmt.Println("   [邮件] 发给用户", userID, "：", msg)
	return nil
}

// 实现 2：短信通知
type SMSNotifier struct{}

func (SMSNotifier) Notify(userID int, msg string) error {
	fmt.Println("   [短信] 发给用户", userID, "：", msg)
	return nil
}

type Wechater struct{}

func (Wechater) Notify(userID int, msg string) error {
	fmt.Println("   [微信] 发给用户", userID, "：", msg)
	return nil
}

// ---------- 用法①：接口作为函数参数（面向接口编程）----------
// 邮件能传、短信能传、以后加个微信通知 —— 这个函数永远不用改
func sendWelcome(n Notifier, userID int) {
	n.Notify(userID, "欢迎注册！")
}

// ---------- 用法②：接口作为结构体字段（依赖注入）----------
// OrderService 不自己 new 通知器，而是"被塞进来" —— 换实现不动业务代码
type OrderService struct {
	Notifier Notifier // 依赖注入：外部决定用邮件还是短信
	// Logger    Logger // 实际项目里还会注入日志、数据库等依赖
}

// PlaceOrder 下单成功后通知用户（业务代码只认 Notifier 契约）
func (s *OrderService) PlaceOrder(userID int) {
	fmt.Println("   [订单] 用户", userID, "下单成功")
	s.Notifier.Notify(userID, "您的订单已支付")
}

func demoUse() {
	fmt.Println("========== 02 接口的三种使用场景 ==========")

	// ---------- 用法①：函数参数 ----------
	fmt.Println("① 函数参数（面向接口编程）：")
	sendWelcome(EmailNotifier{}, 1)
	sendWelcome(SMSNotifier{}, 2)
	sendWelcome(Wechater{}, 3)

	// ---------- 用法②：结构体字段（依赖注入）----------
	fmt.Println("② 结构体字段（依赖注入）：")
	svc := &OrderService{Notifier: EmailNotifier{}} // 注入邮件
	svc.PlaceOrder(100)
	svc.Notifier = SMSNotifier{} // 随时换成短信，业务代码一行没动
	svc.PlaceOrder(101)

	svc.Notifier = Wechater{}
	svc.PlaceOrder(102)

	// ---------- 用法③：函数返回值（工厂）----------
	fmt.Println("③ 函数返回值（工厂函数）：")
	apiClient := newClient("mock") // 测试环境：假实现
	apiClient.Request()
	apiClient = newClient("real") // 生产环境：真实现
	apiClient.Request()
	// 调用方只依赖 APIClient 契约，具体是哪个实现由工厂决定

	// ---------- 补充：值接收者 vs 指针接收者实现接口 ----------
	//! 值接收者实现 → 值和指针都满足接口（Go 自动 &d 调方法）
	//! 指针接收者实现 → 只有指针满足，值不满足（编译直接报错）
	// 上面的 Notify 都是值接收者，所以 &EmailNotifier{} 也能传进来：
	var n Notifier = &EmailNotifier{}
	n.Notify(999, "指针也满足接口（值接收者）")

	fmt.Println("\n🔑 三落点：函数参数 / 字段注入 / 工厂返回；铁律：接收接口，返回结构体")
}

// ---------- 用法③ 的工厂 ----------

// APIClient 接口：客户端契约
type APIClient interface {
	Request()
}

// RealClient 生产实现：请求真实接口
type RealClient struct{}

func (RealClient) Request() { fmt.Println("   [real] 请求真实接口") }

// MockClient 测试实现：返回假数据
type MockClient struct{}

func (MockClient) Request() { fmt.Println("   [mock] 返回假数据，测试用") }

// newClient 工厂：根据配置返回不同实现（依赖注入的"选择器"）
// 返回类型写接口：调用方只依赖能力，具体实现由工厂决定
func newClient(mode string) APIClient {
	if mode == "mock" {
		return MockClient{}
	}
	return RealClient{}
}
