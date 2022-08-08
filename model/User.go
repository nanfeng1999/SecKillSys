package model

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	MinUserNameLen = 1
	MinPasswordLen = 1
	NormalCustomer = "customer"
	NormalSeller   = "seller"
)

type LoginUser struct {
	Username string	// 账号
	Password string // 密码
}

type RegisterUser struct {
	LoginUser
	Kind     string // 指示 是商家 还是 用户
}

// User 存入数据库中的结构体类型
type User struct {
	Id       int     `gorm:"primary_key;auto_increment"`
	Username string  `gorm:"type:varchar(20)"`
	Kind     string  `gorm:"type:varchar(20)"`
	Password string  `gorm:"type:varchar(32)"`
}

// IsCustomer 判断用户是否是消费者
func (user User)IsCustomer() bool {
	return user.Kind == NormalCustomer
}

// IsSeller 判断用户是否是售出者
func (user User)IsSeller() bool {
	return user.Kind == NormalSeller
}

// IsValidKind 判断用户类型是否合法
func IsValidKind(kind string) bool {
	return kind == NormalCustomer || kind == NormalSeller
}

// GetMD5 获取字符串的MD5结果
func GetMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}