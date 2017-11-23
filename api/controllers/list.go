package controllers

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/gin-gonic/gin"
	"sync"
)

// GetListCacheKeyHandler API which gets list cache entry by key.
func GetListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		act.Await(
			pid,
			func(wg *sync.WaitGroup) *actor.PID {
				return createListReplyActor(c, wg)
			},
			func(replyPid *actor.PID) interface{} {
				return &act.GetListCacheKeyMessage{Key: key, ReplyTo: replyPid}
			})
	}
}

// DeleteListCacheKeyHandler API which deletes list cache entry by key.
func DeleteListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		act.Await(
			pid,
			func(wg *sync.WaitGroup) *actor.PID {
				return createListReplyActor(c, wg)
			},
			func(replyPid *actor.PID) interface{} {
				return &act.DeleteListCacheKeyMessage{Key: key, ReplyTo: replyPid}
			})
	}
}

// PostListCacheKeyHandler API which posts new list value.
func PostListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		var json contracts.NewListCacheValuesContract
		if err := c.ShouldBindJSON(&json); err == nil {
			act.Await(
				pid,
				func(wg *sync.WaitGroup) *actor.PID {
					return createListReplyActor(c, wg)
				},
				func(replyPid *actor.PID) interface{} {
					return &act.PostListCacheKeyMessage{
						Key:     json.Key,
						Values:  json.Values,
						TTL:     api.ParseDuration(json.TTL),
						ReplyTo: replyPid}
				})
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
			act.Await(
				pid,
				func(wg *sync.WaitGroup) *actor.PID {
					return createListReplyActor(c, wg)
				},
				func(replyPid *actor.PID) interface{} {
					return &act.PutListCacheValueMessage{Key: key, NewValue: json.Value, OriginalValue: value, ReplyTo: replyPid}
				})
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
			act.Await(
				pid,
				func(wg *sync.WaitGroup) *actor.PID {
					return createListReplyActor(c, wg)
				},
				func(replyPid *actor.PID) interface{} {
					return &act.PostListCacheValueMessage{Key: key, NewValue: json.Value, ReplyTo: replyPid}
				})
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
		act.Await(
			pid,
			func(wg *sync.WaitGroup) *actor.PID {
				return createListReplyActor(c, wg)
			},
			func(replyPid *actor.PID) interface{} {
				return &act.DeleteListCacheValueMessage{Key: key, Value: value, ReplyTo: replyPid}
			})
	}
}

func dispatchListReply(c *gin.Context, ctx actor.Context, wg *sync.WaitGroup) {
	switch s := ctx.Message().(type) {
	case act.GetListCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.OK(c, contracts.ListCacheValueContract{Key: s.Key, Values: s.Values})
		} else {
			api.NotFound(c, fmt.Sprintf("key '%s' was not found", s.Key))
		}
		break
	case act.DeleteListCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.NoContent(c)
		} else {
			api.NotFound(c, fmt.Sprintf("key '%s' was not found", s.Key))
		}
		break
	case act.PostListCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.Created(c)
		} else {
			api.Bad(c, fmt.Sprintf("key '%s' was already used", s.Key))
		}
		break
	case act.PutListCacheValueReply:
		defer wg.Done()
		if s.Success {
			api.NoContent(c)
		} else {
			api.Bad(c, fmt.Sprintf("list value '%s' of key '%s' was already changed or never existed", s.OriginalValue, s.Key))
		}
		break
	case act.PostListCacheValueReply:
		defer wg.Done()
		if s.Success {
			api.Created(c)
		} else {
			api.Bad(c, fmt.Sprintf("key '%s' was not found", s.Key))
		}
		break
	case act.DeleteListCacheValueReply:
		defer wg.Done()
		if s.Success {
			api.NoContent(c)
		} else {
			api.Bad(c, fmt.Sprintf("list value '%s' of key '%s' was already deleted or never existed", s.DeletedValue, s.Key))
		}
		break
	}
}

func createListReplyActor(c *gin.Context, wg *sync.WaitGroup) *actor.PID {
	return actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
		dispatchListReply(c, ctx, wg)
	}))
}
