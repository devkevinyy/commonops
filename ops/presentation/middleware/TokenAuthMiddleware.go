package middleware

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/exception"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 2:34 PM
 * @Desc:
 */

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId int32
		var authValid bool
		authToken := getUserToken(c)
		parseToken, err := jwt.Parse(authToken, func(*jwt.Token) (interface{}, error) {
			return []byte(conf.SecretSalt), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespData{Code: 0, Msg: exception.TokenInvalidException.Error(), Data: nil})
		}
		if parseToken != nil {
			var parmMap = parseToken.Claims.(jwt.MapClaims)
			var expire = (int64)(parmMap["exp"].(float64))
			var now = time.Now().Unix()
			if expire < now {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespData{Code: 0, Msg: exception.TokenInvalidException.Error(), Data: nil})
			}
			var userInfo = parmMap["userInfo"].(map[string]interface{})
			userId = int32(userInfo["userId"].(float64))
			if userValid := models.IsUserValid(userId); !userValid {
				c.AbortWithStatusJSON(http.StatusForbidden, utils.RespData{Code: 0, Msg: exception.UserInvalidException.Error(), Data: nil})
			}
			for k, v := range userInfo {
				c.Set(k, v)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespData{Code: 0, Msg: exception.TokenInvalidException.Error(), Data: nil})
		}

		userAction := fmt.Sprintf("%s:%s", c.Request.Method, utils.ExtractUriPath(c.Request.RequestURI))
		if authValid = models.IsUserHasActionPermision(userId, userAction); !authValid {
			msg := fmt.Sprintf("用户操作无权限: %s", userAction)
			c.AbortWithStatusJSON(http.StatusAccepted, utils.RespData{Code: -1, Msg: msg, Data: nil})
		}
		fmt.Printf("用户权限校验: %s, %v \n", userAction, authValid)
		c.Next()
	}
}

func getUserToken(c *gin.Context) string {
	var authToken = c.GetHeader("Authorization")
	if strings.TrimSpace(authToken) != "" {
		return authToken
	}
	return c.Query("token")
}