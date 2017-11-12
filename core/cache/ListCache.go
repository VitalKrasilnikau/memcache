package cache

import (
	"time"
	"log"
)

// ListCacheEntry is a string cache data item stored in the memory cache.
type ListCacheEntry struct {
	Values      []string
	ExpireAfter int64
	Added       int64
	Updated     int64
	Persisted   bool
}

// ListCache is a single-thread in-memory cache based on map[string]ListCacheEntry.
type ListCache struct {
	Map map[string]ListCacheEntry
}

// TryGet returns the value if contains the key specified.
func (c *ListCache) TryGet(key string) (bool, []string) {
	v, ok := c.getValueWithExpiration(key)
	return ok, v.Values
}

// TryGetSnapshot returns the value if contains the key specified.
func (c *ListCache) TryGetSnapshot(key string) (bool, ListCacheEntry) {
	v, ok := c.getValueWithExpiration(key)
	return ok, v
}

// TryAdd add new value to the cache by the key specified if the key is not already used.
func (c *ListCache) TryAdd(key string, values []string, ttl time.Duration) bool {
	_, ok := c.getValueWithExpiration(key)
	if !ok {
		nowTime := time.Now()
		now := nowTime.Unix()
		var expireAfter int64
		if ttl > 0 {
			expireAfter = nowTime.Add(ttl).Unix()
		}
		c.Map[key] = ListCacheEntry{Values: values, Updated: now, Added: now, ExpireAfter: expireAfter}
	}
	return !ok
}

// TryAddFromSnapshot add new value to the cache by the key specified if the key is not already used.
func (c *ListCache) TryAddFromSnapshot(key string, entry ListCacheEntry) bool {
	_, ok := c.Map[key]
	if !ok {
		c.Map[key] = entry
	}
	return !ok
}

// TryDelete deletes the value by the key specified if the key is already used.
func (c *ListCache) TryDelete(key string) (bool, []string) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		delete(c.Map, key)
	}
	return ok, v.Values
}

// TryUpdateValue updates the existing value in the list by the key.
func (c *ListCache) TryUpdateValue(key string, newValue string, originalValue string) (bool, []string) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		newValues, updated := c.replaceValueInArray(v.Values, newValue, originalValue, false)
		if updated {
			entry := ListCacheEntry{
				Values:      newValues,
				Added:       v.Added,
				ExpireAfter: v.ExpireAfter,
				Persisted:   v.Persisted,
				Updated:     time.Now().Unix()}
			c.Map[key] = entry
		}
		return updated, v.Values
	}
	return false, v.Values
}

// TryDeleteValue deletes the existing value in the list by the key.
func (c *ListCache) TryDeleteValue(key string, value string) (bool, []string) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		newValues, deleted := c.replaceValueInArray(v.Values, "", value, true)
		if deleted {
			entry := ListCacheEntry{
				Values:      newValues,
				Added:       v.Added,
				ExpireAfter: v.ExpireAfter,
				Persisted:   v.Persisted,
				Updated:     time.Now().Unix()}
			c.Map[key] = entry
		}
		return deleted, v.Values
	}
	return false, v.Values
}

// TryAddValue adds the value to the list by the key.
func (c *ListCache) TryAddValue(key string, newValue string) (bool, []string) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		entry := ListCacheEntry{
			Values:      append(v.Values, newValue),
			Added:       v.Added,
			ExpireAfter: v.ExpireAfter,
			Persisted:   v.Persisted,
			Updated:     time.Now().Unix()}
		c.Map[key] = entry
		return true, v.Values
	}
	return false, v.Values
}

// GetKeys returns all the keys in the map.
func (c *ListCache) GetKeys() []string {
	var keySlice []string
	now := time.Now().Unix()
	for key, v := range c.Map {
		if v.ExpireAfter == 0 || now <= v.ExpireAfter {
			keySlice = append(keySlice, key)
		}
	}
	if keySlice == nil {
		return make([]string, 0)
	}
	return keySlice
}

func (c *ListCache) getValueWithExpiration(key string) (ListCacheEntry, bool) {
	v, ok := c.Map[key]
	if ok {
		if v.ExpireAfter == 0 {
			return v, ok // no expiration
		}
		now := time.Now().Unix()
		if now > v.ExpireAfter {
			delete(c.Map, key)
			log.Printf("[ListCache] key %s had expired and was removed", key)
			return v, false // expired
		}
		return v, ok // not expired
	}
	return v, ok
}

func (c *ListCache) replaceValueInArray(a []string, newValue string, originalValue string, delete bool) ([]string, bool) {
	var ret []string
	updated := false
	if a != nil {
		for _, i := range a {
			if i == originalValue {
				if !delete {
					ret = append(ret, newValue)
				}
				updated = true
			} else {
				ret = append(ret, i)
			}
		}
	}
	return ret, updated
}
