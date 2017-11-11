package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	_ "github.com/VitalKrasilnikau/memcache/api/docs"
	"github.com/VitalKrasilnikau/memcache/core/actors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// GetGetStringCacheKeyHandler API which gets string cache entry by key.
// @Description gets string cache entry by key
// @Accept   json
// @Produce  json
// @Param    new-key	path	string	true	"new-key"
// @Success 200 {object} contracts.StringCacheValueContract	"key and corresponding value"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/string/{new-key} [get]
func GetGetStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		res, e := pid.RequestFuture(&act.GetStringCacheKeyMessage{Key: key}, 50*time.Millisecond).Result()
		if e == nil {
			s, ok := res.(act.GetStringCacheKeyReply)
			if ok {
				if s.Success {
					c.JSON(http.StatusOK, contracts.StringCacheValueContract{Key: s.Key, Value: s.Value})
				} else {
					c.JSON(http.StatusNotFound, contracts.ErrorContract{Status: fmt.Sprintf("key '%s' was not found", key)})
				}
			} else {
				c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: "can't parse data"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: e.Error()})
		}
	}
}

// GetDeleteStringCacheKeyHandler API which deletes string cache entry by key.
// @Description deletes string cache entry by key
// @Accept   json
// @Produce  json
// @Param    key	path	string	true	"key"
// @Success 204 {string} string "no content"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/string/{key} [delete]
func GetDeleteStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		res, e := pid.RequestFuture(&act.DeleteStringCacheKeyMessage{Key: key}, 50*time.Millisecond).Result()
		if e == nil {
			s, ok := res.(act.DeleteStringCacheKeyReply)
			if ok {
				if s.Success {
					c.String(http.StatusNoContent, "")
				} else {
					c.JSON(http.StatusNotFound, contracts.ErrorContract{Status: fmt.Sprintf("key '%s' was not found", key)})
				}
			} else {
				c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: "can't parse data"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: e.Error()})
		}
	}
}

// GetGetStringCacheKeysHandler API which gets all string cache keys.
// @Description gets all string cache keys
// @Accept   json
// @Produce  json
// @Success 200 {object} contracts.StringCacheKeysContract	"string cache keys"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string [get]
func GetGetStringCacheKeysHandler(pid *act.BroadcastStringKeysGroup) func(*gin.Context) {
	return func(c *gin.Context) {
		res, err := pid.Request(&act.GetStringCacheKeysMessage{}, 50*time.Millisecond)
		if (err == nil) {
			c.JSON(http.StatusOK, contracts.StringCacheKeysContract{Keys: res.Keys})
		} else {
			c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: err.Error()})
		}
	}
}

// GetPostStringCacheKeyHandler API which posts new string value.
// @Description posts new string value
// @Accept   json
// @Produce  json
// @Param    body	body	contracts.NewStringCacheValueContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string/ [post]
func GetPostStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		var json contracts.NewStringCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			message := &act.PostStringCacheKeyMessage{Key: json.Key, Value: json.Value, TTL: parseDurationFromJSON(json.TTL)}
			res, e := pid.RequestFuture(message, 50*time.Millisecond).Result()
			if e == nil {
				s, ok := res.(act.PostStringCacheKeyReply)
				if ok {
					if s.Success {
						c.String(http.StatusCreated, "")
					} else {
						c.JSON(http.StatusBadRequest, contracts.ErrorContract{Status: fmt.Sprintf("key '%s' was already used", json.Key)})
					}
				} else {
					c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: "can't parse data"})
				}
			} else {
				c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: e.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, contracts.ErrorContract{Status: fmt.Sprintf("malformed request: %s", err.Error())})
		}
	}
}

// GetPutStringCacheKeyHandler API which updates existing string value by the key and new value specified.
// @Description updates existing string value by the key and new value specified
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    body	body	contracts.UpdateStringCacheValueContract	true	"body"
// @Success 204 {string} string	"no content"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string/{update-key} [put]
func GetPutStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Param("key")
		var json contracts.UpdateStringCacheValueContract
		if err := c.ShouldBindJSON(&json); err == nil {
			entry := &act.PutStringCacheKeyMessage{Key: key, NewValue: json.NewValue, OriginalValue: json.OriginalValue}
			res, e := pid.RequestFuture(entry, 50*time.Millisecond).Result()
			if e == nil {
				s, ok := res.(act.PutStringCacheKeyReply)
				if ok {
					if s.Success {
						c.String(http.StatusNoContent, "")
					} else {
						c.JSON(http.StatusBadRequest, contracts.ErrorContract{Status: fmt.Sprintf("key '%s' was already changed to '%s'", key, s.OriginalValue)})
					}
				} else {
					c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: "can't parse data"})
				}
			} else {
				c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: e.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, contracts.ErrorContract{Status: fmt.Sprintf("malformed request: %s", err.Error())})
		}
	}
}

func parseDurationFromJSON(json string) time.Duration {
	return 0
}

func main() {
	pid, bpid, cpid := act.NewStringCacheActorCluster("memcache", 10)
	router := gin.Default()
	api := router.Group("/api")
	{
		str := api.Group("/string")
		{
			str.GET("/", GetGetStringCacheKeysHandler(cpid))
			str.GET("/:key", GetGetStringCacheKeyHandler(pid))
			str.POST("/", GetPostStringCacheKeyHandler(pid))
			str.PUT("/:key", GetPutStringCacheKeyHandler(pid))
			str.DELETE("/:key", GetDeleteStringCacheKeyHandler(pid))
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stop the service and save data to DB", sig)
			bpid.Stop()
			time.Sleep(1*time.Second)
			os.Exit(0)
		}
	}()

	router.Run(":8080")
}
