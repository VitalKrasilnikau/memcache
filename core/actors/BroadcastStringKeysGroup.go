package act

import (
	"errors"
	"github.com/AsynkronIT/protoactor-go/actor"
	"strings"
	"sync"
)

// BroadcastStringKeysGroup holds PIDs of actors to get string keys in the group.
type BroadcastStringKeysGroup struct {
	Routee []*actor.PID
}

// Request selects string keys from all actors in the group in parallel.
func (g *BroadcastStringKeysGroup) Request() (GetCacheKeysReply, error) {
	var keys []string
	var errorsStrings []string
	var resultError error
	type futureResult struct {
		reply GetCacheKeysReply
		e     error
	}
	ch := make(chan futureResult, len(g.Routee))
	for _, pid := range g.Routee {
		go func(pid *actor.PID) {
			Await(
				pid,
				func(wg *sync.WaitGroup) *actor.PID {
					return actor.Spawn(actor.FromFunc(func(ctx actor.Context) {
						if s, ok := ctx.Message().(GetCacheKeysReply); ok {
							resp := futureResult{reply: s}
							ch <- resp
						}
					}))
				},
				func(replyPid *actor.PID) interface{} {
					return &GetCacheKeysMessage{ReplyTo: replyPid}
				})
		}(pid)
	}
	for range g.Routee {
		res := <-ch
		if res.e == nil {
			if len(res.reply.Keys) > 0 {
				keys = append(keys, res.reply.Keys...)
			}
		} else {
			errorsStrings = append(errorsStrings, res.e.Error())
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
