package web

import "github.com/gin-gonic/gin"

//设置全局配置的中间件
func ConfigMiddleware(m map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range m {
			c.Set(k, v)
			c.Next()
		}
	}
}
