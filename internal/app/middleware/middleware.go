package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// в виду удобства были использованы встроенные инструменты из Gin,
// но возможность создания кастомного Writer была также рассмотрена
// type (
// 	responseData struct {
// 		status int
// 		size   int
// 	}

// 	loggingResponseWriter struct {
// 		http.ResponseWriter
// 		responseData *responseData
// 	}
// )

// func (r *loggingResponseWriter) Write(b []byte) (int, error) {
// 	size, err := r.ResponseWriter.Write(b)
// 	r.responseData.size += size
// 	return size, err
// }

// func (r *loggingResponseWriter) WriteHeader(statusCode int) {
// 	r.ResponseWriter.WriteHeader(statusCode)
// 	r.responseData.status = statusCode
// }

func WithLogging(logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()
		duration := time.Since(start)

		logger.Info(
			"request details",
			zap.String("uri", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("size", c.Writer.Size()),
			zap.String("location", c.Writer.Header().Get("Location")),
		)
	}
}
