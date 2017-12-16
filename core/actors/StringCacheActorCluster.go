package act

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
)

// NewStringCacheActorCluster is a constructor function for the cluster of StringCacheActor.
func NewStringCacheActorCluster(clusterName string, nodeNumber int, usePersistence bool, isRemote bool) (*actor.PID, *BroadcastStopGroup, *BroadcastStringKeysGroup) {
	var nodes = make([]*actor.PID, nodeNumber)
	for i := 0; i < nodeNumber; i++ {
		if isRemote {
			nodes[i] = actor.NewPID(fmt.Sprintf("127.0.0.1:%d", i + 100), fmt.Sprintf("strings%d", i))
		} else {
			nodes[i] = factory.CreateStringCacheActor(clusterName, fmt.Sprintf("strings%d", i), usePersistence)
		}
	}
	return actor.Spawn(router.NewConsistentHashGroup(nodes...)),
		NewBroadcastStopGroup(nodes),
		NewBroadcastStringKeysGroup(nodes)
}

// NewStringCacheActor creates actor instance for remote connection.
func NewStringCacheActor(clusterName string, nodeNumber int, usePersistence bool) *actor.PID {
		return factory.CreateStringCacheActor(clusterName, fmt.Sprintf("strings%d", nodeNumber), usePersistence)
}
