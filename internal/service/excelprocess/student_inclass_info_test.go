package excelprocess

import (
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestGetStudentsInfo(t *testing.T) {
	f, _ := excelize.OpenFile("test.xlsx")
	rows, _ := f.GetRows("数据导出")
	idenStus, unidenStus, err := getStudentsInfo(rows[5:], 2020)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range idenStus {
		fmt.Printf("[%d] %v\n", k, v)
	}
	fmt.Printf("\n=========无法识别==========\n")
	for i, v := range unidenStus {
		fmt.Printf("[%d] %v\n", i, v)
	}
}

func TestGetStuEntryTime(t *testing.T) {
	str := []string{"2022/03/22 10:24:52", "2022/03/22 10:24:04", "2022/03/22 10:23:04",
		"2022/03/22 10:21:43", "2022/03/22 10:13:01", "2022/03/22 10:13:00", "2022/03/22 10:10:32",
		"2022/03/22 10:06:16", "2022/03/22 10:06:02", "2022/03/22 10:05:14", "2022/03/22 10:03:43",
		"2022/03/22 10:03:33", "2022/03/22 10:02:14"}
	for i, v := range str {
		t, err := getStuEntryTime(v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("[%v] %v\n", i, *t)
	}
}

func TestGetStuWatchDuration(t *testing.T) {
	str := []string{"48分钟", "49分钟", "48分钟", "不足一分钟"}
	for i, v := range str {
		res, err := getStuWatchDuration(v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("[%v] type:%T %v\n", i, res, res)
	}
}

func TestGetStudentNameAndNum(t *testing.T) {
	str := []string{"0407郭锋", "桂园中学初二7班周浩宇", "03郑伟杭", "马云", "九班张钰",
		"0321解畅然", "8班黄佳宁", "123456",
	}
	for i, v := range str {
		flag, stuName, stuClass, stuNum := getStuNameAndNum(v)
		if flag == false {
			fmt.Printf("[%v] flag:%t %v\n", i, flag, v)
		} else {
			fmt.Printf("[%v] flag:%t name:%s class:%v num:%v\n", i, flag, stuName, stuClass, stuNum)
		}
	}
}
