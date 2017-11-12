package cache

import (
	"time"
)

// CacheEntryData holds basic info about cache item.
type CacheEntryData struct {
	ExpireAfter int64
	Added       int64
	Updated     int64
	Persisted   bool
}

// NewCacheEntryData creates new CacheEntryData structure.
func NewCacheEntryData(ttl time.Duration) CacheEntryData {
	nowTime := time.Now()
	now := nowTime.Unix()
	var expireAfter int64
	if ttl > 0 {
		expireAfter = nowTime.Add(ttl).Unix()
	}
	return CacheEntryData{
		ExpireAfter: expireAfter,
		Added: now,
		Updated: now,
		Persisted: false}
}

// UpdateCacheEntryData returns an updated copy of CacheEntryData.
func UpdateCacheEntryData(v CacheEntryData) CacheEntryData {
	return CacheEntryData{
		Added:       v.Added,
		ExpireAfter: v.ExpireAfter,
		Persisted:   v.Persisted,
		Updated:     time.Now().Unix()}
}

// IsCacheEntryExpired returns true if ttl defined for this entry had elapced.
func IsCacheEntryExpired(v CacheEntryData) bool {
	now := time.Now().Unix()
	return v.ExpireAfter != 0 && now > v.ExpireAfter
}