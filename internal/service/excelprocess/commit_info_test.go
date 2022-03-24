package excelprocess

import (
	"fmt"
	"student-attendance/internal/store"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestGetCommitInfo(t *testing.T) {
	f, _ := excelize.OpenFile("test.xlsx")
	ci := new(store.CommitInfo)
	getCommitInfo(f, ci)
	fmt.Println(ci)
}
