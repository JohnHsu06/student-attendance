package mysql

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("ERROR RECORD NOT FOUND")
)

var db = &gorm.DB{}

func InitMysql() {
	// 连接Mysql数据库
	dsn := "root:abc123@tcp(127.0.0.1:3306)/attendance_db?charset=utf8mb4&parseTime=True&loc=Local"
	//单独定义err,以避免调用 gorm.Open 时用:=导致db变量发生代码遮蔽的问题
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	//关闭错误输出

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(15)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Minute)
}
