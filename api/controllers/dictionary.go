package controllers

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/VitalKrasilnikau/memcache/core/cache"
	"github.com/gin-gonic/gin"
	"sync"
)

// GetDictionaryCacheKeyHandler API which gets dictionary cache entry by key.
func GetDictionaryCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		var wg sync.WaitGroup
		wg.Add(1)
		replyPid := actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
			s, ok := ctx.Message().(act.GetDictionaryCacheKeyReply)
			if ok {
				if s.Success {
					api.OK(c, contracts.DictionaryCacheValueContract{Key: s.Key, Values: toDto(s.Values)})
				} else {
					api.NotFound(c, fmt.Sprintf("key '%s' was not found", key))
				}
				wg.Done()
			}
		}))
		pid.Tell(&act.GetDictionaryCacheKeyMessage{Key: key, ReplyTo: replyPid})
		wg.Wait()
		replyPid.Stop()
	}
}

// DeleteDictionaryCacheKeyHandler API which deletes list cache entry by key.
func DeleteDictionaryCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		var wg sync.WaitGroup
		wg.Add(1)
		replyPid := actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
			s, ok := ctx.Message().(act.DeleteDictionaryCacheKeyReply)
			if ok {
				if s.Success {
					api.NoContent(c)
				} else {
					api.NotFound(c, fmt.Sprintf("key '%s' was not found", key))
				}
				wg.Done()
			}
		}))
		pid.Tell(&act.DeleteDictionaryCacheKeyMessage{Key: key, ReplyTo: replyPid})
		wg.Wait()
		replyPid.Stop()
	}
}

// PostDictionaryCacheKeyHandler API which posts new list value.
func PostDictionaryCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		var json contracts.NewDictionaryCacheValuesContract
		if err := c.ShouldBindJSON(&json); err == nil {
			var wg sync.WaitGroup
			wg.Add(1)
			replyPid := actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
				s, ok := ctx.Message().(act.PostDictionaryCacheKeyReply)
				if ok {
					if s.Success {
						api.Created(c)
					} else {
						api.Bad(c, fmt.Sprintf("key '%s' was already used", json.Key))
					}
					wg.Done()
				}
			}))
			message := &act.PostDictionaryCacheKeyMessage{
				Key:     json.Key,
				Values:  fromDto(json.Values),
				TTL:     api.ParseDuration(json.TTL),
				ReplyTo: replyPid}
			pid.Tell(message)
			wg.Wait()
			replyPid.Stop()
		} else {
			api.Bad(c, fmt.Sprintf("malformed request: %s", err.Error()))
		}
	}
}

// PutDictionaryCacheValueHandler API which updates existing string value in the list by the key and new value specified in body.
func PutDictionaryCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		subkey := c.Param("subkey")
		var json contracts.UpdateDictionaryCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			var wg sync.WaitGroup
			wg.Add(1)
			replyPid := actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
				s, ok := ctx.Message().(act.PutDictionaryCacheValueReply)
				if ok {
					if s.Success {
						api.NoContent(c)
					} else {
						api.Bad(c, fmt.Sprintf("dictionary subkey '%s' of key '%s' was already changed or never existed", subkey, key))
					}
					wg.Done()
				}
			}))
			entry := &act.PutDictionaryCacheValueMessage{
				Key:           key,
				SubKey:        subkey,
				NewValue:      json.Value,
				OriginalValue: json.Original,
				ReplyTo:       replyPid}
			pid.Tell(entry)
			wg.Wait()
			replyPid.Stop()
		} else {
			api.Bad(c, fmt.Sprintf("malformed request: %s", err.Error()))
		}
	}
}

// PostDictionaryCacheValueHandler API which updates existing list value by the key and new added value specified in body.
func PostDictionaryCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		var json contracts.AddDictionaryCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			var wg sync.WaitGroup
			wg.Add(1)
			replyPid := actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
				s, ok := ctx.Message().(act.PostDictionaryCacheValueReply)
				if ok {
					if s.Success {
						api.Created(c)
					} else {
						api.Bad(c, fmt.Sprintf("key '%s' was not found", key))
					}
					wg.Done()
				}
			}))
			entry := &act.PostDictionaryCacheValueMessage{
				Key:      key,
				NewValue: cache.KeyValue{Key: json.Value.Key, Value: json.Value.Value},
				ReplyTo:  replyPid}
			pid.Tell(entry)
			wg.Wait()
			replyPid.Stop()
		} else {
			api.Bad(c, fmt.Sprintf("malformed request: %s", err.Error()))
		}
	}
}

// DeleteDictionaryCacheValueHandler API which updates existing list value by the key and deleting the value specified.
func DeleteDictionaryCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("subkey")
		var wg sync.WaitGroup
		wg.Add(1)
		replyPid := actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
			s, ok := ctx.Message().(act.DeleteDictionaryCacheValueReply)
			if ok {
				if s.Success {
					api.NoContent(c)
				} else {
					api.Bad(c, fmt.Sprintf("dictionary value '%s' of key '%s' was already deleted or never existed", value, key))
				}
				wg.Done()
			}
		}))
		entry := &act.DeleteDictionaryCacheValueMessage{Key: key, SubKey: value, ReplyTo: replyPid}
		pid.Tell(entry)
		wg.Wait()
		replyPid.Stop()
	}
}

func toDto(values []cache.KeyValue) []contracts.DictionaryKeyValueContract {
	var a = make([]contracts.DictionaryKeyValueContract, len(values))
	for i, v := range values {
		a[i] = contracts.DictionaryKeyValueContract{Key: v.Key, Value: v.Value}
	}
	return a
}

func fromDto(values []contracts.DictionaryKeyValueContract) []cache.KeyValue {
	var a = make([]cache.KeyValue, len(values))
	for i, v := range values {
		a[i] = cache.KeyValue{Key: v.Key, Value: v.Value}
	}
	return a
}
