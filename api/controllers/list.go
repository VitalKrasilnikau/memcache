package controllers

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/gin-gonic/gin"
)

// GetListCacheKeyHandler API which gets list cache entry by key.
func GetListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		res, e := pid.RequestFuture(&act.GetListCacheKeyMessage{Key: key}, defTimeout).Result()
		if e == nil {
			s, ok := res.(act.GetListCacheKeyReply)
			if ok {
				if s.Success {
					api.OK(c, contracts.ListCacheValueContract{Key: s.Key, Values: s.Values})
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

// DeleteListCacheKeyHandler API which deletes list cache entry by key.
func DeleteListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		res, e := pid.RequestFuture(&act.DeleteListCacheKeyMessage{Key: key}, defTimeout).Result()
		if e == nil {
			s, ok := res.(act.DeleteListCacheKeyReply)
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

// PostListCacheKeyHandler API which posts new list value.
func PostListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		var json contracts.NewListCacheValuesContract
		if err := c.ShouldBindJSON(&json); err == nil {
			message := &act.PostListCacheKeyMessage{Key: json.Key, Values: json.Values, TTL: api.ParseDuration(json.TTL)}
			res, e := pid.RequestFuture(message, defTimeout).Result()
			if e == nil {
				s, ok := res.(act.PostListCacheKeyReply)
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

// PutListCacheValueHandler API which updates existing string value in the list by the key and new value specified in body.
func PutListCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")
		var json contracts.UpdateListCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			entry := &act.PutListCacheValueMessage{Key: key, NewValue: json.Value, OriginalValue: value}
			res, e := pid.RequestFuture(entry, defTimeout).Result()
			if e == nil {
				s, ok := res.(act.PutListCacheValueReply)
				if ok {
					if s.Success {
						api.NoContent(c)
					} else {
						api.Bad(c, fmt.Sprintf("list value '%s' of key '%s' was already changed or never existed", value, key))
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

// PostListCacheValueHandler API which updates existing list value by the key and new added value specified in body.
func PostListCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		var json contracts.UpdateListCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			entry := &act.PostListCacheValueMessage{Key: key, NewValue: json.Value}
			res, e := pid.RequestFuture(entry, defTimeout).Result()
			if e == nil {
				s, ok := res.(act.PostListCacheValueReply)
				if ok {
					if s.Success {
						api.Created(c)
					} else {
						api.Bad(c, fmt.Sprintf("key '%s' was not found", key))
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

// DeleteListCacheValueHandler API which updates existing list value by the key and deleting the value specified.
func DeleteListCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")
		entry := &act.DeleteListCacheValueMessage{Key: key, Value: value}
		res, e := pid.RequestFuture(entry, defTimeout).Result()
		if e == nil {
			s, ok := res.(act.DeleteListCacheValueReply)
			if ok {
				if s.Success {
					api.NoContent(c)
				} else {
					api.Bad(c, fmt.Sprintf("list value '%s' of key '%s' was already deleted or never existed", value, key))
				}
			} else {
				api.Error(c, "can't parse data")
			}
		} else {
			api.Error(c, e.Error())
		}
	}
}
