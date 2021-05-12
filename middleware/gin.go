package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gotrace "github.com/lyouthzzz/go-trace"
)

func ServerTracingInterceptor() gin.HandlerFunc {
	propagator := gotrace.GetPropagator()
	tracer := gotrace.GetGlobalTracer()

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		carrier := gotrace.HTTPCarrier(c.Request.Header)
		
		ctx = propagator.Extract(ctx, carrier)
		ctx, span := tracer.StartSpan(ctx, c.Request.Method)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if c.Writer.Status() != http.StatusOK {
			if c.Errors != nil && len(c.Errors) != 0 {
				span.AddError(c.Errors[0])
			}
		}
	}
}
