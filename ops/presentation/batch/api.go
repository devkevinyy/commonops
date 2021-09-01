package batch

import (
	"net/http"
	"strconv"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

// IGetAllNodeTree get servers by user role id for user.
func IGetAllNodeTree(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	data, err := models.GetAllNodeTree(uint(userId))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}
