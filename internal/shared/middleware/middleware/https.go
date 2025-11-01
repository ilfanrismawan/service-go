package middleware

import (
	"net/http"
	"service/internal/config"

	"github.com/gin-gonic/gin"
)

// HTTPSRedirectMiddleware enforces HTTPS in production by redirecting HTTP requests.
// It respects X-Forwarded-Proto when behind a reverse proxy/ingress.
func HTTPSRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Config != nil && config.Config.Environment == "production" {
			proto := c.GetHeader("X-Forwarded-Proto")
			if proto == "" {
				if c.Request.TLS == nil {
					redirectToHTTPS(c)
					return
				}
			} else if proto != "https" {
				redirectToHTTPS(c)
				return
			}
		}
		c.Next()
	}
}

func redirectToHTTPS(c *gin.Context) {
	url := *c.Request.URL
	url.Scheme = "https"
	url.Host = c.Request.Host
	c.Redirect(http.StatusMovedPermanently, url.String())
	c.Abort()
}


