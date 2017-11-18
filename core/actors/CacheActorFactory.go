package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/VitalKrasilnikau/memcache/core/cache"
	"github.com/VitalKrasilnikau/memcache/core/repository"
)

// CacheActorFactory is a factory which creates actors and their dependencies.
type CacheActorFactory struct{}

// CreateStringCacheActor is a constructor function for StringCacheActor.
func (f CacheActorFactory) CreateStringCacheActor(clusterName string, nodeName string, usePersistence bool) *actor.PID {
	a := StringCacheActor{ClusterName: clusterName, NodeName: nodeName}
	a.Cache = cache.StringCache{Map: make(map[string]cache.StringCacheEntry)}
	if usePersistence {
		a.DB = repo.StringCacheRepository{Host: "localhost", DBName: a.ClusterName, ColName: a.NodeName}
	} else {
		a.DB = repo.EmptyStringCacheRepository{}
	}
	a.restoreSnapshot()
	props := actor.FromInstance(&a)
	return actor.Spawn(props)
}

// CreateListCacheActor is a constructor function for ListCacheActor.
func (f CacheActorFactory) CreateListCacheActor(clusterName string, nodeName string, usePersistence bool) *actor.PID {
	a := ListCacheActor{ClusterName: clusterName, NodeName: nodeName}
	a.Cache = cache.ListCache{Map: make(map[string]cache.ListCacheEntry)}
	if usePersistence {
		a.DB = repo.ListCacheRepository{Host: "localhost", DBName: a.ClusterName, ColName: a.NodeName}
	} else {
		a.DB = repo.EmptyListCacheRepository{}
	}
	a.restoreSnapshot()
	props := actor.FromInstance(&a)
	return actor.Spawn(props)
}

// CreateDictionaryCacheActor is a constructor function for DictionaryCacheActor.
func (f CacheActorFactory) CreateDictionaryCacheActor(clusterName string, nodeName string, usePersistence bool) *actor.PID {
	a := DictionaryCacheActor{ClusterName: clusterName, NodeName: nodeName}
	a.Cache = cache.DictionaryCache{Map: make(map[string]cache.DictionaryCacheEntry)}
	if usePersistence {
		a.DB = repo.DictionaryCacheRepository{Host: "localhost", DBName: a.ClusterName, ColName: a.NodeName}
	} else {
		a.DB = repo.EmptyDictionaryCacheRepository{}
	}
	a.restoreSnapshot()
	props := actor.FromInstance(&a)
	return actor.Spawn(props)
}
