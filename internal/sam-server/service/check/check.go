package check

import (
	"errors"
	"fmt"
	"student-attendance/internal/pkg/model"
	"time"
)

var (
	ErrIncompleteInfo = errors.New("INCOMPLETE INFORMATION")
	ErrIllogicalTime  = errors.New("ILLOGICAL TIME")
)

// CheckUploadInfo 检查用户上传的基本信息是否合理
func CheckUploadInfo(ui *model.UploadInfo) error {
	//检查是否有上传空信息
	if ui.User == "" || ui.Subject == "" || ui.StartTime == "" || ui.EndTime == "" {
		return ErrIncompleteInfo
	}
	//检查上下课时间是否合理
	if err := timeCheck(ui.StartTime, ui.EndTime); err != nil {
		return err
	}

	return nil
}

// timeCheck 检查用户上传的上课时间与下课时间是否合理
func timeCheck(startTimeStr, endTimeStr string) error {
	timeForm := "15:04"
	startTime, err := time.Parse(timeForm, startTimeStr)
	if err != nil {
		return ErrIllogicalTime
	}
	endTime, err := time.Parse(timeForm, endTimeStr)
	if err != nil {
		return ErrIllogicalTime
	}
	//判断结束时间是否在开始时间之前
	if endTime.Before(startTime) {
		return ErrIllogicalTime
	}
	//判断课堂时间是否过短,当前的判断标准是15分钟
	d := endTime.Sub(startTime)
	fmt.Println(d)
	if int(d/time.Minute) < 15 {
		return ErrIllogicalTime
	}
	return nil
}
