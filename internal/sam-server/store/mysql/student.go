package mysql

import (
	"errors"
	"student-attendance/internal/pkg/model"

	"gorm.io/gorm"
)

type student struct {
	gorm.Model
	model.StudentInfo
}

// GetStuID 根据学生的年级、班别、学号和姓名从数据库中获取学生的软件内部ID
func GetStuID(grade uint16, class int8, num int8, name string) (uint, error) {
	// cond结构体存储学生的基本信息，作为查询条件
	cond := new(student)
	cond.StuGrade = grade
	cond.StuClass = class
	cond.StuNumber = num
	cond.StuName = name

	stu := new(student)
	res := db.Where(cond).First(stu)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, ErrRecordNotFound
	}
	return stu.ID, nil

}
