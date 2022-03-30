package router

import (
	"io"
	"os"
	"student-attendance/internal/sam-server/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	// 记录到文件。
	f, _ := os.Create("gin.log")
	// 同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	//装载静态文件
	r.Static("/static", "../../web/static/") //调试路径
	// r.Static("/static", "web/static/") //部署路径

	//渲染模板
	r.LoadHTMLGlob("../../web/templates/*") //调试路径
	// r.LoadHTMLGlob("web/templates/*") //部署路径

	// 为 multipart forms 设置较低的内存限制
	r.MaxMultipartMemory = 3
	//路由
	r.GET("/home", api.HomeHandler)
	r.GET("/help", api.HelpHandler)
	r.POST("/upload", api.UploadHandler)
	return r
}
