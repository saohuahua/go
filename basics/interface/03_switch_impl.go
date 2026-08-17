// 03_switch_impl.go —— 换实现 + mock：接口在分层架构里的核心价值（重点）
//
// 后端分层架构（handler → service → repository）里，
// repository（数据访问层）几乎总是定义成接口：
//
//	生产环境：MySQL 实现（真连数据库）
//	测试环境：内存 mock 实现（不连库、秒回假数据）
//
// 业务层只认 repository 接口 —— 换数据库、写测试，业务代码一行不用动。
// 这就是接口在真实项目里最大的用处，也是面试官最爱问的"为什么要分层"。
package main

import "fmt"

// ---------- 1. 定义契约：用户数据仓库 ----------

// User 数据模型
type User struct {
	ID   int
	Name string
}

// UserRepository 数据访问契约（实际项目里放在 repository 层）
type UserRepository interface {
	GetByID(id int) (*User, error)
	Save(u *User) error
}

// ---------- 2. 生产实现：MySQL ----------
// 真实项目里这里写 SQL / GORM 查询，这里用 map 假装数据库，重点看接口用法
type MySQLUserRepo struct {
	users map[int]*User
}

func (r *MySQLUserRepo) GetByID(id int) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("用户 %d 不存在", id)
	}
	return u, nil
}

func (r *MySQLUserRepo) Save(u *User) error {
	fmt.Println("   [MySQL] 已写入数据库：", u.Name)
	r.users[u.ID] = u
	return nil
}

// ---------- 3. 测试实现：内存 mock ----------
// 单测用：不连数据库、秒回假数据；还能记录"被调了几次"（测试断言用）
type MockUserRepo struct {
	calls int // GetByID 被调用的次数（测试里断言"确实调用了"）
}

func (r *MockUserRepo) GetByID(id int) (*User, error) {
	r.calls++
	if id == 99 {
		return &User{ID: id, Name: "测试用户"}, nil
	}
	return nil, fmt.Errorf("mock：用户 %d 不存在", id)
}

func (r *MockUserRepo) Save(u *User) error {
	fmt.Println("   [mock] 假装保存成功：", u.Name)
	return nil
}

// ---------- 业务层：只认 UserRepository 契约 ----------

// getUserProfile 业务逻辑，完全不知道下面是 MySQL 还是 mock
func getUserProfile(repo UserRepository, id int) string {
	u, err := repo.GetByID(id)
	if err != nil {
		return "（查不到：" + err.Error() + "）"
	}
	return u.Name
}

func demoSwitchImpl() {
	fmt.Println("========== 03 换实现 + mock：分层架构的接口 ==========")

	// 生产：MySQL 实现
	fmt.Println("① 生产环境（MySQL 实现）：")
	mysqlRepo := &MySQLUserRepo{users: map[int]*User{}}
	mysqlRepo.Save(&User{ID: 1, Name: "小明"})
	fmt.Println("   getUserProfile →", getUserProfile(mysqlRepo, 1))
	fmt.Println("   getUserProfile →", getUserProfile(mysqlRepo, 999)) // 查不到的错误路径

	// 测试：mock 实现（不连库）
	fmt.Println("② 测试环境（mock 实现）：")
	mockRepo := &MockUserRepo{}
	fmt.Println("   getUserProfile →", getUserProfile(mockRepo, 99))
	fmt.Println("   getUserProfile →", getUserProfile(mockRepo, 1))
	fmt.Println("   mock 的 GetByID 被调了", mockRepo.calls, "次（测试断言：确认调用了）")

	// ---------- 关键结论 ----------
	// 同一个 getUserProfile，两种实现都能跑 —— 业务代码从头到尾不知道自己连的是啥
	// 以后换 Postgres、加 Redis 缓存 —— 只要实现 UserRepository，业务层一行不用改
	fmt.Println("\n🔑 repository 接口 = 分层架构的地基：换实现 / 写测试全靠它")
}
