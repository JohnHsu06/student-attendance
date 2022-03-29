package model

import "time"

// StudentInfo是用于确认一个学生身份的基本信息
type StudentInfo struct {
	StuGrade  uint16 //如2021,表示2021级,即初一学生;同理，2020即初二,2019即初三
	StuClass  int8
	StuNumber int8
	StuName   string `gorm:"type:varchar(40)"`
}

// TeacherInfo是用于确认一个教师身份的基本信息
type TeacherInfo struct {
	Grade         uint16
	Subject       int8
	ClassTeacher  int8   //班主任,0代表不是,正数代表所在班级
	SubjectLeader int8   //备课组长,0代表不是,正数代表对应的学科代码
	GradeLeader   bool   //年级级长
	Admin         bool   //管理员
	SuperAdmin    bool   //超级管理员
	Name          string `gorm:"type:varchar(40)"`
}

// CommitInfo用于存储每次提交的课堂记录的主要信息
type CommitInfo struct {
	Subject            int8
	WeekNum            int8
	ExpectedArrivalNum uint16
	ActualArrivalNum   uint16
	WatchNum           uint16
	Grade              uint16
	AttendanceRate     string
	BroadcastTime      *time.Time
	ClassStartTime     *time.Time
	ClassEndTime       *time.Time
	TeacherName        string
	ClassTheme         string
	AttendanceClasses  []int8
}

// UploadInfo用于存储伴随考勤表一起提交的用户信息
type UploadInfo struct {
	User      string `form:"teacher-name"`
	Subject   string `form:"subject"`
	StartTime string `form:"start-time"`
	EndTime   string `form:"end-time"`
}
