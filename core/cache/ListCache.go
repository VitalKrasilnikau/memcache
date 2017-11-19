package cache

import (
	"log"
	"time"
)

// ListCacheEntry is a string cache data item stored in the memory cache.
type ListCacheEntry struct {
	Values []string
	CacheEntryData
}

// IListCache is an interface for ListCache.
type IListCache interface {
	TryGet(key string) (bool, []string)
	TryAdd(key string, values []string, ttl time.Duration) bool
	TryDelete(key string) (bool, []string)
	TryUpdateValue(key string, newValue string, originalValue string) (bool, []string)
	TryDeleteValue(key string, value string) (bool, []string)
	TryAddValue(key string, newValue string) (bool, []string)
	GetKeys() []string
}

// IListCachePersistence is an interface for persisting ListCache.
type IListCachePersistence interface {
	TryGetSnapshot(key string) (bool, ListCacheEntry)
	TryAddFromSnapshot(key string, entry ListCacheEntry) bool
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
		c.Map[key] = ListCacheEntry{Values: values, CacheEntryData: NewCacheEntryData(ttl)}
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
				Values:         newValues,
				CacheEntryData: UpdateCacheEntryData(v.CacheEntryData)}
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
				Values:         newValues,
				CacheEntryData: UpdateCacheEntryData(v.CacheEntryData)}
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
			Values:         append(v.Values, newValue),
			CacheEntryData: UpdateCacheEntryData(v.CacheEntryData)}
		c.Map[key] = entry
		return true, v.Values
	}
	return false, v.Values
}

// GetKeys returns all the keys in the map.
func (c *ListCache) GetKeys() []string {
	var keySlice []string
	for key, v := range c.Map {
		if !IsCacheEntryExpired(v.CacheEntryData) {
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
		if IsCacheEntryExpired(v.CacheEntryData) {
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
