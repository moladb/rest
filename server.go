package rest

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moladb/ginprom"
)

type Server struct {
	config          Config
	router          *gin.Engine
	httpSrv         *http.Server
	decorateHandler func(apiGroup, api string, h gin.HandlerFunc) gin.HandlerFunc
	registry        *serviceRegistry
}

func NewServer(config Config) *Server {
	s := &Server{
		config:   config,
		router:   gin.Default(),
		registry: newServiceRegistry(),
	}
	if config.EnableMetrics {
		s.decorateHandler = func(apiGroup, resource string, h gin.HandlerFunc) gin.HandlerFunc {
			var path string
			if apiGroup == "" {
				path = "/" + strings.TrimLeft(resource, "/")
			} else {
				path = "/" + strings.Trim(apiGroup, "/") + "/" + strings.TrimLeft(resource, "/")
			}
			return ginprom.WithMetrics(path, h)
		}
	} else {
		s.decorateHandler = func(_, _ string, h gin.HandlerFunc) gin.HandlerFunc {
			return h
		}
	}
	return s
}

func (s *Server) RegisterServiceGroup(svc ServiceGroup) {
	apiGroup := strings.Trim(svc.GetAPIGroup(), "/")
	group := s.router.Group("/" + apiGroup)
	handlers := svc.ListHandlers()
	for _, h := range handlers {
		group.Handle(h.Method, "/"+strings.TrimLeft(h.Path, "/"),
			s.decorateHandler(apiGroup, h.Resource.Name, h.HandlerFunc))
		s.registry.AddGroupResource(apiGroup, h.Resource)
	}
}

func (s *Server) RegisterService(svc Service) {
	handlers := svc.ListHandlers()
	for _, h := range handlers {
		s.router.Handle(h.Method, "/"+strings.TrimLeft(h.Path, "/"),
			s.decorateHandler("", h.Resource.Name, h.HandlerFunc))
		s.registry.AddResource(h.Resource)
	}
}

func (s *Server) Run() error {
	if s.config.EnableDebug {
		s.RegisterServiceGroup(newPProfService())
	}
	if s.config.EnableMetrics {
		s.RegisterService(newMetricsService())
	}
	if s.config.EnableDiscovery {
		s.RegisterService(newDiscoveryService(s.registry))
	}
	s.httpSrv = &http.Server{
		Addr:    s.config.BindAddr,
		Handler: s.router,
	}
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(s.config.GraceShutdownTimeoutS)*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server GraceShutdown:", err)
	}
}
