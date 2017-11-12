package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	_ "github.com/VitalKrasilnikau/memcache/api/docs"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"os"
	"os/signal"
	"time"
)

var defTimeout = 50*time.Millisecond

/* String handlers */

// GetStringCacheKeyHandler API which gets string cache entry by key.
// @Description gets string cache entry by key
// @Accept   json
// @Produce  json
// @Param    key	path	string	true	"key"
// @Success 200 {object} contracts.StringCacheValueContract	"key and corresponding value"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/string/{key} [get]
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
// @Description deletes string cache entry by key
// @Accept   json
// @Produce  json
// @Param    deleted-key	path	string	true	"deleted-key"
// @Success 204 {string} string "no content"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/string/{deleted-key} [delete]
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

// GetCacheKeysHandler API which gets all string cache keys.
// @Description gets all string cache keys
// @Accept   json
// @Produce  json
// @Success 200 {object} contracts.CacheKeysContract	"string cache keys"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string [get]
func GetCacheKeysHandler(pid *act.BroadcastStringKeysGroup) func(*gin.Context) {
	return func(c *gin.Context) {
		res, e := pid.Request(&act.GetCacheKeysMessage{}, defTimeout)
		if (e == nil) {
			api.OK(c, contracts.CacheKeysContract{Keys: res.Keys})
		} else {
			api.Error(c, e.Error())
		}
	}
}

// PostStringCacheKeyHandler API which posts new string value.
// @Description posts new string value
// @Accept   json
// @Produce  json
// @Param    body	body	contracts.NewStringCacheValueContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string/ [post]
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
// @Description updates existing string value by the key and new value specified
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    body	body	contracts.UpdateStringCacheValueContract	true	"body"
// @Success 204 {string} string	"no content"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string/{update-key} [put]
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

/* List handlers */

// GetListCacheKeyHandler API which gets list cache entry by key.
// @Description gets list cache entry by key
// @Accept   json
// @Produce  json
// @Param    key	path	string	true	"key"
// @Success 200 {object} contracts.ListCacheValueContract	"key and corresponding list value"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/list/{key} [get]
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
// @Description deletes list cache entry by key
// @Accept   json
// @Produce  json
// @Param    deleted-key	path	string	true	"deleted-key"
// @Success 204 {string} string "no content"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/list/{deleted-key} [delete]
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

// GetListKeysHandler API which gets all list cache keys.
// @Description gets all list cache keys
// @Accept   json
// @Produce  json
// @Success 200 {object} contracts.CacheKeysContract	"string cache keys"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list [get]
func GetListKeysHandler(pid *act.BroadcastStringKeysGroup) func(*gin.Context) {
	return GetCacheKeysHandler(pid)
}

// PostListCacheKeyHandler API which posts new list value.
// @Description posts new list value
// @Accept   json
// @Produce  json
// @Param    body	body	contracts.NewListCacheValuesContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/ [post]
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
// @Description updates existing string value in the list by the key and new value specified in body
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    update-value	path	string	true	"update-value"
// @Param    body	body	contracts.UpdateListCacheValueContract	true	"body"
// @Success 204 {string} string	"no content"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/{update-key}/{update-value} [put]
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
						api.Bad(c, fmt.Sprintf("list value '%s' of key '%s' was already changed or never existed", key, value))
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
// @Description updates existing list value by the key and new added value specified in body
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    body	body	contracts.UpdateListCacheValueContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/{update-key} [post]
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
// @Description updates existing list value by the key and deleting the value specified
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    delete-value	path	string	true	"delete-value"
// @Success 204 {string} string	"no content"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/{update-key}/{delete-value} [delete]
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
					api.Bad(c, fmt.Sprintf("list value '%s' of key '%s' was already deleted or never existed", key, value))
				}
			} else {
				api.Error(c, "can't parse data")
			}
		} else {
			api.Error(c, e.Error())
		}
	}
}

func main() {
	pid, bpid, cpid := act.NewStringCacheActorCluster("memcache", 10)
	lpid, lbpid, lcpid := act.NewListCacheActorCluster("memcache", 10)
	router := gin.Default()
	api := router.Group("/api")
	{
		str := api.Group("/string")
		{
			str.GET("/", GetCacheKeysHandler(cpid))
			str.GET("/:key", GetStringCacheKeyHandler(pid))
			str.POST("/", PostStringCacheKeyHandler(pid))
			str.PUT("/:key", PutStringCacheKeyHandler(pid))
			str.DELETE("/:key", DeleteStringCacheKeyHandler(pid))
		}
		list := api.Group("/list")
		{
			list.GET("/", GetCacheKeysHandler(lcpid))
			list.GET("/:key", GetListCacheKeyHandler(lpid))
			list.POST("/", PostListCacheKeyHandler(lpid))
			list.POST("/:key", PostListCacheValueHandler(lpid))
			list.PUT("/:key/:value", PutListCacheValueHandler(lpid))
			list.DELETE("/:key", DeleteListCacheKeyHandler(lpid))
			list.DELETE("/:key/:value", DeleteListCacheValueHandler(lpid))
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stop the service and save data to DB", sig)
			bpid.Stop()
			lbpid.Stop()
			time.Sleep(1*time.Second)
			os.Exit(0)
		}
	}()

	router.Run(":8080")
}
