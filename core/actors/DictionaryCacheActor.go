package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/core/cache"
	"github.com/VitalKrasilnikau/memcache/core/repository"
	"log"
	"time"
)

// GetDictionaryCacheKeyMessage is used to get the dictionary cache entry.
type GetDictionaryCacheKeyMessage struct {
	Key     string
	ReplyTo *actor.PID
}

// Hash is used for partitioning in actor cluster.
func (m *GetDictionaryCacheKeyMessage) Hash() string {
	return m.Key
}

// GetDictionaryCacheKeyReply is a reply message for GetDictionaryCacheKeyMessage.
type GetDictionaryCacheKeyReply struct {
	Key     string
	Values  []cache.KeyValue
	Success bool
}

// Hash is used for partitioning in actor cluster.
func (m *GetDictionaryCacheKeyReply) Hash() string {
	return m.Key
}

// DeleteDictionaryCacheKeyMessage is used to request the cache item deletion.
type DeleteDictionaryCacheKeyMessage struct {
	Key     string
	ReplyTo *actor.PID
}

// Hash is used for partitioning in actor cluster.
func (m *DeleteDictionaryCacheKeyMessage) Hash() string {
	return m.Key
}

// DeleteDictionaryCacheKeyReply is a reply message for DeleteDictionaryCacheKeyMessage.
type DeleteDictionaryCacheKeyReply struct {
	Key           string
	DeletedValues []cache.KeyValue
	Success       bool
}

// PostDictionaryCacheKeyMessage is used to add new dictionary cache entry.
type PostDictionaryCacheKeyMessage struct {
	Key     string
	Values  []cache.KeyValue
	TTL     time.Duration
	ReplyTo *actor.PID
}

// Hash is used for partitioning in actor cluster.
func (m *PostDictionaryCacheKeyMessage) Hash() string {
	return m.Key
}

// PostDictionaryCacheKeyReply is a reply message for PostDictionaryCacheKeyMessage.
type PostDictionaryCacheKeyReply struct {
	Key     string
	Success bool
}

// Hash is used for partitioning in actor cluster.
func (m *PostDictionaryCacheKeyReply) Hash() string {
	return m.Key
}

// PutDictionaryCacheValueMessage is used to request cache entry update by original value.
type PutDictionaryCacheValueMessage struct {
	Key           string
	SubKey        string
	NewValue      string
	OriginalValue string
	ReplyTo       *actor.PID
}

// Hash is used for partitioning in actor cluster.
func (m *PutDictionaryCacheValueMessage) Hash() string {
	return m.Key
}

// PutDictionaryCacheValueReply is a reply message for PutDictionaryCacheValueMessage.
type PutDictionaryCacheValueReply struct {
	Key           string
	SubKey        string
	Success       bool
	NewValue      string
	OriginalValue string
}

// DeleteDictionaryCacheValueMessage is used to request cache entry update by deleting the original value.
type DeleteDictionaryCacheValueMessage struct {
	Key     string
	SubKey  string
	ReplyTo *actor.PID
}

// Hash is used for partitioning in actor cluster.
func (m *DeleteDictionaryCacheValueMessage) Hash() string {
	return m.Key
}

// DeleteDictionaryCacheValueReply is a reply message for DeleteDictionaryCacheValueMessage.
type DeleteDictionaryCacheValueReply struct {
	Key          string
	SubKey       string
	DeletedValue cache.KeyValue
	Success      bool
}

// PostDictionaryCacheValueMessage is used to request cache entry update by original value.
type PostDictionaryCacheValueMessage struct {
	Key      string
	NewValue cache.KeyValue
	ReplyTo  *actor.PID
}

// Hash is used for partitioning in actor cluster.
func (m *PostDictionaryCacheValueMessage) Hash() string {
	return m.Key
}

// PostDictionaryCacheValueReply is a reply message for PostDictionaryCacheValueMessage.
type PostDictionaryCacheValueReply struct {
	Key        string
	Success    bool
	AddedValue cache.KeyValue
}

// DictionaryCacheActor manages partitioned dictionary cache and its persistence.
type DictionaryCacheActor struct {
	ClusterName    string
	NodeName       string
	Cache          cache.IDictionaryCache
	CachePersister cache.IDictionaryCachePersistence
	DB             repo.IDictionaryCacheRepository
}

// Receive is DictionaryCacheActor messages handler.
func (a *DictionaryCacheActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *GetDictionaryCacheKeyMessage:
		ok, v := a.Cache.TryGet(msg.Key)
		context.Tell(msg.ReplyTo, GetDictionaryCacheKeyReply{Key: msg.Key, Values: v, Success: ok})
		break
	case *DeleteDictionaryCacheKeyMessage:
		ok, v := a.Cache.TryDelete(msg.Key)
		context.Tell(msg.ReplyTo, DeleteDictionaryCacheKeyReply{Key: msg.Key, DeletedValues: v, Success: ok})
		log.Printf("[DictionaryCacheActor] Deleted %s", msg.Key)
		break
	case *GetCacheKeysMessage:
		context.Respond(GetCacheKeysReply{Keys: a.Cache.GetKeys()})
		break
	case *PostDictionaryCacheKeyMessage:
		ok := a.Cache.TryAdd(msg.Key, msg.Values, msg.TTL)
		context.Tell(msg.ReplyTo, PostDictionaryCacheKeyReply{Key: msg.Key, Success: ok})
		log.Printf("[DictionaryCacheActor] Created %s [%v]", msg.Key, msg.TTL)
		break
	case *PostDictionaryCacheValueMessage:
		ok, _ := a.Cache.TryAddValue(msg.Key, msg.NewValue)
		context.Tell(msg.ReplyTo, PostDictionaryCacheValueReply{Key: msg.Key, Success: ok, AddedValue: msg.NewValue})
		if ok {
			log.Printf("[DictionaryCacheActor] Added value %s to list %s", msg.NewValue, msg.Key)
		}
		break
	case *PutDictionaryCacheValueMessage:
		ok, _ := a.Cache.TryUpdateValue(msg.Key, msg.SubKey, msg.NewValue, msg.OriginalValue)
		context.Tell(msg.ReplyTo, PutDictionaryCacheValueReply{Key: msg.Key, Success: ok, NewValue: msg.NewValue, OriginalValue: msg.OriginalValue, SubKey: msg.SubKey})
		if ok {
			log.Printf("[DictionaryCacheActor] Updated value %s to %s in list %s", msg.OriginalValue, msg.NewValue, msg.Key)
		}
		break
	case *DeleteDictionaryCacheValueMessage:
		ok, del := a.Cache.TryDeleteValue(msg.Key, msg.SubKey)
		context.Tell(msg.ReplyTo, DeleteDictionaryCacheValueReply{Key: msg.Key, DeletedValue: del, Success: ok, SubKey: msg.SubKey})
		if ok {
			log.Printf("[DictionaryCacheActor] Deleted subkey %s in dictionary %s", msg.SubKey, msg.Key)
		}
		break
	case *actor.Stopping:
		a.persistSnapshot()
		break
	}
}

func (a *DictionaryCacheActor) restoreSnapshot() {
	for _, entry := range a.DB.GetAll() {
		mappedItem := cache.DictionaryCacheEntry{
			Map: cache.ToMap(fromDB(entry.Values)),
			CacheEntryData: cache.CacheEntryData{
				Added:       entry.Added,
				Updated:     entry.Updated,
				ExpireAfter: entry.ExpireAfter,
				Persisted:   true}}
		a.CachePersister.TryAddFromSnapshot(entry.Key, mappedItem)
	}
}

func (a *DictionaryCacheActor) persistSnapshot() {
	var newItems []repo.DictionaryCacheDBEntry
	var updatedItems []repo.DictionaryCacheDBEntry
	for _, k := range a.Cache.GetKeys() {
		ok, v := a.CachePersister.TryGetSnapshot(k)
		if ok {
			mappedItem := repo.DictionaryCacheDBEntry{
				Key:         k,
				Values:      toDB(cache.FromMap(v.Map)),
				Added:       v.Added,
				Updated:     v.Updated,
				ExpireAfter: v.ExpireAfter}
			if v.Persisted {
				updatedItems = append(updatedItems, mappedItem)
			} else {
				newItems = append(newItems, mappedItem)
			}
		}
	}
	if newItems == nil {
		newItems = make([]repo.DictionaryCacheDBEntry, 0)
	}
	if updatedItems == nil {
		updatedItems = make([]repo.DictionaryCacheDBEntry, 0)
	}
	a.DB.SaveAll(newItems, updatedItems)
}

func toDB(values []cache.KeyValue) []repo.DictionaryValueDBEntry {
	var a = make([]repo.DictionaryValueDBEntry, len(values))
	for i, v := range values {
		a[i] = repo.DictionaryValueDBEntry{Key: v.Key, Value: v.Value}
	}
	return a
}

func fromDB(values []repo.DictionaryValueDBEntry) []cache.KeyValue {
	var a = make([]cache.KeyValue, len(values))
	for i, v := range values {
		a[i] = cache.KeyValue{Key: v.Key, Value: v.Value}
	}
	return a
}
