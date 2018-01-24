package v1

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moladb/rest"
)

type EchoService struct {
	m         sync.Mutex
	health    bool
	lastCheck time.Time
}

func NewEchoService() *EchoService {
	return &EchoService{
		health:    true,
		lastCheck: time.Now(),
	}
}

func (s *EchoService) GetAPIGroup() string {
	return "/v1"
}

func (s *EchoService) ListHandlers() []rest.Handler {
	return []rest.Handler{
		{
			Resource: rest.Resource{
				Name:   "echo",
				Method: "GET",
				Path:   "/echo/*msg",
			},
			HandlerFunc: func(c *gin.Context) {
				msg := c.Param("msg")
				c.String(http.StatusOK, msg)
			},
		},
		{
			Resource: rest.Resource{
				Name:   "health",
				Method: "GET",
				Path:   "/health",
			},
			HandlerFunc: func(c *gin.Context) {
				// alive 120s dead 10s
				var (
					code int    = http.StatusOK
					msg  string = "OK"
				)
				s.m.Lock()
				defer s.m.Unlock()
				// change to unhealth
				if s.health && time.Now().After(s.lastCheck.Add(120*time.Second)) {
					s.health = false
					s.lastCheck = time.Now()
					code = http.StatusOK
					msg = "OK"
				} else if !s.health && time.Now().After(s.lastCheck.Add(30*time.Second)) {
					s.health = true
					s.lastCheck = time.Now()
					code = http.StatusServiceUnavailable
					msg = "RetryLater"
				}
				c.String(code, msg)
			},
		},
	}
}
