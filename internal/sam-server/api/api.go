package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"student-attendance/internal/pkg/model"
	"student-attendance/internal/sam-server/service/check"
	"student-attendance/internal/sam-server/service/controller"
	"student-attendance/internal/sam-server/store/mysql"

	"github.com/gin-gonic/gin"
)

// HomeHandler 提供用户访问入口
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

// UploadHandler 处理用户提交的考勤数据
func UploadHandler(c *gin.Context) {
	ui := new(model.UploadInfo)
	if err := c.ShouldBind(ui); err != nil {
		//这里将来可能要加错误处理逻辑
		fmt.Println(err)
	}
	//校验提交的信息是否齐全与合理
	if err := check.CheckUploadInfo(ui); err != nil {
		fmt.Println(err)
	}

	//写校验User和Subject是否相符的逻辑
	_, err := mysql.GetTeacherByNameNSubject(ui.User, ui.Subject)
	if err != nil {
		fmt.Println(err)
	}

	//写校验上传文件是否相符的逻辑
	file, err := c.FormFile("excel-file")
	if err != nil {
		fmt.Println(err) //没有上传文件的逻辑
	}
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		fmt.Println("上传的文件不合规格的逻辑")
	}

	//保存上传的考勤文件
	dst := "./excelfiles" + file.Filename
	log.Println(file.Filename)
	c.SaveUploadedFile(file, dst)

	//执行业务逻辑
	res, err := controller.Process(dst, ui)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, res)
}
