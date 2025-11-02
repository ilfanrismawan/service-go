package middleware

import (
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// InitSentry initializes Sentry if DSN is provided.
func InitSentry(dsn string, environment string) error {
	if dsn == "" {
		return nil
	}
	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		TracesSampleRate: 0.1,
	})
}

// SentryMiddleware captures panics and request errors to Sentry.
func SentryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clone hub per request
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetTag("path", c.FullPath())
		if rid, exists := c.Get("request_id"); exists {
			hub.Scope().SetTag("request_id", rid.(string))
		}

		defer func() {
			if rec := recover(); rec != nil {
				hub.Recover(rec)
				hub.Flush(2 * time.Second)
				panic(rec)
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if e.Err != nil {
					hub.CaptureException(e.Err)
				}
			}
			hub.Flush(2 * time.Second)
		}
	}
}
