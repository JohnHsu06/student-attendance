package controller

import (
	"fmt"
	"student-attendance/internal/pkg/model"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tmpT := time.Now()
	ci := new(model.CommitInfo)
	ci.BroadcastTime = &tmpT

	res := parseTime("09:30", ci)
	fmt.Println(res)

}

func TestGetAttendanceRate(t *testing.T) {
	var ac uint16 = 253
	var ep uint16 = 270
	fmt.Println(getAttendanceRate(ac, ep))
}
