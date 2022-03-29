package router

import (
	"student-attendance/internal/sam-server/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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
	r.POST("/upload", api.UploadHandler)
	return r
}
