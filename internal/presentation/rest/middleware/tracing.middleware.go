package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func TracingMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))

		span := tracer.StartSpan(
			fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			ext.RPCServerOption(spanCtx),
		)

		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.EscapedPath())

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		defer func() {
			status := c.Writer.Status()
			ext.HTTPStatusCode.Set(span, uint16(status))

			if status >= 500 && status <= 599 {
				ext.Error.Set(span, true)
			}

			if status >= 400 && status <= 599 && len(c.Errors) > 0 {
				err := c.Errors.Last().Err
				span.LogKV(
					"event", "error",
					"message", err.Error(),
					"error.kind", fmt.Sprintf("%T", err),
				)
			}

			span.Finish()
		}()

		c.Next()
	}
}
