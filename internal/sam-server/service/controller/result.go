package controller

import (
	"strconv"
	"student-attendance/internal/pkg/model"
	"student-attendance/internal/pkg/utils"
	"student-attendance/internal/sam-server/service/excelprocess"
	"student-attendance/internal/sam-server/store/mysql"
	"time"
)

type resultInfo struct {
	Ci             *model.CommitInfo
	Unidens        []*excelprocess.UnidentifiedStudentInfo
	AbsentStus     map[uint]*mysql.Student
	LateStus       map[uint]*excelprocess.IdentifiedStudentInfo
	EarlyLeaveStus map[uint]*excelprocess.IdentifiedStudentInfo
}

// makeResult 生成用于在html页面上显示的结果
func (r *resultInfo) makeResult() map[string]interface{} {
	Result := make(map[string]interface{})
	Result["Ci"] = r.makeMainInfo("2022-02-21")
	Result["Unidens"] = r.Unidens
	Result["ClassList"] = r.makeClassList()
	class := make([]map[string]interface{}, 0, 10)
	for _, classNum := range r.Ci.AttendanceClasses {
		class = append(class, r.makeClassInfo(classNum))
	}
	Result["Class"] = class
	return Result
}

// makeMainInfo根据Ci生成合适html页面显示的信息
func (r *resultInfo) makeMainInfo(firstWeekDay string) map[string]interface{} {
	ci := make(map[string]interface{})
	ci["Subject"] = utils.GetSubjectFromCode(r.Ci.Subject)
	ci["ExpectedArrivalNum"] = r.Ci.ExpectedArrivalNum
	ci["ActualArrivalNum"] = r.Ci.ActualArrivalNum
	ci["WatchNum"] = r.Ci.WatchNum
	ci["AttendanceRate"] = strconv.FormatFloat(r.Ci.AttendanceRate*100, 'f', 2, 64) + "%"
	ci["BroadcastTime"] = r.Ci.BroadcastTime.Format("2006-01-02 15:04:05")
	ci["ClassStartTime"] = r.Ci.ClassStartTime.Format("15:04")
	ci["ClassEndTime"] = r.Ci.ClassEndTime.Format("15:04")
	ci["TeacherName"] = r.Ci.TeacherName
	ci["ClassTheme"] = r.Ci.ClassTheme
	ci["AttendanceClasses"] = r.Ci.AttendanceClasses

	firstTeachingWeekDay, _ := time.Parse("2006-01-02", firstWeekDay)
	_, firstTeachingWeek := firstTeachingWeekDay.ISOWeek()
	ci["WeekNum"] = r.Ci.WeekNum - int8(firstTeachingWeek) + 1

	classDuration := r.Ci.ClassEndTime.Sub(*r.Ci.ClassStartTime)
	ci["ClassDuration"] = int(classDuration / time.Minute)
	return ci
}

// makeClassList返回各班班号字符串切片
func (r *resultInfo) makeClassList() []string {
	classList := make([]string, 0, 10)
	for _, class := range r.Ci.AttendanceClasses {
		classList = append(classList, strconv.Itoa(int(class)))
	}
	return classList
}

// makeClassInfo 根据班号生成各班结果
func (r *resultInfo) makeClassInfo(classNum int8) map[string]interface{} {
	//生成以问题(缺勤、迟到、早退，外加一项班号为key)的map
	problem := make(map[string]interface{})
	problem["ClassNum"] = classNum
	//该班缺勤人员
	absentList := make([]string, 0, 5)
	for id, stu := range r.AbsentStus {
		if stu.StuClass == classNum {
			absentList = append(absentList, stu.StuName)
			delete(r.AbsentStus, id)
		}
	}
	problem["AbsentList"] = absentList
	//该班迟到人员
	lateList := make(map[string]string)
	for id, stu := range r.LateStus {
		if stu.StuClass == classNum {
			entryTime := stu.EntryTime.Format("15:04:05")
			lateList[stu.StuName] = entryTime
			delete(r.LateStus, id)
		}
	}
	problem["LateList"] = lateList
	//该班早退人员
	earlyLeaveList := make(map[string]string)
	for id, stu := range r.EarlyLeaveStus {
		if stu.StuClass == classNum {
			watchTime := strconv.Itoa(int(stu.WatchDuration)) + "分钟"
			earlyLeaveList[stu.StuName] = watchTime
			delete(r.EarlyLeaveStus, id)
		}
	}
	problem["EarlyLeaveList"] = earlyLeaveList

	return problem
}
