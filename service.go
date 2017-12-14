package rest

import (
	"fmt"
	"strings"

	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	Resource
	HandlerFunc gin.HandlerFunc
}

type Service interface {
	ListHandlers() []Handler
}

type ServiceGroup interface {
	Service
	GetAPIGroup() string
}

//func validateService(svc Service) error {
//	// TODO
//	return nil
//}
//
//func validateServiceGroup(svc ServiceGroup) error {
//	// TODO
//	return nil
//}

type Resource struct {
	Name   string
	Path   string
	Method string
}

type APIGroup struct {
	Name      string
	Resources []Resource
}

type serviceRegistry struct {
	apiGroups map[string]APIGroup
}

func newServiceRegistry() *serviceRegistry {
	return &serviceRegistry{
		apiGroups: make(map[string]APIGroup),
	}
}

func (r *serviceRegistry) AddGroupResource(group string, res Resource) {
	group = strings.Trim(group, "/")
	r.addGroupResource(group, res)
}

func (r *serviceRegistry) AddResource(res Resource) {
	r.addGroupResource("/", res)
}

func (r *serviceRegistry) addGroupResource(group string, res Resource) {
	apiGroup, ok := r.apiGroups[group]
	if !ok {
		apiGroup = APIGroup{Name: group}
	}
	apiGroup.Resources = append(apiGroup.Resources, res)
	r.apiGroups[group] = apiGroup
}

func (r *serviceRegistry) ListAPIGroups() []string {
	names := []string{}
	for k := range r.apiGroups {
		names = append(names, k)
	}
	return names
}

func (r *serviceRegistry) ListGroupResources(apiGroup string) (APIGroup, bool) {
	apiGroup = strings.Trim(apiGroup, "/")
	g, ok := r.apiGroups[apiGroup]
	return g, ok
}

func (r *serviceRegistry) ListResources() []Resource {
	g, ok := r.apiGroups["/"]
	if !ok {
		return []Resource{}
	}
	return g.Resources
}

type discoveryService struct {
	registry *serviceRegistry
}

func newDiscoveryService(r *serviceRegistry) *discoveryService {
	return &discoveryService{registry: r}
}

func (s *discoveryService) ListHandlers() []Handler {
	return []Handler{
		{
			Resource: Resource{
				Name:   "/apis",
				Path:   "/apis",
				Method: "GET",
			},
			HandlerFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK,
					gin.H{
						"apis": s.registry.ListAPIGroups(),
					})
			},
		},
		{
			Resource: Resource{
				Name:   "/apis/",
				Path:   "/apis/",
				Method: "GET",
			},
			HandlerFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, s.registry.ListResources())
			},
		},
		{
			Resource: Resource{
				Name:   "/apis/apigroup",
				Path:   "/apis/:apigroup",
				Method: "GET",
			},
			HandlerFunc: func(c *gin.Context) {
				apiGroup := c.Param("apigroup")
				rs, ok := s.registry.ListGroupResources(apiGroup)
				if !ok {
					c.JSON(http.StatusNotFound,
						gin.H{
							"error": fmt.Sprintf("APIGroup:%s not found", apiGroup),
						})
					return
				}
				c.JSON(http.StatusOK, rs)
			},
		},
	}
}

type pprofService struct{}

func newPProfService() *pprofService {
	return &pprofService{}
}

func (s *pprofService) GetAPIGroup() string {
	return "debug"
}

func (s *pprofService) ListHandlers() []Handler {
	return []Handler{
		{
			Resource: Resource{
				Name:   "pprof",
				Path:   "/pprof",
				Method: "GET",
			},
			HandlerFunc: ginHandlerFunc(pprof.Index),
		},
		{
			Resource: Resource{
				Name:   "/pprof/profile",
				Path:   "/pprof/profile",
				Method: "GET",
			},
			HandlerFunc: ginHandlerFunc(pprof.Profile),
		},
		{
			Resource: Resource{
				Name:   "/pprof/cmdline",
				Path:   "/pprof/cmdline",
				Method: "GET",
			},
			HandlerFunc: ginHandlerFunc(pprof.Cmdline),
		},
	}
}

type metricsService struct{}

func newMetricsService() *metricsService {
	return &metricsService{}
}

func (s *metricsService) ListHandlers() []Handler {
	h := promhttp.Handler()
	return []Handler{
		{
			Resource: Resource{
				Name:   "/metrics",
				Path:   "/metrics",
				Method: "GET",
			},
			HandlerFunc: func(c *gin.Context) {
				h.ServeHTTP(c.Writer, c.Request)
			},
		},
	}
}
