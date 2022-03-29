package mysql

import (
	"fmt"
	"testing"
)

func TestGetStusByGradeNClass(t *testing.T) {
	InitMysql()
	stus := GetStusByGradeNClass(2020, 4)
	for _, v := range stus {
		fmt.Println(v)
	}
}
