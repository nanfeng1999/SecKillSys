package engine

import (
	"SecKillSys/api"
	"SecKillSys/data"
	"SecKillSys/middleware/jwt"
	"SecKillSys/model"
	"encoding/gob"
	"github.com/gin-gonic/gin"
)

func SeckillEngine() *gin.Engine {
	// 开启gin
	router := gin.New()

	// 注册gob编码
	gob.Register(&model.User{})

	// 创建用户相关的路由
	userRouter := router.Group("/api/users")
	userRouter.POST("", api.RegisterUser) //注册

	// 一下请求都需要通过jwt做用户授权 添加中间件
	userRouter.Use(jwt.JWTAuth())//这些请求都需要通过jwt做用户授权 添加中间件
	{
		// 抢优惠券
		userRouter.PATCH("/:username/coupons/:name", api.FetchCoupon)
		// 获取优惠券的信息
		userRouter.GET("/:username/coupons", api.GetCoupons)
		// 添加优惠券
		userRouter.POST("/:username/coupons", api.AddCoupon)
	}

	// 创建验证相关的路由
	authRouter := router.Group("/api/auth") //登录和注销
	{
		// 登录验证
		authRouter.POST("", api.LoginAuth)
		// 退出验证
		authRouter.POST("/logout", api.Logout)
	}

	// 创建测试相关的路由
	testRouter := router.Group("/test")
	{
		testRouter.GET("/", api.Welcome)
		testRouter.GET("/flush", func(context *gin.Context) {
			if _, err := data.FlushAll(); err != nil {
				println("Error when flushAll. " + err.Error())
			} else {
				println("Flushall succeed.")
			}
		})
	}

	// 启动秒杀功能的消费者（用来异步更新数据库）
	api.RunSecKillConsumer()

	return router
}