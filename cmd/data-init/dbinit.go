package main

import (
	"fmt"
	"student-attendance/internal/data-init/service"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//设置年级
	var grade uint16 = 2020

	//初始化数据库连接
	dsn := "root:abc123@tcp(127.0.0.1:3306)/attendance_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	//初始化学生记录表
	stuFile, err := excelize.OpenFile("全年级学生.xlsx")
	if err != nil {
		panic(err)
	}
	err = service.InitStudentsDatabase(stuFile, db, grade)
	if err != nil {
		fmt.Println(err)
	}

	//初始化教师记录表
	tecFile, err := excelize.OpenFile("全年级老师.xlsx")
	if err != nil {
		panic(err)
	}
	err = service.InitTeachersDatabase(tecFile, db, grade)
	if err != nil {
		fmt.Println(err)
	}
}
