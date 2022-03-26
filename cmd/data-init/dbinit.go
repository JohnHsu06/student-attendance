package main

import (
	"fmt"
	"student-attendance/internal/data-init/excelprocess"
	"student-attendance/internal/pkg/model"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	model.StudentInfo
}

func main() {
	//设置学生年级
	var grade uint16 = 2020

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

	db.AutoMigrate(&Student{})

	f, err := excelize.OpenFile("全年级学生.xlsx")
	if err != nil {
		panic(err)
	}
	stuSlice, err := excelprocess.InitStudentsDatabase(f)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range stuSlice {
		student := Student{}
		student.StuGrade = grade
		student.StuClass = v.StuClass
		student.StuNumber = v.StuNumber
		student.StuName = v.StuName
		db.Create(&student)
	}

}
