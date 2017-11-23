package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

// GetCacheKeysMessage is used to request all cache keys.
type GetCacheKeysMessage struct {
	ReplyTo *actor.PID
}

// GetCacheKeysReply is a reply message for GetCacheKeysMessage.
type GetCacheKeysReply struct {
	Keys []string
}
