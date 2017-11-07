package act

import (
	"time"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/core/cache"
	"github.com/VitalKrasilnikau/memcache/core/repository"
)

type GetStringCacheKeyMessage struct{
	Key string
}
type GetStringCacheKeyReply struct{
	Key string
	Value string
	Success bool
}

type DeleteStringCacheKeyMessage struct{
	Key string
}
type DeleteStringCacheKeyReply struct{
	Key string
	DeletedValue string
	Success bool
}

type GetStringCacheKeysMessage struct{}
type GetStringCacheKeysReply struct{
	Keys []string
}

type PostStringCacheKeyMessage struct{
	Key string
	Value string
	Ttl time.Duration
}
type PostStringCacheKeyReply struct{
	Key string
	Success bool
}

type PutStringCacheKeyMessage struct{
	Key string
	NewValue string
	OriginalValue string
}
type PutStringCacheKeyReply struct{
	Key string
	OriginalValue string
	Success bool
}

func NewStringCacheActor(clusterName string) *actor.PID {
	a := StringCacheActor{ ClusterName: clusterName }
	a.Init()
	props := actor.FromInstance(&a)
	return actor.Spawn(props)
}

type StringCacheActor struct{
	ClusterName string
	Cache cache.StringCache
	DB repo.StringCacheRepository
}

func (a *StringCacheActor) Init() {
	a.Cache = cache.StringCache{ Map: make(map[string]cache.StringCacheEntry) }
	a.DB = repo.StringCacheRepository{ Host: "localhost", DBName: a.ClusterName, ColName: "strings" }
	a.restoreSnapshot()
}

func (a *StringCacheActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case GetStringCacheKeyMessage:
		ok, v := a.Cache.TryGet(msg.Key)
		context.Respond(GetStringCacheKeyReply{ Key: msg.Key, Value: v, Success: ok })
		break
	case DeleteStringCacheKeyMessage:
		ok, v := a.Cache.TryDelete(msg.Key)
		context.Respond(DeleteStringCacheKeyReply{ Key: msg.Key, DeletedValue: v, Success: ok })
		break
	case GetStringCacheKeysMessage:
		context.Respond(GetStringCacheKeysReply{ Keys: a.Cache.GetKeys() })
		break
	case PostStringCacheKeyMessage:
		ok := a.Cache.TryAdd(msg.Key, msg.Value, msg.Ttl)
		context.Respond(PostStringCacheKeyReply{ Key: msg.Key, Success: ok })
		break
	case PutStringCacheKeyMessage:
		ok, v := a.Cache.TryUpdate(msg.Key, msg.NewValue, msg.OriginalValue)
		context.Respond(PutStringCacheKeyReply{ Key: msg.Key, OriginalValue: v, Success: ok })
		break
	case *actor.Stopping:
		a.persistSnapshot()
		break
	}
}

func (a *StringCacheActor) restoreSnapshot() {
	for _, entry := range a.DB.GetAll() {
		mappedItem := cache.StringCacheEntry{
			Added: entry.Added,
			Updated: entry.Updated,
			ExpireAfter: entry.ExpireAfter,
			Value: entry.Value,
			Persisted: true }
		a.Cache.TryAddFromSnapshot(entry.Key, mappedItem)
	}
}

func (a *StringCacheActor) persistSnapshot() {
	var newItems []repo.StringCacheDBEntry
	var updatedItems []repo.StringCacheDBEntry
	for _, k := range a.Cache.GetKeys() {
		ok, v := a.Cache.TryGetSnapshot(k)
		if (ok) {
			mappedItem := repo.StringCacheDBEntry{
				Key: k,
				Added: v.Added,
				Updated: v.Updated,
				ExpireAfter: v.ExpireAfter,
				Value: v.Value}
			if v.Persisted {
				updatedItems = append(updatedItems, mappedItem)
			} else {
				newItems = append(newItems, mappedItem)
			}
		}
	}
	if newItems == nil {
		newItems = make([]repo.StringCacheDBEntry, 0)
	}
	if updatedItems == nil {
		updatedItems = make([]repo.StringCacheDBEntry, 0)
	}
	a.DB.SaveAll(newItems, updatedItems)
}