package po

// 系统用户
type UserPO struct {
	// 系统用户id
	IdUser string `json:"idUser"`
	// 登录账号 :
	Account string `json:"account"`
	// 用户密码 :
	UserPwd string `json:"userPwd"`
	// 手机号码:
	Phone string `json:"phone"`
	// 邮箱:
	Email string `json:"email"`
	// 姓名 :
	Name string `json:"name"`
	// 昵称:
	NickName string `json:"nickName"`
	// 性别:
	Gender string `json:"gender"`
	// 启用标志:
	FgActive bool `json:"fgActive"`
}
