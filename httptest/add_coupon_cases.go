package httptest

import (
	"SecKillSys/api"
	"github.com/gavv/httpexpect"
	"net/http"
)

/*
该文件下依赖于注册过的demo用户，需要先调用registerDemoUsers
该文件定义了添加优惠券的各种函数
*/

/* 定义了添加优惠券的表格，函数等等 */
const addCouponPath = "/api/users/{username}/coupons"

type AddCouponForm struct {
	Name        string `json:"name"`        // 谁加这个优惠券
	Amount      int    `json:"amount"`      // 添加的数量
	Description string `json:"description"` // 描述
	Stock       int    `json:"stock"`       // 库存
}

// 定义了demo优惠券
var demoCouponName = "my_coupon"
var demoAmount = 10
var demoStock = 50
var demoAddCouponForm AddCouponForm = AddCouponForm{
	Name:        demoCouponName,                  //
	Amount:      demoAmount,                      //
	Stock:       demoStock,                       //
	Description: "kiana: this is my good coupon", //
}

// 测试添加优惠券时的表格格式 测试格式错误时候的情况 返回结果是否符合预期
func testAddCouponWrongFormat(e *httpexpect.Expect) {
	// 退出
	logout(e)
	// 商家登录
	demoSellerLogin(e)

	// amount值为0的情况下
	amountNotNumberForm := demoAddCouponForm
	amountNotNumberForm.Amount = 0
	e.POST(addCouponPath, demoSellerName).
		WithForm(amountNotNumberForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Amount field wrong format.")

	// stock值为0的情况下
	stockNotNumberForm := demoAddCouponForm
	stockNotNumberForm.Stock = 0
	e.POST(addCouponPath, demoSellerName).
		WithForm(stockNotNumberForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Stock field wrong format.")

	// 优惠券名为空的情况下
	emptyCouponNameForm := demoAddCouponForm
	emptyCouponNameForm.Name = ""
	e.POST(addCouponPath, demoSellerName).
		WithForm(emptyCouponNameForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Coupon name or description should not be empty.")

	// 优惠券描述为空的情况下
	emptyDescriptionForm := demoAddCouponForm
	emptyDescriptionForm.Description = ""
	e.POST(addCouponPath, demoSellerName).
		WithForm(emptyDescriptionForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Coupon name or description should not be empty.")
}

// 测试非商家添加优惠券或为其它用户添加优惠券
func testAddCouponWrongUser(e *httpexpect.Expect) {
	// 登录顾客
	demoCustomerLogin(e)

	// 顾客不可添加优惠券
	e.POST(addCouponPath, demoCustomerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Only sellers can create coupons.")

	// 登录商家
	demoSellerLogin(e)
	// 不可为其它用户添加优惠券
	e.POST(addCouponPath, demoCustomerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Cannot create coupons for other users.")
}

// 测试未登录添加优惠券
func testAddCouponNotLogIn(e *httpexpect.Expect) {
	logout(e)

	e.POST(addCouponPath, demoSellerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Not Logged in.")
}

func testAddCouponDuplicate(e *httpexpect.Expect) {
	resp := demoSellerLogin(e)

	e.POST(addCouponPath, demoSellerName).
		WithForm(demoAddCouponForm).
		WithHeader("Authorization", resp.Header("Authorization").Raw()).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

	// 添加重复优惠券失败
	e.POST(addCouponPath, demoSellerName).
		WithJSON(demoAddCouponForm).
		WithHeader("Authorization", resp.Header("Authorization").Raw()).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Create failed. Maybe (username,coupon name) duplicates")
}

// 添加demo优惠券(事先不得添加过)
func demoAddCoupon(e *httpexpect.Expect) {
	resp := demoSellerLogin(e)

	e.POST(addCouponPath, demoSellerName).
		WithJSON(demoAddCouponForm).
		WithHeader("Authorization", resp.Header("Authorization").Raw()).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}
