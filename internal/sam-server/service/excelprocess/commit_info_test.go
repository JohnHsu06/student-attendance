package excelprocess

import (
	"fmt"
	"student-attendance/internal/pkg/model"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestGetCommitInfo(t *testing.T) {
	f, _ := excelize.OpenFile("test.xlsx")
	ci := new(model.CommitInfo)
	getCommitInfo(f, ci)
	fmt.Println(ci)
}
