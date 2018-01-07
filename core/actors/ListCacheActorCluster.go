package act

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
)

// NewListCacheActorCluster is a constructor function for the cluster of ListCacheActor.
func NewListCacheActorCluster(clusterName string, nodeNumber int, usePersistence bool, isRemote bool) (*actor.PID, *BroadcastStopGroup, *BroadcastStringKeysGroup) {
	var nodes = make([]*actor.PID, nodeNumber)
	for i := 0; i < nodeNumber; i++ {
		if isRemote {
			nodes[i] = actor.NewPID(fmt.Sprintf("127.0.0.1:%d", i + 58000), fmt.Sprintf("lists%d", i))
		} else {
			nodes[i] = factory.CreateListCacheActor(clusterName, fmt.Sprintf("lists%d", i), usePersistence)
		}
	}
	return actor.Spawn(router.NewConsistentHashGroup(nodes...)),
		NewBroadcastStopGroup(nodes),
		NewBroadcastStringKeysGroup(nodes)
}

// NewListCacheActor creates actor instance for remote connection.
func NewListCacheActor(clusterName string, nodeNumber int, usePersistence bool) *actor.PID {
		return factory.CreateListCacheActor(clusterName, fmt.Sprintf("lists%d", nodeNumber), usePersistence)
}
