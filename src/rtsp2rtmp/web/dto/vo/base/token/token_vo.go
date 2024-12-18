package vo

import (
	"time"
)
// 令牌
type TokenVO struct {
	// 令牌主属性
	IdToken string  `json:"idToken"`
	// 用户名称:
	Username string `json:"username"`
	// 昵称:
	NickName string `json:"nickName"`
	// 创建时间:
	CreateTime time.Time `json:"createTime"`
	// 令牌:
	Token string `json:"token"`
	// 过期时间:
	ExpiredTime time.Time `json:"expiredTime"`
	// 用户信息序列化:
	UserInfoString string `json:"userInfoString"`
}