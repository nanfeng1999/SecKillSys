package api

import (
	"SecKillSys/api/dbService"
	"SecKillSys/model"
	"log"
)

type secKillMessage struct {
	username string
	coupon model.Coupon
}

const maxMessageNum = 20000
var SecKillChannel = make(chan secKillMessage, maxMessageNum) //有缓存的channel

func seckillConsumer(){
	for {
		message := <- SecKillChannel
		log.Println("Got one message: " + message.username)

		username := message.username//抢购成功的用户的用户名
		sellerName := message.coupon.Username //优惠券的商家名
		couponName := message.coupon.CouponName//优惠券名

		var err error
		err = dbService.UserHasCoupon(username, message.coupon) //用户优惠券数+1
		if err != nil {
			println("Error when inserting user's coupon. " + err.Error())
		}
		err = dbService.DecreaseOneCouponLeft(sellerName, couponName)//优惠券库存自减1
		if err != nil {
			println("Error when decreasing coupon left. " + err.Error())
		}
	}
}

var isConsumerRun = false
func RunSecKillConsumer() {
	// Only Run one consumer.
	if !isConsumerRun {
		go seckillConsumer()  //开启一个消费者goroutune，作用是接收redis的改动信息，更新数据库
		isConsumerRun = true
	}
}