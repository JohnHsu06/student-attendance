package service

import (
	"strconv"
	"student-attendance/internal/pkg/model"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type student struct {
	gorm.Model
	model.StudentInfo
}

// InitStudentsDatabase 函数从Excel表格中读取出全年级的学生数据
func InitStudentsDatabase(f *excelize.File, db *gorm.DB, grade uint16) error {
	sheetList := f.GetSheetList()
	stuSlice := make([]*model.StudentInfo, 0, 800)
	//遍历每一个班的工作表
	for _, sheetIndex := range sheetList {
		rows, err := f.GetRows(sheetIndex)
		if err != nil {
			return err
		}
		aClassStus := GetStudentsInfo(rows[1:])
		stuSlice = append(stuSlice, aClassStus...)
	}

	db.AutoMigrate(&student{})
	for _, v := range stuSlice {
		student := student{}
		student.StuGrade = grade
		student.StuClass = v.StuClass
		student.StuNumber = v.StuNumber
		student.StuName = v.StuName
		db.Create(&student)
	}
	return nil
}

// GetStudentsInfo 函数从传来的学生二维数组中读取学生的班别、学号和姓名
func GetStudentsInfo(stuRows [][]string) []*model.StudentInfo {
	aClassStusSlice := make([]*model.StudentInfo, 0, 60)
	for _, row := range stuRows {
		if row == nil {
			continue
		}
		stu := &model.StudentInfo{}
		class, num := GetStudentClassAndNum(row[0])
		stu.StuClass = class
		stu.StuNumber = num
		stu.StuName = row[1]
		aClassStusSlice = append(aClassStusSlice, stu)
	}
	return aClassStusSlice
}

func GetStudentClassAndNum(numStr string) (class, num int8) {
	tempClass, _ := strconv.ParseInt(numStr[:2], 10, 8)
	tempNum, _ := strconv.ParseInt(numStr[2:], 10, 8)
	class = int8(tempClass)
	num = int8(tempNum)
	return
}
