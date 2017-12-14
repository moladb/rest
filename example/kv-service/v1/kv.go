package v1

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
	return nil
}

func (s *KVService) GetAPIGroup() string {
	return "/v1"
}

func (s *KVService) ListHandlers() []rest.Handler {
	return []rest.Handler{
		{
			Resource: rest.Resource{
				Name:   "kv",
				Method: "GET",
				Path:   "/kv/*key",
			},
			HandlerFunc: getKVsHandler(s),
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
			HandlerFunc: deleteKVsHandler(s),
		},
	}
}

type KV struct {
	Key   string
	Value string
}

func (s *KVService) getByKey(key string, keyOnly bool) ([]KV, bool) {
	if keyOnly {
		return []KV{{Key: "this-is-key"}}, true
	}
	return []KV{{Key: "this-is-key", Value: "this-is-value"}}, true
}

func (s *KVService) getByPrefix(prefix string, keyOnly bool) []KV {
	return []KV{{Key: "this-is-key-1"}, {Key: "this-is-key-2"}}
}

func (s *KVService) putKV(key, value string) {
}

func (s *KVService) deleteByKey(key string) {
}

func (s *KVService) deleteByPrefix(prefix string) {
}

func getKVsHandler(s *KVService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			key     string
			onlyKey bool
			prefix  bool
		)
		key = c.Param("key")
		_, onlyKey = c.GetQuery("keys")
		_, prefix = c.GetQuery("prefix")
		if prefix {
			kvs := s.getByPrefix(key, onlyKey)
			c.JSON(http.StatusOK, kvs)
		} else {
			kvs, ok := s.getByKey(key, onlyKey)
			if ok {
				c.JSON(http.StatusOK, kvs)
			} else {
				c.Status(http.StatusNotFound)
			}
		}
	}
}

func putKVHandler(s *KVService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			key string
			val struct {
				Value string `json:"value" binding:"required"`
			}
		)
		key = c.Param("key")
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

func deleteKVsHandler(s *KVService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			key    string
			prefix bool
		)
		key = c.Param("key")
		_, prefix = c.GetQuery("prefix")
		if prefix {
			s.deleteByPrefix(key)
		} else {
			s.deleteByKey(key)
		}
		c.Status(http.StatusOK)
	}
}
