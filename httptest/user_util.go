package httptest

import (
	"SecKillSys/api"
	"SecKillSys/model"
	"github.com/gavv/httpexpect"
	"net/http"
)

// 本文件存放了一些demo用户信息，demo用户的注册/登录函数
// 定义了用户登出函数
// 还定义了注册/登录用户的表格

const demoSellerName = "kiana00"
const demoCustomerName = "jinzili01"
const demoArCustomerName = "karsa01" // name of another customer
const demoPassword = "shen6508"

type RegisterForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Kind     string `form:"kind"`
}

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func registerDemoUsers(e *httpexpect.Expect) {
	// 注册商家的账户
	e.POST("/api/users").
		WithJSON(RegisterForm{demoSellerName, demoPassword, model.NormalSeller}).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

	// 注册消费者的账户
	e.POST("/api/users").
		WithJSON(RegisterForm{demoCustomerName, demoPassword, model.NormalCustomer}).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

	// 注册另一个消费者的账户
	e.POST("/api/users").
		WithJSON(RegisterForm{demoArCustomerName, demoPassword, model.NormalCustomer}).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

}

func demoCustomerLogin(e *httpexpect.Expect) *httpexpect.Response {
	res := e.POST("/api/auth").
		WithJSON(LoginForm{demoCustomerName, demoPassword}).
		Expect()
	res.Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "").
		ValueEqual("kind", model.NormalCustomer)

	return res
}

func demoArCustomerLogin(e *httpexpect.Expect) {
	e.POST("/api/auth").
		WithJSON(LoginForm{demoArCustomerName, demoPassword}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "").
		ValueEqual("kind", model.NormalCustomer)
}

func demoSellerLogin(e *httpexpect.Expect) *httpexpect.Response {
	res := e.POST("/api/auth").
		WithJSON(LoginForm{demoSellerName, demoPassword}).
		Expect()
	res.Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "").
		ValueEqual("kind", model.NormalSeller)
	return res
}

func logout(e *httpexpect.Expect) {
	e.POST("/api/auth/logout").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "log out.")
}
