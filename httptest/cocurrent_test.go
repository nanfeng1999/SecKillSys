package httptest

import (
	"SecKillSys/api"
	"SecKillSys/model"
	"fmt"
	"github.com/gavv/httpexpect"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func getExpect(t *testing.T) *httpexpect.Expect{
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://127.0.0.1:20080",
		Reporter: httpexpect.NewAssertReporter(t),

		// use http.Client with a cookie jar and timeout
		Client: &http.Client{
			Jar:     httpexpect.NewJar(),
			Timeout: time.Second * 30,
		},
		// use verbose logging
		//Printers: []httpexpect.Printer{
		//	httpexpect.NewCurlPrinter(t),
		//	httpexpect.NewDebugPrinter(t, true),
		//},
	})

	return e
}

func TestRegister(t *testing.T) {
	e := getExpect(t)

	for i := 0 ; i <= 10000 ; i++{
		e.POST("/api/users").
			WithJSON(RegisterForm{"customer" + strconv.Itoa(i), "123456", model.NormalCustomer}).
			Expect()
	}

	for i := 0 ; i <= 10 ; i++{
		e.POST("/api/users").
			WithJSON(RegisterForm{"seller"+ strconv.Itoa(i), "123456", model.NormalSeller}).
			Expect()
	}
}

func TestAddCoupon(t *testing.T){
	e := getExpect(t)
	var demoAddCouponForm AddCouponForm = AddCouponForm{
		Name:        demoCouponName,
		Amount:      100 ,
		Stock:       demoStock,
		Description: "kiana: this is my good coupon",
	}


	for i := 0 ; i< 1 ; i++{
		resp := e.POST("/api/auth").
			WithJSON(LoginForm{"seller"+ strconv.Itoa(i), "123456"}).
			Expect()

		resp.Status(http.StatusOK).JSON().Object().
			ValueEqual(api.ErrMsgKey, "").
			ValueEqual("kind", model.NormalSeller)

		for j := 0 ; j < 10 ; j++{
			demoAddCouponForm.Name = "my_coupon" + strconv.Itoa(j)+ "seller"+strconv.Itoa(i)
			e.POST(addCouponPath, "seller"+strconv.Itoa(i)).
				WithJSON(demoAddCouponForm).
				WithHeader("Authorization", resp.Header("Authorization").Raw()).
				Expect().
				Status(http.StatusCreated).JSON().Object().
				ValueEqual(api.ErrMsgKey, "")
		}
	}
}

func TestFetch(t *testing.T){
	t1 := time.Now()
	e := getExpect(t)
	var val int32= 0
	wg := new(sync.WaitGroup)
	wg.Add(100)
	for i := 0 ; i< 500 ; i++{
		go func(j int) {
			defer wg.Done()
			resp := e.POST("/api/auth").
				WithJSON(LoginForm{"customer"+ strconv.Itoa(j), "123456"}).
				Expect().Status(http.StatusOK)
			cnt := e.PATCH(fetchCouponPath, "seller0", "my_coupon6seller0").
				WithHeader("Authorization", resp.Header("Authorization").Raw()).
				Expect()
			if cnt.Raw().StatusCode == 200{
				atomic.AddInt32(&val, 1)
			}
		}(i)

	}
	wg.Wait()

	fmt.Println("用户抢到的券总数为",val)
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
