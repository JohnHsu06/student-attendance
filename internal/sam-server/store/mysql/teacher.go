package mysql

import (
	"errors"
	"student-attendance/internal/pkg/code"
	"student-attendance/internal/pkg/model"
	"student-attendance/internal/pkg/utils"

	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	model.TeacherInfo
}

// GetTeacherByNameNSubject 根据姓名与科目在数据库中获取记录
func GetTeacherByNameNSubject(name, sub string) (*Teacher, error) {
	//cond结构体存储教师的姓名与科目,作为查询条件
	cond := new(Teacher)
	cond.Name = name
	tec := new(Teacher)
	var err error

	subCode := utils.GetSubjectCode(sub)
	//根据是否是班会课分开判断
	if subCode == code.ClassMeeting { //是班会课的情况
		res := db.Where("class_teacher <> ? AND name = ?", "0", name).First(tec)
		err = res.Error
	} else { //不是班会课的情况
		cond.Subject = subCode
		res := db.Where(cond).First(tec)
		err = res.Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return tec, nil
}
