package middleware

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 2:38 PM
 * @Desc:
 */

func KubernetesMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiServer, token := models.GetK8sInfoByClusterId(c.GetHeader("ClusterId"))
		if k8sClient, err := service.NewKubernetesService(apiServer, token); err == nil {
			c.Set("k8sCluster", &k8sClient)
		}
		c.Next()
	}
}
