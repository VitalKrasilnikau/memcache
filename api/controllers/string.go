package controllers

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/gin-gonic/gin"
)

// GetStringCacheKeyHandler API which gets string cache entry by key.
func GetStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		res, e := pid.RequestFuture(&act.GetStringCacheKeyMessage{Key: key}, defTimeout).Result()
		if e == nil {
			s, ok := res.(act.GetStringCacheKeyReply)
			if ok {
				if s.Success {
					api.OK(c, contracts.StringCacheValueContract{Key: s.Key, Value: s.Value})
				} else {
					api.NotFound(c, fmt.Sprintf("key '%s' was not found", key))
				}
			} else {
				api.Error(c, "can't parse data")
			}
		} else {
			api.Error(c, e.Error())
		}
	}
}

// DeleteStringCacheKeyHandler API which deletes string cache entry by key.
func DeleteStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		res, e := pid.RequestFuture(&act.DeleteStringCacheKeyMessage{Key: key}, defTimeout).Result()
		if e == nil {
			s, ok := res.(act.DeleteStringCacheKeyReply)
			if ok {
				if s.Success {
					api.NoContent(c)
				} else {
					api.NotFound(c, fmt.Sprintf("key '%s' was not found", key))
				}
			} else {
				api.Error(c, "can't parse data")
			}
		} else {
			api.Error(c, e.Error())
		}
	}
}

// PostStringCacheKeyHandler API which posts new string value.
func PostStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		var json contracts.NewStringCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			message := &act.PostStringCacheKeyMessage{Key: json.Key, Value: json.Value, TTL: api.ParseDuration(json.TTL)}
			res, e := pid.RequestFuture(message, defTimeout).Result()
			if e == nil {
				s, ok := res.(act.PostStringCacheKeyReply)
				if ok {
					if s.Success {
						api.Created(c)
					} else {
						api.Bad(c, fmt.Sprintf("key '%s' was already used", json.Key))
					}
				} else {
					api.Error(c, "can't parse data")
				}
			} else {
				api.Error(c, e.Error())
			}
		} else {
			api.Bad(c, fmt.Sprintf("malformed request: %s", err.Error()))
		}
	}
}

// PutStringCacheKeyHandler API which updates existing string value by the key and new value specified.
func PutStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		var json contracts.UpdateStringCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			entry := &act.PutStringCacheKeyMessage{Key: key, NewValue: json.NewValue, OriginalValue: json.OriginalValue}
			res, e := pid.RequestFuture(entry, defTimeout).Result()
			if e == nil {
				s, ok := res.(act.PutStringCacheKeyReply)
				if ok {
					if s.Success {
						api.NoContent(c)
					} else {
						api.Bad(c, fmt.Sprintf("key '%s' was already changed to '%s'", key, s.OriginalValue))
					}
				} else {
					api.Error(c, "can't parse data")
				}
			} else {
				api.Error(c, e.Error())
			}
		} else {
			api.Bad(c, fmt.Sprintf("malformed request: %s", err.Error()))
		}
	}
}
