package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

// BroadcastStopGroup holds PIDs of actors to stop in the group.
type BroadcastStopGroup struct {
	Routee []*actor.PID
}

// Stop sends stop message to all actors in the group.
func (g *BroadcastStopGroup) Stop() {
	for _, r := range g.Routee {
		r.Stop()
	}
}

// NewBroadcastStopGroup creates new BroadcastStopGroup.
func NewBroadcastStopGroup(routees []*actor.PID) *BroadcastStopGroup {
	return &BroadcastStopGroup{Routee: routees}
}
