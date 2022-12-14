package main

import (
	"SecKillSys/data"
	"SecKillSys/engine"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

const port = 20080

func main() {
	router := engine.SeckillEngine() //路由跳转都写在这里
	defer data.Close()               // 关闭数据库

	go func() { //可视化性能测试
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9876", nil))
	}()

	// 监听该端口
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		println("Error when running server. " + err.Error())
	}
}
