package middleware

import (
	"github.com/gin-gonic/gin"
)

type corsMiddleware struct {
	// middleware class , may be needed by middleware
}

func CorsMiddleware() gin.HandlerFunc {
	cm := &corsMiddleware{
		//to be used
	}

	return cm.handle
}

func (cm *corsMiddleware) handle(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Next()
}
