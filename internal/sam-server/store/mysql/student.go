package mysql

import (
	"errors"
	"fmt"
	"student-attendance/internal/pkg/model"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	model.StudentInfo
}

// GetStuID 根据学生的年级、班别、学号和姓名从数据库中获取学生的软件内部ID
func GetStuID(grade uint16, class int8, num int8, name string) (uint, error) {
	// cond结构体存储学生的基本信息，作为查询条件
	cond := new(Student)
	cond.StuGrade = grade
	cond.StuClass = class
	cond.StuNumber = num
	cond.StuName = name

	stu := new(Student)
	res := db.Where(cond).First(stu)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, ErrRecordNotFound
	}
	return stu.ID, nil
}

// GetStusByGradeNClass 根据年级和班别返回该班的所有学生
func GetStusByGradeNClass(grade uint16, class int8) []Student {
	// cond结构体存储学生的基本信息，作为查询条件
	cond := new(Student)
	cond.StuGrade = grade
	cond.StuClass = class

	students := make([]Student, 0, 70)
	res := db.Where(cond).Find(&students)
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	return students
}
