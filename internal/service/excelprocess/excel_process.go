package excelprocess

import (
	"errors"
	"student-attendance/internal/store"
	"time"

	"github.com/xuri/excelize/v2"
)

type excelResult struct {
	CommitInfo           *store.CommitInfo
	IdentifiedStudents   map[uint]*IdentifiedStudentInfo
	UnidentifiedStudents []*UnidentifiedStudentInfo
}

// IdentifiedStudentInfo 用于存储Excel表格中可以辨别身份学生的课堂信息
type IdentifiedStudentInfo struct {
	store.StudentInfo
	EntryTime     *time.Time
	WatchDuration uint16
	TencentID     uint64
}

// UnidentifiedStudentInfo 用于存储Excel表格中无法辨别身份学生的课堂信息
type UnidentifiedStudentInfo struct {
	NameStr       string
	EntryTime     *time.Time
	WatchDuration uint16
	TencentID     uint64
}

var (
	ErrIncompleteExcelFile = errors.New("CLASSINFO EXCEL FILE LACKS IMPORTANT INFORMATION")
)

// 腾讯课堂导出的学生考勤数据Excel表格中的默认工作表名
var sheetName = "数据导出"

// readClassInfo从上传的考勤Excel表格中读出课堂基本信息和学生信息
func readClassInfo(fileName string, ci *store.CommitInfo) (*excelResult, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//判断课堂信息表中是否含有"数据导出"工作表
	if n := f.GetSheetIndex(sheetName); n == -1 {
		return nil, ErrIncompleteExcelFile
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
		CommitInfo:           ci,
		IdentifiedStudents:   identifiedStudents,
		UnidentifiedStudents: unidentifiedStudents,
	}
	return excelRes, nil
}
