package mysql

import (
	"fmt"
	"testing"
)

func TestGetStuID(t *testing.T) {
	InitMysql()
	id, err := GetStuID(2020, 4, 1, "蔡华龙")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}
