package httptest

import (
	"SecKillSys/api"
	"github.com/gavv/httpexpect"
	"net/http"
)

var fetchCouponPath = "/api/users/{username}/coupons/{name}"

func fetchDemoCouponSuccess(e *httpexpect.Expect, resp *httpexpect.Response) {
	e.PATCH(fetchCouponPath, demoSellerName, demoCouponName).
		WithHeader("Authorization", resp.Header("Authorization").Raw()).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}

func fetchDemoCouponFail(e *httpexpect.Expect,resp *httpexpect.Response ) {
	e.PATCH(fetchCouponPath, demoSellerName, demoCouponName).
		WithHeader("Authorization", resp.Header("Authorization").Raw()).
		Expect().
		Status(http.StatusNoContent).
		Body().Empty()
}