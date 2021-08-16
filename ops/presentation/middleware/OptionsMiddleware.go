package middleware

import (
	"github.com/gin-gonic/gin"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 4:28 PM
 * @Desc:
 */

func OptionsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
		"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, "+
		"Token, Language, From, Cookie, OperationCode, ClusterId")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}