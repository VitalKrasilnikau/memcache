package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/core/cache"
	"github.com/VitalKrasilnikau/memcache/core/repository"
	"log"
	"time"
)

// GetListCacheKeyMessage is used to get the list cache entry.
type GetListCacheKeyMessage struct {
	Key string
}

// Hash is used for partitioning in actor cluster.
func (m *GetListCacheKeyMessage) Hash() string {
	return m.Key
}

// GetListCacheKeyReply is a reply message for GetListCacheKeyMessage.
type GetListCacheKeyReply struct {
	Key     string
	Values  []string
	Success bool
}

// Hash is used for partitioning in actor cluster.
func (m *GetListCacheKeyReply) Hash() string {
	return m.Key
}

// DeleteListCacheKeyMessage is used to request the cache item deletion.
type DeleteListCacheKeyMessage struct {
	Key string
}

// Hash is used for partitioning in actor cluster.
func (m *DeleteListCacheKeyMessage) Hash() string {
	return m.Key
}

// DeleteListCacheKeyReply is a reply message for DeleteListCacheKeyMessage.
type DeleteListCacheKeyReply struct {
	Key           string
	DeletedValues []string
	Success       bool
}

// PostListCacheKeyMessage is used to add new cache entry.
type PostListCacheKeyMessage struct {
	Key    string
	Values []string
	TTL    time.Duration
}

// Hash is used for partitioning in actor cluster.
func (m *PostListCacheKeyMessage) Hash() string {
	return m.Key
}

// PostListCacheKeyReply is a reply message for PostListCacheKeyMessage.
type PostListCacheKeyReply struct {
	Key     string
	Success bool
}

// Hash is used for partitioning in actor cluster.
func (m *PostListCacheKeyReply) Hash() string {
	return m.Key
}

// PutListCacheValueMessage is used to request cache entry update by original value.
type PutListCacheValueMessage struct {
	Key           string
	NewValue      string
	OriginalValue string
}

// Hash is used for partitioning in actor cluster.
func (m *PutListCacheValueMessage) Hash() string {
	return m.Key
}

// PutListCacheValueReply is a reply message for PutListCacheValueMessage.
type PutListCacheValueReply struct {
	Key           string
	Success       bool
	NewValue      string
	OriginalValue string
}

// DeleteListCacheValueMessage is used to request cache entry update by deleting the original value.
type DeleteListCacheValueMessage struct {
	Key   string
	Value string
}

// Hash is used for partitioning in actor cluster.
func (m *DeleteListCacheValueMessage) Hash() string {
	return m.Key
}

// DeleteListCacheValueReply is a reply message for DeleteListCacheValueMessage.
type DeleteListCacheValueReply struct {
	Key          string
	DeletedValue string
	Success      bool
}

// PostListCacheValueMessage is used to request cache entry update by original value.
type PostListCacheValueMessage struct {
	Key      string
	NewValue string
}

// Hash is used for partitioning in actor cluster.
func (m *PostListCacheValueMessage) Hash() string {
	return m.Key
}

// PostListCacheValueReply is a reply message for PostListCacheValueMessage.
type PostListCacheValueReply struct {
	Key        string
	Success    bool
	AddedValue string
}

// NewListCacheActor is a constructor function for ListCacheActor.
func NewListCacheActor(clusterName string, nodeName string, usePersistence bool) *actor.PID {
	a := ListCacheActor{ClusterName: clusterName, NodeName: nodeName}
	a.Init(usePersistence)
	props := actor.FromInstance(&a)
	return actor.Spawn(props)
}

// ListCacheActor manages partitioned string cache and its persistence.
type ListCacheActor struct {
	ClusterName string
	NodeName    string
	Cache       cache.IListCache
	DB          repo.IListCacheRepository
}

// Init restores cache entries snapshot from DB.
func (a *ListCacheActor) Init(usePersistence bool) {
	a.Cache = cache.ListCache{Map: make(map[string]cache.ListCacheEntry)}
	if usePersistence {
		a.DB = repo.ListCacheRepository{Host: "localhost", DBName: a.ClusterName, ColName: a.NodeName}
	} else {
		a.DB = repo.EmptyListCacheRepository{}
	}
	a.restoreSnapshot()
}

// Receive is ListCacheActor messages handler.
func (a *ListCacheActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *GetListCacheKeyMessage:
		ok, v := a.Cache.TryGet(msg.Key)
		context.Respond(GetListCacheKeyReply{Key: msg.Key, Values: v, Success: ok})
		break
	case *DeleteListCacheKeyMessage:
		ok, v := a.Cache.TryDelete(msg.Key)
		context.Respond(DeleteListCacheKeyReply{Key: msg.Key, DeletedValues: v, Success: ok})
		log.Printf("[ListCacheActor] Deleted %s", msg.Key)
		break
	case *GetCacheKeysMessage:
		context.Respond(GetCacheKeysReply{Keys: a.Cache.GetKeys()})
		break
	case *PostListCacheKeyMessage:
		ok := a.Cache.TryAdd(msg.Key, msg.Values, msg.TTL)
		context.Respond(PostListCacheKeyReply{Key: msg.Key, Success: ok})
		log.Printf("[ListCacheActor] Created %s [%v]", msg.Key, msg.TTL)
		break
	case *PostListCacheValueMessage:
		ok, _ := a.Cache.TryAddValue(msg.Key, msg.NewValue)
		context.Respond(PostListCacheValueReply{Key: msg.Key, Success: ok, AddedValue: msg.NewValue})
		if ok {
			log.Printf("[ListCacheActor] Added value %s to list %s", msg.NewValue, msg.Key)
		}
		break
	case *PutListCacheValueMessage:
		ok, _ := a.Cache.TryUpdateValue(msg.Key, msg.NewValue, msg.OriginalValue)
		context.Respond(PutListCacheValueReply{Key: msg.Key, Success: ok, NewValue: msg.NewValue, OriginalValue: msg.OriginalValue})
		if ok {
			log.Printf("[ListCacheActor] Updated value %s to %s in list %s", msg.OriginalValue, msg.NewValue, msg.Key)
		}
		break
	case *DeleteListCacheValueMessage:
		ok, _ := a.Cache.TryDeleteValue(msg.Key, msg.Value)
		context.Respond(DeleteListCacheValueReply{Key: msg.Key, DeletedValue: msg.Value, Success: ok})
		if ok {
			log.Printf("[ListCacheActor] Deleted value %s in list %s", msg.Value, msg.Key)
		}
		break
	case *actor.Stopping:
		a.persistSnapshot()
		break
	}
}

func (a *ListCacheActor) restoreSnapshot() {
	for _, entry := range a.DB.GetAll() {
		mappedItem := cache.ListCacheEntry{
			Values: entry.Values,
			CacheEntryData: cache.CacheEntryData{
				Added:       entry.Added,
				Updated:     entry.Updated,
				ExpireAfter: entry.ExpireAfter,
				Persisted:   true}}
		a.Cache.TryAddFromSnapshot(entry.Key, mappedItem)
	}
}

func (a *ListCacheActor) persistSnapshot() {
	var newItems []repo.ListCacheDBEntry
	var updatedItems []repo.ListCacheDBEntry
	for _, k := range a.Cache.GetKeys() {
		ok, v := a.Cache.TryGetSnapshot(k)
		if ok {
			mappedItem := repo.ListCacheDBEntry{
				Key:         k,
				Added:       v.CacheEntryData.Added,
				Updated:     v.CacheEntryData.Updated,
				ExpireAfter: v.CacheEntryData.ExpireAfter,
				Values:      v.Values}
			if v.Persisted {
				updatedItems = append(updatedItems, mappedItem)
			} else {
				newItems = append(newItems, mappedItem)
			}
		}
	}
	if newItems == nil {
		newItems = make([]repo.ListCacheDBEntry, 0)
	}
	if updatedItems == nil {
		updatedItems = make([]repo.ListCacheDBEntry, 0)
	}
	a.DB.SaveAll(newItems, updatedItems)
}
