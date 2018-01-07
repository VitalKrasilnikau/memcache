package act

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
)

// NewDictionaryCacheActorCluster is a constructor function for the cluster of DictionaryCacheActor.
func NewDictionaryCacheActorCluster(clusterName string, nodeNumber int, usePersistence bool, isRemote bool) (*actor.PID, *BroadcastStopGroup, *BroadcastStringKeysGroup) {
	var nodes = make([]*actor.PID, nodeNumber)
	for i := 0; i < nodeNumber; i++ {
		if isRemote {
			nodes[i] = actor.NewPID(fmt.Sprintf("127.0.0.1:%d", i + 60000), fmt.Sprintf("strings%d", i))
		} else {
			nodes[i] = factory.CreateDictionaryCacheActor(clusterName, fmt.Sprintf("dictionaries%d", i), usePersistence)
		}
	}
	return actor.Spawn(router.NewConsistentHashGroup(nodes...)),
		NewBroadcastStopGroup(nodes),
		NewBroadcastStringKeysGroup(nodes)
}

// NewDictionaryCacheActor creates actor instance for remote connection.
func NewDictionaryCacheActor(clusterName string, nodeNumber int, usePersistence bool) *actor.PID {
		return factory.CreateDictionaryCacheActor(clusterName, fmt.Sprintf("dictionaries%d", nodeNumber), usePersistence)
}

