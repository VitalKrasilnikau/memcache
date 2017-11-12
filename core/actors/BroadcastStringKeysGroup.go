package act

import (
	"errors"
	"github.com/AsynkronIT/protoactor-go/actor"
	"strings"
	"time"
)

// BroadcastStringKeysGroup holds PIDs of actors to get string keys in the group.
type BroadcastStringKeysGroup struct {
	Routee []*actor.PID
}

// Request selects string keys from all actors in the group.
func (g *BroadcastStringKeysGroup) Request(message *GetCacheKeysMessage, timeout time.Duration) (GetCacheKeysReply, error) {
	var keys []string
	var errorsStrings []string
	var resultError error
	for _, pid := range g.Routee {
		res, e := pid.RequestFuture(message, timeout).Result()
		if e == nil {
			s, ok := res.(GetCacheKeysReply)
			if ok && len(s.Keys) > 0 {
				keys = append(keys, s.Keys...)
			}
		} else {
			errorsStrings = append(errorsStrings, e.Error())
		}
	}
	if errorsStrings != nil {
		resultError = errors.New(strings.Join(errorsStrings, "\n"))
	}
	if keys == nil {
		keys = make([]string, 0)
	}
	return GetCacheKeysReply{Keys: keys}, resultError
}

// NewBroadcastStringKeysGroup creates new BroadcastStringKeysGroup.
func NewBroadcastStringKeysGroup(routees []*actor.PID) *BroadcastStringKeysGroup {
	return &BroadcastStringKeysGroup{Routee: routees}
}
