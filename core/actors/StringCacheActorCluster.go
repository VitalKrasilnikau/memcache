package act

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
)

// NewStringCacheActorCluster is a constructor function for the cluster of StringCacheActor.
func NewStringCacheActorCluster(clusterName string, nodeNumber int, usePersistence bool) (*actor.PID, *BroadcastStopGroup, *BroadcastStringKeysGroup) {
	var nodes = make([]*actor.PID, nodeNumber)
	for i := 0; i < nodeNumber; i++ {
		nodes[i] = factory.CreateStringCacheActor(clusterName, fmt.Sprintf("strings%d", i), usePersistence)
	}
	return actor.Spawn(router.NewConsistentHashGroup(nodes...)),
		NewBroadcastStopGroup(nodes),
		NewBroadcastStringKeysGroup(nodes)
}
