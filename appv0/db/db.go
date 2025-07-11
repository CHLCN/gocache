package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewDb() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root", "123456", "127.0.0.1:3306", "gocache")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	DB = conn
}

var Rdb *redis.Client

func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

	//测试链接是否正常
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("无法连接到 Redis 服务器: %v\n", err)
		return
	}

	Rdb = rdb
}
