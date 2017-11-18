package api

import (
	"fmt"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Error is 500 status response handler.
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, contracts.ErrorContract{Status: message})
}

// Bad is 400 status response handler.
func Bad(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, contracts.ErrorContract{Status: message})
}

// NotFound is 404 status response handler.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, contracts.ErrorContract{Status: message})
}

// NoContent is 204 status response handler.
func NoContent(c *gin.Context) {
	c.String(http.StatusNoContent, "")
}

// Created is 201 status response handler.
func Created(c *gin.Context) {
	c.String(http.StatusCreated, "")
}

// OK is 200 status response handler.
func OK(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

// DurationToString converts Go duration into string of format hh:mm:ss
func DurationToString(ttl time.Duration) string {
	return fmt.Sprintf("%.0f:%.0f:%.0f", ttl.Hours(), ttl.Minutes(), ttl.Seconds())
}

// ParseDuration converts string values like "01:05:20" or "11:05" (hh:mm) to go duration.
func ParseDuration(durationString string) time.Duration {
	if durationString == "" {
		return 0
	}
	parts := strings.Split(durationString, ":")
	switch len(parts) {
	case 2:
		h, e := strconv.Atoi(parts[0])
		if e != nil || h < 0 || h > 23 {
			return 0
		}
		m, e := strconv.Atoi(parts[1])
		if e != nil || m < 0 || m > 59 {
			return 0
		}
		return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
	case 3:
		h, e := strconv.Atoi(parts[0])
		if e != nil || h < 0 || h > 23 {
			return 0
		}
		m, e := strconv.Atoi(parts[1])
		if e != nil || m < 0 || m > 59 {
			return 0
		}
		s, e := strconv.Atoi(parts[2])
		if e != nil || s < 0 || s > 59 {
			return 0
		}
		return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	default:
		return 0
	}
}
