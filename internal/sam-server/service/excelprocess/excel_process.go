package excelprocess

import (
	"errors"
	"log"
	"student-attendance/internal/pkg/model"
	"time"

	"github.com/xuri/excelize/v2"
)

type excelResult struct {
	// CommitInfo           *model.CommitInfo
	IdentifiedStudents   map[uint]*IdentifiedStudentInfo
	UnidentifiedStudents []*UnidentifiedStudentInfo
}

// IdentifiedStudentInfo 用于存储Excel表格中可以辨别身份学生的课堂信息
type IdentifiedStudentInfo struct {
	model.StudentInfo
	EntryTime     *time.Time
	WatchDuration uint16
	LeaveTime     *time.Time
	SignIn        bool
	TencentID     uint64
}

// UnidentifiedStudentInfo 用于存储Excel表格中无法辨别身份学生的课堂信息
type UnidentifiedStudentInfo struct {
	NameStr       string
	EntryTime     *time.Time
	WatchDuration uint16
	LeaveTime     *time.Time
	SignIn        bool
	TencentID     uint64
}

var (
	ErrIncompleteExcelFile = errors.New("CLASSINFO EXCEL FILE LACKS IMPORTANT INFORMATION")
)

// 腾讯课堂导出的学生考勤数据Excel表格中的默认工作表名
var sheetName = "数据导出"

// ReadClassInfo从上传的考勤Excel表格中读出课堂基本信息和学生信息
func ReadClassInfo(filePath string, ci *model.CommitInfo) (*excelResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()

	//判断课堂信息表中是否含有"数据导出"工作表
	if n := f.GetSheetIndex(sheetName); n == -1 {
		return nil, err
	}

	//获取课堂基本信息
	if err = getCommitInfo(f, ci); err != nil {
		return nil, err
	}
	//获取所有到堂学生的信息
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	stuRows := rows[5:] //学生信息是从工作表的第6行开始的
	identifiedStudents, unidentifiedStudents, err := getStudentsInfo(stuRows, ci.Grade)
	if err != nil {
		return nil, err
	}
	excelRes := &excelResult{
		// CommitInfo:           ci,
		IdentifiedStudents:   identifiedStudents,
		UnidentifiedStudents: unidentifiedStudents,
	}
	return excelRes, nil
}
