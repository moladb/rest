package v0

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/moladb/rest"
)

const maxDataLen int = 512 * 1024

type KVService struct {
	kvs    map[string]string
	kvLock sync.RWMutex
}

func NewKVService() *KVService {
	return &KVService{
		kvs: make(map[string]string),
	}
}

func (s *KVService) GetAPIGroup() string {
	return "/v0"
}

func (s *KVService) ListHandlers() []rest.Handler {
	return []rest.Handler{
		{
			Resource: rest.Resource{
				Name:   "kv",
				Method: "GET",
				Path:   "/kv/*key",
			},
			HandlerFunc: getKVHandler(s),
		},
		{
			Resource: rest.Resource{
				Name:   "kv",
				Method: "PUT",
				Path:   "/kv/*key",
			},
			HandlerFunc: putKVHandler(s),
		},
		{
			Resource: rest.Resource{
				Name:   "kv",
				Method: "DELETE",
				Path:   "/kv/*key",
			},
			HandlerFunc: deleteKVHandler(s),
		},
	}
}

func (s *KVService) getKV(key string) (string, bool) {
	s.kvLock.RLock()
	defer s.kvLock.RUnlock()
	val, ok := s.kvs[key]
	return val, ok
}

func (s *KVService) putKV(key, val string) {
	s.kvLock.Lock()
	defer s.kvLock.Unlock()
	s.kvs[key] = val
}

func (s *KVService) deleteKV(key string) {
	s.kvLock.Lock()
	defer s.kvLock.Unlock()
	delete(s.kvs, key)
}

func getKVHandler(s *KVService) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		val, ok := s.getKV(key)
		if ok {
			c.JSON(http.StatusOK, gin.H{"value": val})
			return
		}
		c.Status(http.StatusNotFound)
	}
}

func putKVHandler(s *KVService) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		var val struct {
			Value string `json:"value" binding:"required"`
		}
		if err := c.ShouldBindWith(&val, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(val.Value) > maxDataLen {
			c.JSON(http.StatusBadRequest, gin.H{"error": "exceed max_data_len(512K)"})
			return
		}
		s.putKV(key, val.Value)
		c.Status(http.StatusOK)
	}
}

func deleteKVHandler(s *KVService) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		s.deleteKV(key)
		c.Status(http.StatusOK)
	}
}
