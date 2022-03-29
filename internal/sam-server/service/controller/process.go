package controller

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"student-attendance/internal/pkg/model"
	"student-attendance/internal/pkg/utils"
	"student-attendance/internal/sam-server/service/excelprocess"
	"student-attendance/internal/sam-server/store/mysql"
	"time"
)

func Process(dst string, ui *model.UploadInfo) (*model.CommitInfo, error) {
	//创建一个ci用于存储这一次记录提交的主要信息
	ci := new(model.CommitInfo)
	ci.Subject = utils.GetSubjectCode(ui.Subject)
	ci.TeacherName = ui.User
	//通过查询教师表拿到年级信息
	t, _ := mysql.GetTeacherByNameNSubject(ui.User, ui.Subject)
	ci.Grade = t.Grade

	//调用Excel处理函数,读取学生信息,同时将Excel表里的课堂基本信息写入ci
	excelRes, err := excelprocess.ReadClassInfo(dst, ci)
	if err != nil {
		return nil, err
	}
	//根据Excel处理后ci里的直播日期，将直播开始与结束时间转换成time.Time形式
	ci.ClassStartTime = parseTime(ui.StartTime, ci)
	ci.ClassEndTime = parseTime(ui.EndTime, ci)

	//从可识别学生中获取上课的班级列表
	ci.AttendanceClasses = getClasses(excelRes.IdentifiedStudents)
	//从上课的班级列表获取应到学生名单、应到学生人数
	expectedStu := getExpectedStudents(ci.Grade, ci.AttendanceClasses)
	ci.ExpectedArrivalNum = uint16(len(expectedStu))
	//应到学生名单与可识别学生名单比对得到缺勤学生与实际到达人数,此操作后应到学生名单里的人员为缺勤名单
	ci.ActualArrivalNum = getAbsentStusNArrivalNum(expectedStu, excelRes.IdentifiedStudents)
	//从实到学生和应到学生得到学生出勤率
	ci.AttendanceRate = getAttendanceRate(ci.ActualArrivalNum, ci.ExpectedArrivalNum)

	return ci, nil
}

// parseTime 将ui中的短格式日期字符串与excel中读出的课堂时间结合解析出可以用于时间比较的time.Time格式
func parseTime(timeStr string, ci *model.CommitInfo) *time.Time {
	preTimeStr := ci.BroadcastTime.Format("2006-01-02")
	timeStr = strings.Join([]string{preTimeStr, timeStr, "+0800 CST"}, " ")
	res, err := time.Parse("2006-01-02 15:04 -0700 MST", timeStr)
	if err != nil {
		fmt.Println(err)
	}
	return &res
}

// getClasses 从可识别学生map中得到上课的班级列表
func getClasses(IdenStudents map[uint]*excelprocess.IdentifiedStudentInfo) []int8 {
	count := make(map[int]uint8)
	//遍历所有可识别学生,统计各班出现的人数
	for _, stu := range IdenStudents {
		count[int(stu.StuClass)]++
	}
	//学生出现人数大于15人的班级纳入班级列表中
	temp := make([]int, 0, 20)
	for class, stuCount := range count {
		if stuCount > 15 {
			temp = append(temp, class)
		}
	}
	//升序排序
	sort.Ints(temp)
	classes := make([]int8, 0, 20)
	for _, v := range temp {
		classes = append(classes, int8(v))
	}
	return classes
}

// getExpectedStudents 根据班级列表查询数据库得到所有应到的学生名单
func getExpectedStudents(grade uint16, classes []int8) map[uint]*mysql.Student {
	expectedStus := make(map[uint]*mysql.Student)
	for _, class := range classes {
		//根据年级和班别查询该班所有学生
		aClassStu := mysql.GetStusByGradeNClass(grade, class)
		//将该班所有学生加入应到学生名单
		for _, stu := range aClassStu {
			expectedStus[stu.ID] = &stu
		}
	}
	return expectedStus
}

// getAbsentStusNArrivalNum 通过应到学生名单与可识别学生名单比对得到缺勤学生名单与实到人数
func getAbsentStusNArrivalNum(expectedStus map[uint]*mysql.Student,
	IdenStus map[uint]*excelprocess.IdentifiedStudentInfo) uint16 {
	var arrivalCount uint16
	//遍历可识别学生,删除应到学生名单中的对应记录
	for ID := range IdenStus {
		_, ok := expectedStus[ID]
		if ok { //该学生已到
			delete(expectedStus, ID)
			arrivalCount++
		}
	}
	return arrivalCount
}

// getAttendanceRate 根据ci中实到人数与应到人数的比值算出出勤率
func getAttendanceRate(actualArrivalNum, expectedArrivalNum uint16) string {
	attendanceRate := float64(actualArrivalNum) / float64(expectedArrivalNum)
	return strconv.FormatFloat(attendanceRate*100, 'f', 2, 64) + "%"
}
