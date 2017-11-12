package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/api/controllers"
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

/* String handlers for swagger */

// GetStringCacheKeyHandler .
// @Description gets string cache entry by key
// @Summary gets string cache entry by key
// @Accept   json
// @Produce  json
// @Param    key	path	string	true	"key"
// @Success 200 {object} contracts.StringCacheValueContract	"key and corresponding value"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/string/{key} [get]
func GetStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.GetStringCacheKeyHandler(pid)
}

// DeleteStringCacheKeyHandler .
// @Description deletes string cache entry by key
// @Summary deletes string cache entry by key
// @Accept   json
// @Produce  json
// @Param    deleted-key	path	string	true	"deleted-key"
// @Success 204 {string} string "no content"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/string/{deleted-key} [delete]
func DeleteStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.DeleteStringCacheKeyHandler(pid)
}

// GetCacheKeysHandler .
// @Description gets all string cache keys
// @Summary gets all string cache keys
// @Accept   json
// @Produce  json
// @Success 200 {object} contracts.CacheKeysContract	"string cache keys"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string [get]
func GetCacheKeysHandler(pid *act.BroadcastStringKeysGroup) func(*gin.Context) {
	return controllers.GetCacheKeysHandler(pid)
}

// PostStringCacheKeyHandler .
// @Description posts new string value
// @Summary posts new string value
// @Accept   json
// @Produce  json
// @Param    body	body	contracts.NewStringCacheValueContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string/ [post]
func PostStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.PostStringCacheKeyHandler(pid)
}

// PutStringCacheKeyHandler .
// @Description updates existing string value by the key and new value specified
// @Summary updates existing string value by the key and new value specified
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    body	body	contracts.UpdateStringCacheValueContract	true	"body"
// @Success 204 {string} string	"no content"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/string/{update-key} [put]
func PutStringCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.PutStringCacheKeyHandler(pid)
}

/* List handlers for swagger */

// GetListCacheKeyHandler .
// @Description gets list cache entry by key
// @Summary gets list cache entry by key
// @Accept   json
// @Produce  json
// @Param    key	path	string	true	"key"
// @Success 200 {object} contracts.ListCacheValueContract	"key and corresponding list value"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/list/{key} [get]
func GetListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.GetListCacheKeyHandler(pid)
}

// DeleteListCacheKeyHandler .
// @Description deletes list cache entry by key
// @Summary deletes list cache entry by key
// @Accept   json
// @Produce  json
// @Param    deleted-key	path	string	true	"deleted-key"
// @Success 204 {string} string "no content"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Failure 404 {object} contracts.ErrorContract "key was not found"
// @Router /api/list/{deleted-key} [delete]
func DeleteListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.DeleteListCacheKeyHandler(pid)
}

// GetListKeysHandler .
// @Description gets all list cache keys
// @Summary gets all list cache keys
// @Accept   json
// @Produce  json
// @Success 200 {object} contracts.CacheKeysContract	"string cache keys"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list [get]
func GetListKeysHandler(pid *act.BroadcastStringKeysGroup) func(*gin.Context) {
	return controllers.GetCacheKeysHandler(pid)
}

// PostListCacheKeyHandler .
// @Description posts new list value
// @Summary posts new list value
// @Accept   json
// @Produce  json
// @Param    body	body	contracts.NewListCacheValuesContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/ [post]
func PostListCacheKeyHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.PostListCacheKeyHandler(pid)
}

// PutListCacheValueHandler .
// @Description updates existing string value in the list by the key and new value specified in body
// @Summary updates existing string value in the list by the key and new value specified in body
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
	return controllers.PutListCacheValueHandler(pid)
}

// PostListCacheValueHandler .
// @Description updates existing list value by the key and new added value specified in body
// @Summary updates existing list value by the key and new added value specified in body
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    body	body	contracts.UpdateListCacheValueContract	true	"body"
// @Success 201 {string} string	"created"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/{update-key} [post]
func PostListCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.PostListCacheValueHandler(pid)
}

// DeleteListCacheValueHandler .
// @Description updates existing list value by the key and deleting the value specified
// @Summary updates existing list value by the key and deleting the value specified
// @Accept   json
// @Produce  json
// @Param    update-key	path	string	true	"update-key"
// @Param    delete-value	path	string	true	"delete-value"
// @Success 204 {string} string	"no content"
// @Failure 400 {object} contracts.ErrorContract "bad request"
// @Failure 500 {object} contracts.ErrorContract "server error"
// @Router /api/list/{update-key}/{delete-value} [delete]
func DeleteListCacheValueHandler(pid *actor.PID) func(*gin.Context) {
	return controllers.DeleteListCacheValueHandler(pid)
}

// @title Memory cache based on Go Swagger API
// @version 1.0
// @description This is a memory cache based on Go.
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
			list.GET("/", GetListKeysHandler(lcpid))
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
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}()

	router.Run(":8080")
}
