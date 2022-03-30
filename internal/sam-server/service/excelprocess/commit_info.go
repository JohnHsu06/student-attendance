package excelprocess

import (
	"strconv"
	"strings"
	"student-attendance/internal/pkg/model"
	"time"

	"github.com/xuri/excelize/v2"
)

// getCommitInfo 函数获取课堂的基本信息
func getCommitInfo(f *excelize.File, ci *model.CommitInfo) error {
	//获取直播(开始)时间
	broadcastTimeStr, err := f.GetCellValue(sheetName, "A3")
	if err != nil {
		return ErrIncompleteExcelFile
	}
	tempTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", broadcastTimeStr)
	if err != nil {
		return ErrIncompleteExcelFile
	}
	ci.BroadcastTime = &tempTime
	//获取直播时间是当年的第几周
	_, weekNum := tempTime.ISOWeek()
	ci.WeekNum = int8(weekNum)

	//获取观看直播的人数
	watchNumStr, err := f.GetCellValue(sheetName, "C3")
	if err != nil {
		return ErrIncompleteExcelFile
	}
	watchNumStr = strings.TrimSuffix(watchNumStr, ".0")
	tempNum, err := strconv.ParseUint(watchNumStr, 10, 16)
	if err != nil {
		return ErrIncompleteExcelFile
	}
	ci.WatchNum = uint16(tempNum)

	//获取课程主题
	ci.ClassTheme, err = f.GetCellValue(sheetName, "C6")
	if err != nil {
		return ErrIncompleteExcelFile
	}

	return nil
}
