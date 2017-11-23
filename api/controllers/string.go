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

// GetStringCacheKeyHandler API which gets string cache entry by key.
func GetStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		act.Await(
			pid,
			func(wg *sync.WaitGroup) *actor.PID {
				return createStringReplyActor(c, wg)
			},
			func(replyPid *actor.PID) interface{} {
				return &act.GetStringCacheKeyMessage{Key: key, ReplyTo: replyPid}
			})
	}
}

// DeleteStringCacheKeyHandler API which deletes string cache entry by key.
func DeleteStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		act.Await(
			pid,
			func(wg *sync.WaitGroup) *actor.PID {
				return createStringReplyActor(c, wg)
			},
			func(replyPid *actor.PID) interface{} {
				return &act.DeleteStringCacheKeyMessage{Key: key, ReplyTo: replyPid}
			})
	}
}

// PostStringCacheKeyHandler API which posts new string value.
func PostStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		var json contracts.NewStringCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			act.Await(
				pid,
				func(wg *sync.WaitGroup) *actor.PID {
					return createStringReplyActor(c, wg)
				},
				func(replyPid *actor.PID) interface{} {
					return &act.PostStringCacheKeyMessage{
						Key:     json.Key,
						Value:   json.Value,
						TTL:     api.ParseDuration(json.TTL),
						ReplyTo: replyPid}
				})
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
			act.Await(
				pid,
				func(wg *sync.WaitGroup) *actor.PID {
					return createStringReplyActor(c, wg)
				},
				func(replyPid *actor.PID) interface{} {
					return &act.PutStringCacheKeyMessage{
						Key:           key,
						NewValue:      json.NewValue,
						OriginalValue: json.OriginalValue,
						ReplyTo:       replyPid}
				})
		} else {
			api.Bad(c, fmt.Sprintf("malformed request: %s", err.Error()))
		}
	}
}

func dispatchStringReply(c *gin.Context, ctx actor.Context, wg *sync.WaitGroup) {
	switch s := ctx.Message().(type) {
	case act.GetStringCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.OK(c, contracts.StringCacheValueContract{Key: s.Key, Value: s.Value})
		} else {
			api.NotFound(c, fmt.Sprintf("key '%s' was not found", s.Key))
		}
		break
	case act.DeleteStringCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.NoContent(c)
		} else {
			api.NotFound(c, fmt.Sprintf("key '%s' was not found", s.Key))
		}
		break
	case act.PostStringCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.Created(c)
		} else {
			api.Bad(c, fmt.Sprintf("key '%s' was already used", s.Key))
		}
		break
	case act.PutStringCacheKeyReply:
		defer wg.Done()
		if s.Success {
			api.NoContent(c)
		} else {
			api.Bad(c, fmt.Sprintf("key '%s' was already changed to '%s'", s.Key, s.OriginalValue))
		}
		break
	}
}

func createStringReplyActor(c *gin.Context, wg *sync.WaitGroup) *actor.PID {
	return actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
		dispatchStringReply(c, ctx, wg)
	}))
}
