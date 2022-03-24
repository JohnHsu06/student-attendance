package excelprocess

import (
	"strconv"
	"strings"
	"time"
)

// getStudentsInfo 函数从工作表中读出每一位学生的基本信息和课堂信息
func getStudentsInfo(stuRows [][]string, stuGrade uint16) (map[uint]*IdentifiedStudentInfo,
	[]*UnidentifiedStudentInfo, error) {
	identifiedStudents := make(map[uint]*IdentifiedStudentInfo, len(stuRows))
	unidentifiedStudents := make([]*UnidentifiedStudentInfo, 0, 10)
	var tempCount uint = 1 //临时充当学生的内部ID，待数据库内容补全后应删除

	//循环工作表学生部分的遍历每一行
	for _, row := range stuRows {
		if row == nil { //跳过可能出现的空行
			continue
		}
		//从E列的单元格读取学生的腾讯课堂ID
		tencentID, err := strconv.ParseUint(row[4], 10, 64)
		if err != nil {
			return nil, nil, ErrIncompleteExcelFile
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

		//从D列的单元格中读取学生的班别、学号和姓名；并根据是否读取成功进行分类
		nameStr := row[3]
		flag, stuName, stuClass, stuNum := getStuNameAndNum(nameStr)
		if !flag { //无法读取出班别、学号和姓名的情况
			unidenStu := &UnidentifiedStudentInfo{
				NameStr:       nameStr,
				EntryTime:     entryTime,
				WatchDuration: watchDuration,
				TencentID:     tencentID,
			}
			unidentifiedStudents = append(unidentifiedStudents, unidenStu)
		} else { //正确读出学生的班别、学号和姓名的情况
			idenStu := &IdentifiedStudentInfo{
				EntryTime:     entryTime,
				WatchDuration: watchDuration,
				TencentID:     tencentID,
			}
			idenStu.StuName = stuName
			idenStu.StuClass = stuClass
			idenStu.StuNumber = stuNum
			idenStu.StuGrade = stuGrade
			//这里应该有查询数据库，获取学生内部ID以作为map(可避免出现学生重名的情况)的key的步骤，待数据库内容写好后补全
			//且应该补全重复登录导致的时长变短
			identifiedStudents[tempCount] = idenStu
			tempCount++
			//
		}
	}
	return identifiedStudents, unidentifiedStudents, nil
}

// getStuEntryTime 函数读取学生进入直播间的时间
func getStuEntryTime(str string) (*time.Time, error) {
	str += " +0800 CST"
	entryTime, err := time.Parse("2006/01/02 15:04:05 -0700 MST", str)
	if err != nil {
		return nil, ErrIncompleteExcelFile
	}
	return &entryTime, nil
}

// getStuWatchDuration 函数读取学生观看直播的时长
func getStuWatchDuration(str string) (uint16, error) {
	var watchDuration uint16
	if str == "不足一分钟" {
		watchDuration = 1
	} else {
		str = strings.TrimRight(str, "分钟")
		tempDuraton, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return 0, ErrIncompleteExcelFile
		}
		watchDuration = uint16(tempDuraton)
	}
	return watchDuration, nil
}

// getStuNameAndNum 函数从工作表D列读取学生的班别、学号和姓名。若返回的布尔值为false,说明无法从单元格中读出上述信息
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

// getCn 函数获取字符串中的中文部分
func getCn(str string) (cnStr string) {
	r := []rune(str)
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnStr += string(r[i])
		}
	}
	return
}
