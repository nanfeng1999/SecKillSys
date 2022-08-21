package data

import (
	"SecKillSys/conf"
	"SecKillSys/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var Db *gorm.DB

func initMysql(config conf.AppConfig) {
	fmt.Println("Load dbService config...")

	// 设置连接相关的参数
	dbType := config.App.Database.Type
	usr := config.App.Database.User
	pwd := config.App.Database.Password
	address := config.App.Database.Address
	dbName := config.App.Database.DbName
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		usr, pwd, address, dbName)

	//创建一个数据库的连接，因为docker中的mysql服务启动时延，一开始需要尝试重试连接
	fmt.Println("Init dbService connections...")
	var err error
	for Db, err = gorm.Open(dbType, dbLink); err != nil; Db, err = gorm.Open(dbType, dbLink) {
		log.Println("Failed to connect database: ", err.Error())
		log.Println("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	// 初始化数据库
	user := model.User{}      // 定义用户
	coupon := &model.Coupon{} // 定义优惠券

	// 设置连接池连接数
	Db.DB().SetMaxOpenConns(config.App.Database.MaxOpen)
	Db.DB().SetMaxIdleConns(config.App.Database.MaxIdle)

	// 创建表
	tables := []interface{}{user, coupon}

	for _, table := range tables {
		// 如果不存在表的话 那么自动建表
		if !Db.HasTable(table) {
			Db.AutoMigrate(table)
		}
	}

	// 删除所有记录
	if config.App.FlushAllForTest {
		println("FlushAllForTest is true. Delete records of all tables.")
		for _, table := range tables {
			Db.Delete(table)
		}
	}

	// 创建唯一索引
	Db.Model(user).AddUniqueIndex("username_index", "username")    // 用户的用户名唯一
	Db.Model(coupon).AddUniqueIndex("coupon_index", "coupon_name") // 优惠券的(用户名, 优惠券名)唯一

	println("---Mysql connection is initialized.---")
}
