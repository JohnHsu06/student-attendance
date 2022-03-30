package excelprocess

import (
	"encoding/json"
	"strconv"
	"strings"
	"student-attendance/internal/sam-server/store/mysql"
	"time"
)

// getStudentsInfo 从工作表中读出每一位学生的基本信息和课堂信息
func getStudentsInfo(stuRows [][]string, stuGrade uint16) (map[uint]*IdentifiedStudentInfo,
	[]*UnidentifiedStudentInfo, error) {
	identifiedStudents := make(map[uint]*IdentifiedStudentInfo, len(stuRows))
	unidentifiedStudents := make([]*UnidentifiedStudentInfo, 0, 10)

	//循环工作表学生部分的遍历每一行
	for _, row := range stuRows {
		if row == nil { //跳过可能出现的空行
			continue
		}
		//从E列的单元格读取学生的腾讯课堂ID
		tencentID, err := strconv.ParseUint(row[4], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		//从G列读取学生进入直播间的时间
		entryTime, err := getStuEntryTime(row[6])
		if err != nil {
			return nil, nil, err
		}
		//从H列读取学生观看直播时长
		watchDuration, err := getStuWatchDuration(row[7])
		if err != nil {
			return nil, nil, err
		}
		//推算学生退出直播间的时间
		leaveTime := calLeaveTime(entryTime, watchDuration)
		//从K列读取学生是否有签到
		signIn := true
		if len(row) == 11 { //代表该老师发起了签到
			signIn, err = checkSignIn(row[10])
			if err != nil {
				return nil, nil, err
			}
		}

		//从D列的单元格中读取学生的班别、学号和姓名；并根据是否读取成功进行分类
		nameStr := row[3]
		flag, stuName, stuClass, stuNum := getStuNameAndNum(nameStr)
		//如果出现无法读取出班别、学号和姓名的情况
		if !flag {
			unidenStu := makeUnidenStuRec(nameStr, entryTime, watchDuration, leaveTime, signIn, tencentID)
			unidentifiedStudents = append(unidentifiedStudents, unidenStu)
			continue
		}

		//以下为可以正确读出学生的班别、学号和姓名的情况
		//查询数据库，获取学生内部ID以作为map的key,以避免以姓名为key出现学生重名的情况
		sID, err := mysql.GetStuID(stuGrade, stuClass, stuNum, stuName)
		//防御性检查，避免出现表格中学生姓名符合格式但数据库中无法查询到记录的情况
		//如无法匹配数据库，将学生加入到无法识别名单中
		if err == mysql.ErrRecordNotFound || sID == 0 {
			unidenStu := makeUnidenStuRec(nameStr, entryTime, watchDuration, leaveTime, signIn, tencentID)
			unidentifiedStudents = append(unidentifiedStudents, unidenStu)
			continue
		}
		//查看是否有学生跳转登录腾讯课堂导致存在多条记录的情况
		cmp, ok := identifiedStudents[sID]
		if ok && cmp.WatchDuration > watchDuration { //已有课堂记录且观看时长比新记录更长,直接忽略新纪录
			continue
		}
		idenStu := makeIdenStuRec(entryTime, watchDuration, leaveTime, signIn,
			tencentID, stuName, stuClass, stuNum, stuGrade)
		identifiedStudents[sID] = idenStu
	}
	return identifiedStudents, unidentifiedStudents, nil
}

// getStuEntryTime 读取学生进入直播间的时间
func getStuEntryTime(str string) (*time.Time, error) {
	str += " +0800 CST"
	entryTime, err := time.Parse("2006/01/02 15:04:05 -0700 MST", str)
	if err != nil {
		return nil, err
	}
	return &entryTime, nil
}

// getStuWatchDuration 读取学生观看直播的时长
func getStuWatchDuration(str string) (uint16, error) {
	var watchDuration uint16
	if str == "不足一分钟" {
		watchDuration = 1
	} else {
		str = strings.TrimRight(str, "分钟")
		tempDuraton, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return 0, err
		}
		watchDuration = uint16(tempDuraton)
	}
	return watchDuration, nil
}

// calLeaveTime 根据学生进入课堂的时间与观看时长推算出其退出时间
func calLeaveTime(entryTime *time.Time, d uint16) *time.Time {
	watchDuration := time.Duration(d) * time.Minute
	leaveTime := entryTime.Add(watchDuration)
	return &leaveTime
}

// checkSignIn 检查学生是否有签到
func checkSignIn(signInStr string) (bool, error) {
	check := make(map[string]string)
	err := json.Unmarshal([]byte(signInStr), &check)
	if err != nil {
		return false, err
	}
	for _, v := range check {
		if v != "是" { // v != "是"代表没有签到
			return false, nil // 返回bool值为false代表没有签到
		}
	}
	return true, nil //返回bool值为true代表有签到
}

// getStuNameAndNum 从工作表D列读取学生的班别、学号和姓名。若返回的布尔值为false,说明无法从单元格中读出上述信息
func getStuNameAndNum(str string) (flag bool, stuName string, stuClass int8, stuNum int8) {
	//读取中文名
	str = strings.TrimSpace(str)
	stuName = getCn(str)
	if stuName == "" { //无法读出中文姓名
		return
	}

	//读取班别与学号
	r := []rune(str)
	if len(r) <= 4 { //字符串长度不足4位
		return
	}
	for i := 0; i < 4; i++ { //前4位是否为数字
		if r[i] < 48 || r[i] > 57 {
			return
		}
	}
	tempClass, err := strconv.ParseInt(str[:2], 10, 8)
	if err != nil {
		return
	}
	if tempClass == 0 || tempClass > 15 { //班别数据合理性校验
		return
	}
	stuClass = int8(tempClass)
	tempNum, err := strconv.ParseInt(str[2:4], 10, 8)
	if err != nil {
		return
	}
	if tempNum == 0 || tempNum > 60 { //学号数据合理性校验
		return
	}
	stuNum = int8(tempNum)

	flag = true
	return
}

// getCn 获取字符串中的中文部分
func getCn(str string) (cnStr string) {
	r := []rune(str)
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnStr += string(r[i])
		}
	}
	return
}

// makeUnidenStuRec 传入一个无法识别学生的信息返回对应的一个结构体
func makeUnidenStuRec(nameStr string, entryTime *time.Time, watchDuration uint16, leaveTime *time.Time,
	signIn bool, tencentID uint64) *UnidentifiedStudentInfo {
	unidenStu := &UnidentifiedStudentInfo{
		NameStr:       nameStr,
		EntryTime:     entryTime,
		WatchDuration: watchDuration,
		LeaveTime:     leaveTime,
		SignIn:        signIn,
		TencentID:     tencentID,
	}
	return unidenStu
}

// makeIdenStuRec 传入一个可以识别学生的信息返回一个对应的结构体
func makeIdenStuRec(entryTime *time.Time, watchDuration uint16, leaveTime *time.Time, signIn bool, tencentID uint64,
	stuName string, stuClass int8, stuNum int8, stuGrade uint16) *IdentifiedStudentInfo {
	idenStu := &IdentifiedStudentInfo{
		EntryTime:     entryTime,
		WatchDuration: watchDuration,
		LeaveTime:     leaveTime,
		SignIn:        signIn,
		TencentID:     tencentID,
	}
	idenStu.StuName = stuName
	idenStu.StuClass = stuClass
	idenStu.StuNumber = stuNum
	idenStu.StuGrade = stuGrade
	return idenStu
}
