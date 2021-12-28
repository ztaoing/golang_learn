package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "root:root@tcp(192.168.0.104:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	// 设置全局的logger，这个logger在执行每个SQL语句的时候会打印每一行SQL语句
	sqlLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,   //慢SQL阀值
		LogLevel:                  logger.Silent, // 日志的级别
		IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录没有找到）的错误
		Colorful:                  false,
	})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})
	if err != nil {
		panic(err)
	}
	// 新建会话模式

}
