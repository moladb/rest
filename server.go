package rest

import (
	"context"
	"log"
	"net/http"
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
			return ginprom.WithMetrics(fullResourcePath(apiGroup, resource), h)
		}
	} else {
		s.decorateHandler = func(_, _ string, h gin.HandlerFunc) gin.HandlerFunc {
			return h
		}
	}
	return s
}

func (s *Server) RegisterService(svc Service) {
	s.RegisterServiceGroup(DefaultGroup, svc)
}

func (s *Server) RegisterServiceGroup(group string, svc Service) {
	group = normalizePath(group)
	var r gin.IRoutes = s.router
	if group != DefaultGroup {
		r = s.router.Group(group)
	}
	handlers := svc.ListHandlers()
	for _, h := range handlers {
		r.Handle(h.Method, normalizePath(h.Path),
			s.decorateHandler(group, h.Resource.Name, h.HandlerFunc))
		s.registry.addGroupResource(group, h.Resource)
	}
}

func (s *Server) Run() error {
	if s.config.EnableDebug {
		s.RegisterServiceGroup("/debug", newPProfService())
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
