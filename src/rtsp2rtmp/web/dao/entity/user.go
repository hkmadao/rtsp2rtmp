package entity

// 系统用户
type User struct {
	// 系统用户id
	IdUser string `orm:"pk;column(id_user)" json:"idUser"`
	// 登录账号 :
	Account string `orm:"column(account)" json:"account"`
	// 用户密码 :
	UserPwd string `orm:"column(user_pwd)" json:"userPwd"`
	// 手机号码:
	Phone string `orm:"column(phone)" json:"phone"`
	// 邮箱:
	Email string `orm:"column(email)" json:"email"`
	// 姓名 :
	Name string `orm:"column(name)" json:"name"`
	// 昵称:
	NickName string `orm:"column(nick_name)" json:"nickName"`
	// 性别:
	Gender string `orm:"column(gender)" json:"gender"`
	// 启用标志:
	FgActive bool `orm:"column(fg_active)" json:"fgActive"`
}
