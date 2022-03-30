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
	"time"

	"github.com/gin-gonic/gin"
)

const (
	incompleteInfo string = "您上传的信息不全,请补全信息后重新上传"
	illogicalTime  string = "您上传的课堂时间不合理,可能是课堂时间过短或下课时间早于上课时间,请选择合理的时间重新上传"
	userNotFound   string = "抱歉,没有找到您的用户记录,请填写您的真实姓名与任教学科"
	incorrectFile  string = "无法读取您上传的考勤记录,请上传未经修改的腾讯课堂原始考勤文件(以.xlsx结尾的Excel文件)"
)

// HomeHandler 提供用户访问入口
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

func HelpHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "help.html", nil)
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
		switch err {
		case check.ErrIncompleteInfo:
			c.HTML(http.StatusBadRequest, "error.html", incompleteInfo)
			return
		case check.ErrIllogicalTime:
			c.HTML(http.StatusBadRequest, "error.html", illogicalTime)
			return
		}
	}

	//写校验User和Subject是否相符的逻辑
	_, err := mysql.GetTeacherByNameNSubject(ui.User, ui.Subject)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", userNotFound)
		return
	}

	//写校验上传文件是否相符的逻辑
	file, err := c.FormFile("excel-file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", incorrectFile)
		return
	}
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		c.HTML(http.StatusBadRequest, "error.html", incorrectFile)
		return
	}

	//保存上传的考勤文件
	timeStr := time.Now().Format("06-01-02_15_04_")
	dst := strings.Join([]string{"../excelfiles/", ui.User, timeStr, file.Filename}, "")
	log.Println(file.Filename)
	c.SaveUploadedFile(file, dst)

	//执行业务逻辑
	res, err := controller.Process(dst, ui)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", incorrectFile)
		return
	}
	c.HTML(http.StatusOK, "result.html", res)

}
