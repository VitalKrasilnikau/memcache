package controllers

import (
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/gin-gonic/gin"
)

// GetCacheKeysHandler API which gets all string cache keys.
func GetCacheKeysHandler(pid *act.BroadcastStringKeysGroup) func(*gin.Context) {
	return func(c *gin.Context) {
		res, e := pid.Request()
		if e == nil {
			api.OK(c, contracts.CacheKeysContract{Keys: res.Keys})
		} else {
			api.Error(c, e.Error())
		}
	}
}
