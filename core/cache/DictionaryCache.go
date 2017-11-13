package cache

import (
	"log"
	"time"
)

// DictionaryCacheEntry is a string cache data item stored in the memory cache.
type DictionaryCacheEntry struct {
	Map map[string]string
	CacheEntryData
}

// KeyValue is a dictionary cache item.
type KeyValue struct {
	Key   string
	Value string
}

// IDictionaryCache is an interface for DictionaryCache.
type IDictionaryCache interface {
	TryGet(key string) (bool, []KeyValue)
	TryGetSnapshot(key string) (bool, DictionaryCacheEntry)
	TryAdd(key string, values []KeyValue, ttl time.Duration) bool
	TryAddFromSnapshot(key string, entry DictionaryCacheEntry) bool
	TryDelete(key string) (bool, []KeyValue)
	TryUpdateValue(key string, subKey string, newValue string, originalValue string) (bool, []KeyValue)
	TryDeleteValue(key string, subKey string) (bool, KeyValue)
	TryAddValue(key string, newValue KeyValue) (bool, []KeyValue)
	GetKeys() []string
}

// DictionaryCache is a single-thread in-memory cache based on map[string]DictionaryCacheEntry.
type DictionaryCache struct {
	Map map[string]DictionaryCacheEntry
}

// TryGet returns the value if contains the key specified.
func (c DictionaryCache) TryGet(key string) (bool, []KeyValue) {
	v, ok := c.getValueWithExpiration(key)
	return ok, FromMap(v.Map)
}

// TryGetSnapshot returns the value if contains the key specified.
func (c DictionaryCache) TryGetSnapshot(key string) (bool, DictionaryCacheEntry) {
	v, ok := c.getValueWithExpiration(key)
	return ok, v
}

// TryAdd add new value to the cache by the key specified if the key is not already used.
func (c DictionaryCache) TryAdd(key string, values []KeyValue, ttl time.Duration) bool {
	_, ok := c.getValueWithExpiration(key)
	if !ok {
		c.Map[key] = DictionaryCacheEntry{Map: ToMap(values), CacheEntryData: NewCacheEntryData(ttl)}
	}
	return !ok
}

// TryAddFromSnapshot add new value to the cache by the key specified if the key is not already used.
func (c DictionaryCache) TryAddFromSnapshot(key string, entry DictionaryCacheEntry) bool {
	_, ok := c.Map[key]
	if !ok {
		c.Map[key] = entry
	}
	return !ok
}

// TryDelete deletes the value by the key specified if the key is already used.
func (c DictionaryCache) TryDelete(key string) (bool, []KeyValue) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		delete(c.Map, key)
	}
	return ok, FromMap(v.Map)
}

// TryUpdateValue updates the existing value in the list by the key.
func (c DictionaryCache) TryUpdateValue(key string, subKey string, newValue string, originalValue string) (bool, []KeyValue) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		prevValue, exists := v.Map[subKey]
		if exists {
			sameOriginal := prevValue == originalValue
			if sameOriginal {
				v.Map[subKey] = newValue
				entry := DictionaryCacheEntry{
					Map:            v.Map,
					CacheEntryData: UpdateCacheEntryData(v.CacheEntryData)}
				c.Map[key] = entry
			}
			return exists && sameOriginal, FromMap(v.Map)
		}
	}
	return false, FromMap(v.Map)
}

// TryDeleteValue deletes the existing value in the list by the key.
func (c DictionaryCache) TryDeleteValue(key string, subKey string) (bool, KeyValue) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		del, exists := v.Map[subKey]
		if exists {
			delete(v.Map, subKey)
			entry := DictionaryCacheEntry{
				Map:            v.Map,
				CacheEntryData: UpdateCacheEntryData(v.CacheEntryData)}
			c.Map[key] = entry
			return exists, KeyValue{Key: subKey, Value: del}
		}
	}
	return false, KeyValue{}
}

// TryAddValue adds the value to the list by the key.
func (c DictionaryCache) TryAddValue(key string, newValue KeyValue) (bool, []KeyValue) {
	v, ok := c.getValueWithExpiration(key)
	if ok {
		ok2 := addKeyToMap(v.Map, newValue)
		if ok2 {
			entry := DictionaryCacheEntry{
				Map:            v.Map,
				CacheEntryData: UpdateCacheEntryData(v.CacheEntryData)}
			c.Map[key] = entry
			return true, FromMap(entry.Map)
		}
	}
	return false, FromMap(v.Map)
}

// GetKeys returns all the keys in the map.
func (c DictionaryCache) GetKeys() []string {
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

func (c DictionaryCache) getValueWithExpiration(key string) (DictionaryCacheEntry, bool) {
	v, ok := c.Map[key]
	if ok {
		if IsCacheEntryExpired(v.CacheEntryData) {
			delete(c.Map, key)
			log.Printf("[DictionaryCache] key %s had expired and was removed", key)
			return v, false // expired
		}
		return v, ok // not expired
	}
	return v, ok
}

func addKeyToMap(m map[string]string, v KeyValue) bool {
	_, ok := m[v.Key]
	if !ok {
		m[v.Key] = v.Value
		return true
	}
	return false
}

// ToMap converts key/value array to map.
func ToMap(values []KeyValue) map[string]string {
	var m = make(map[string]string, len(values))
	for _, v := range values {
		addKeyToMap(m, v)
	}
	return m
}

// FromMap converts map to key/value array.
func FromMap(aMap map[string]string) []KeyValue {
	var a = make([]KeyValue, len(aMap))
	if aMap == nil {
		return a
	}
	i := 0
	for k, v := range aMap {
		a[i] = KeyValue{k, v}
		i++
	}
	return a
}
