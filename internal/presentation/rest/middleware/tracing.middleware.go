package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}).Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.EscapedPath()),
		)

		c.Request = c.Request.WithContext(ctx)

		defer func() {
			status := c.Writer.Status()
			span.SetAttributes(attribute.Int("http.status_code", status))

			if len(c.Errors) > 0 {
				err := c.Errors.Last().Err
				if status >= 500 && status <= 599 {
					span.SetStatus(codes.Error, err.Error())
				}

				if status >= 400 && status <= 599 {
					span.RecordError(err)
				}
			}
		}()

		c.Next()
	}
}
