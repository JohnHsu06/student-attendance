package main

import (
	"student-attendance/internal/sam-server/router"
	"student-attendance/internal/sam-server/store/mysql"
)

func main() {
	mysql.InitMysql()
	r := router.SetupRouter()
	r.Run(":317")
}
