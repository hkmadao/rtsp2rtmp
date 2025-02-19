package entity

import (
	"time"
)
// 令牌
type Token struct {
	// 令牌主属性
	IdToken string  `orm:"pk;column(id_sys_token)" json:"idToken"`
	// 用户名称:
	Username string `orm:"column(username)" json:"username"`
	// 昵称:
	NickName string `orm:"column(nick_name)" json:"nickName"`
	// 创建时间:
	CreateTime time.Time `orm:"column(create_time)" json:"createTime"`
	// 令牌:
	Token string `orm:"column(token)" json:"token"`
	// 过期时间:
	ExpiredTime time.Time `orm:"column(expired_time)" json:"expiredTime"`
	// 用户信息序列化:
	UserInfoString string `orm:"column(user_info_string)" json:"userInfoString"`
}
