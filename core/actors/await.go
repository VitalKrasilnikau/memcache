package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"sync"
)

// ReplyProducer is a function which creates reply actor.
type ReplyProducer func(*sync.WaitGroup) *actor.PID

// MessageProducer is a function which creates request message.
type MessageProducer func() interface{}

// Await is used to request reply from new actor created in the controller.
func Await(pid *actor.PID, getReply ReplyProducer, getMessage MessageProducer) {
	var wg sync.WaitGroup
	wg.Add(1)
	replyPid := getReply(&wg)
	pid.Request(getMessage(), replyPid)
	wg.Wait()
	replyPid.Stop()
}
