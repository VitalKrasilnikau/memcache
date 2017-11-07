package cache

import (
	"time"
)

// StringCacheEntry is a string cache data item stored in the memory cache.
type StringCacheEntry struct {
	Value       string
	ExpireAfter int64
	Added       int64
	Updated     int64
	Persisted   bool
}

// StringCache is a single-thread in-memory cache based on map[string]string.
type StringCache struct {
	Map map[string]StringCacheEntry
}

// TryGet returns the value if contains the key specified.
func (c *StringCache) TryGet(key string) (bool, string) {
	v, ok := c.getValueWithExpiration(key)
	return ok, v.Value
}

// TryGetSnapshot returns the value if contains the key specified.
func (c *StringCache) TryGetSnapshot(key string) (bool, StringCacheEntry) {
	v, ok := c.getValueWithExpiration(key)
	return ok, v
}

// TryAdd add new value to the cache by the key specified if the key is not already used.
func (c *StringCache) TryAdd(key string, value string, ttl time.Duration) bool {
	_, ok := c.getValueWithExpiration(key)
	if !ok {
		nowTime := time.Now()
		now := nowTime.Unix()
		var expireAfter int64
		if ttl > 0 {
			expireAfter = nowTime.Add(ttl).Unix()
		}
		c.Map[key] = StringCacheEntry{Value: value, Updated: now, Added: now, ExpireAfter: expireAfter}
	}
	return !ok
}

// TryAddFromSnapshot add new value to the cache by the key specified if the key is not already used.
func (c *StringCache) TryAddFromSnapshot(key string, entry StringCacheEntry) bool {
	_, ok := c.Map[key]
	if !ok {
		c.Map[key] = entry
	}
	return !ok
}

// TryDelete deletes the value by the key specified if the key is already used.
func (c *StringCache) TryDelete(key string) (bool, string) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		delete(c.Map, key)
	}
	return ok, v.Value
}

// TryUpdate updates the value by the key only if the value is the same as the original value.
// If not, you should repeat the update operation using optimistic concurrency loop
func (c *StringCache) TryUpdate(key string, newValue string, originalValue string) (bool, string) {
	v, ok := c.getValueWithExpiration(key)
	if ok && v.Value == originalValue {
		entry := StringCacheEntry{
			Value:       newValue,
			Added:       v.Added,
			ExpireAfter: v.ExpireAfter,
			Persisted:   v.Persisted,
			Updated:     time.Now().Unix()}
		c.Map[key] = entry
		return true, v.Value
	}
	return false, v.Value
}

// GetKeys returns all the keys in the map.
func (c *StringCache) GetKeys() []string {
	var keySlice []string
	now := time.Now().Unix()
	for key, v := range c.Map {
		if now <= v.ExpireAfter {
			keySlice = append(keySlice, key)
		}
	}
	if keySlice == nil {
		return make([]string, 0)
	}
	return keySlice
}

func (c *StringCache) getValueWithExpiration(key string) (StringCacheEntry, bool) {
	v, ok := c.Map[key]
	if ok {
		if v.ExpireAfter == 0 {
			return v, ok // no expiration
		}
		now := time.Now().Unix()
		if now > v.ExpireAfter {
			delete(c.Map, key)
			return v, false // expired
		}
		return v, ok // not expired
	}
	return v, ok
}
