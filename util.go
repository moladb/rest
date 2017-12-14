package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ginHandlerFunc(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
