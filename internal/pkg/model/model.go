package model

import "time"

// StudentInfo是用于确认一个学生身份的基本信息
type StudentInfo struct {
	StuGrade  uint16 //如2021,表示2021级,即初一学生;同理，2020即初二,2019即初三
	StuClass  int8
	StuNumber int8
	StuName   string `gorm:"type:varchar(40)"`
}

// CommitInfo用于存储每次提交的课堂记录的主要信息
type CommitInfo struct {
	Subject            int8
	WeekNum            int8
	ExpectedArrivalNum uint16
	ActualArrivalNum   uint16
	WatchNum           uint16
	Grade              uint16
	AttendanceRate     float64
	BroadcastTime      *time.Time
	ClassStartTime     *time.Time
	ClassEndTime       *time.Time
	TeacherName        string
	ClassTheme         string
	AttendanceClasses  []int8
}
