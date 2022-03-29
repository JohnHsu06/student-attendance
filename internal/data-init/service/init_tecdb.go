package service

import (
	"strconv"
	"student-attendance/internal/pkg/model"
	"student-attendance/internal/pkg/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type teacher struct {
	gorm.Model
	model.TeacherInfo
}

func InitTeachersDatabase(f *excelize.File, db *gorm.DB, grade uint16) error {
	rows, err := f.GetRows("教师表")
	if err != nil {
		return err
	}
	db.AutoMigrate(&teacher{})

	tecRows := rows[1:]
	for _, row := range tecRows {
		tec := new(teacher)
		tec.Grade = grade
		tec.Subject = utils.GetSubjectCode(row[0])
		tec.Name = row[1]
		//获取是否是班主任
		if strToInt8(row[2]) != 0 {
			tec.ClassTeacher = strToInt8(row[2])
		}
		//获取是否是备课组长
		if strToInt8(row[3]) == 1 {
			tec.SubjectLeader = tec.Subject
		}
		//获取是否是级长
		if strToInt8(row[4]) == 1 {
			tec.GradeLeader = true
		}
		//获取是否是管理员
		if strToInt8(row[5]) == 1 {
			tec.Admin = true
		}
		//获取是否是超管
		if strToInt8(row[6]) == 1 {
			tec.SuperAdmin = true
		}
		db.Create(tec)
	}
	return nil
}

// strToInt8 将字符串转成int8格式
func strToInt8(str string) int8 {
	temp, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		panic(err)
	}
	return int8(temp)
}
