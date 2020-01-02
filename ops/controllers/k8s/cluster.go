package k8s

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/forms"
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
)

func IPostKubeConfigFileUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "上传失败", Data: nil})
		return
	}

	path, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/conf/%s", path, untils.GetUUID())
	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = out.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "上传成功", Data: filePath})
}

func PostK8sCluster(c *gin.Context) {
	var req forms.K8sForm
	err := c.Bind(&req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	err = models.AddK8sCluster(req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK,  untils.RespData{Code:0, Msg:"success", Data: nil})
}

func DeleteK8sCluster(c *gin.Context) {
	var req forms.K8sForm
	err := c.Bind(&req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	err = models.DeleteK8s(req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK,  untils.RespData{Code:0, Msg:"success", Data: nil})
}

func GetK8sCluster(c *gin.Context) {
	var req forms.K8sForm
	err := c.Bind(&req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	totalCount := models.GetK8sCount()
	data := models.GetK8sList()
	c.JSON(http.StatusOK,  untils.RespData{Code:0, Msg:"success", Data: map[string]interface{}{
		"totalCount": totalCount,
		"k8sData": data,
	}})
}
