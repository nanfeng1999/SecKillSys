package data

import (
	"SecKillSys/conf"
	"fmt"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

// 开启redis连接池
func initRedisConnection(config conf.AppConfig) {
	client = redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Address,
		Password: config.App.Redis.Password, // It's ok if password is "".
		DB:       0,                         // use default DB
	})

	if _, err := FlushAll(); err != nil {
		println("Error when flushAll. " + err.Error())
	}
}

// 用于测试
func FlushAll() (string, error) {
	return client.FlushAll().Result()
}

func PrepareScript(script string) string {
	// 检验脚本是否存在
	scriptsExists, err := client.ScriptExists(script).Result()
	if err != nil {
		panic("Failed to check if script exists: " + err.Error())
	}

	// 如果不存在 那么载入脚本
	fmt.Println(scriptsExists)
	if !scriptsExists[0] {
		scriptSHA, err := client.ScriptLoad(script).Result()
		if err != nil {
			panic("Failed to load script " + script + " err: " + err.Error())
		}
		// 返回脚本的哈希
		return scriptSHA
	}

	// 如果脚本存在直接返回空
	print("Script Exists.")
	return ""
}

func EvalSHA(sha string, args []string) (interface{}, error) {
	val, err := client.EvalSha(sha, args).Result()
	if err != nil {
		print("Error executing evalSHA... " + err.Error())
		return nil, err
	}
	return val, nil
}

// redis operation SET
func SetForever(key string, value interface{}) (string, error) {
	val, err := client.Set(key, value, 0).Result() // expiration表示无过期时间
	return val, err
}

// redis operation hmset
func SetMapForever(key string, field map[string]interface{}) (bool, error) {
	return client.HMSet(key, field).Result()
}

// redis operation hmget
func GetMap(key string, fields ...string) ([]interface{}, error) {
	return client.HMGet(key, fields...).Result()
}

// redis SADD
func SetAdd(key string, field string) (int64, error) {
	return client.SAdd(key, field).Result()
}

// redis SISMEMBER
func SetIsMember(key string, field string) (bool, error) {
	return client.SIsMember(key, field).Result()
}

// redis SMEMBERS
func GetSetMembers(key string) ([]string, error) {
	return client.SMembers(key).Result()
}
