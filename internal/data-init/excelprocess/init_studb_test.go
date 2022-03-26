package excelprocess

import (
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestInitStudentsDatabase(t *testing.T) {
	f, _ := excelize.OpenFile("全年级学生.xlsx")
	slice, _ := InitStudentsDatabase(f)
	for _, v := range slice {
		fmt.Println(v)
	}
}
