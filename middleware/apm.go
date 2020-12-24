package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pi-sin/go-repo-structure/config"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

type apmMiddleware struct {
	tracer         *apm.Tracer
	enableApm      bool
	requestIgnorer apmhttp.RequestIgnorerFunc
}

func ApmMiddleware() gin.HandlerFunc {
	am := &apmMiddleware{
		tracer:         apm.DefaultTracer,
		enableApm:      config.GetConfig().GetBool("apm_enabled"),
		requestIgnorer: apmhttp.DefaultServerRequestIgnorer(),
	}

	return am.handle
}

func (m *apmMiddleware) handle(c *gin.Context) {
	if !m.enableApm || m.requestIgnorer(c.Request) {
		c.Next()
		return
	}

	/**
	Code for APM Handling
	*/
	c.Next()
}
