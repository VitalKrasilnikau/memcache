package api

import (
	"os"
	"strconv"
)

const (
	defaultPort = "8080"
	noDb        = "no-db"
	actorNumber = 10
)

// CommandArgs is a structure holding parameters from console.
type CommandArgs struct {
	UsePersistence bool
	Port           string
	ActorNumber    int
	IsRemote			 bool
}

// NewCommandArgs parses the console parameters.
func NewCommandArgs() CommandArgs {
	args := os.Args[1:]
	switch len(args) {
	case 1:
		if args[0] == noDb {
			return CommandArgs{Port: defaultPort, UsePersistence: false, ActorNumber: actorNumber}
		}
		_, e := strconv.Atoi(args[0])
		if e == nil {
			return CommandArgs{Port: args[0], UsePersistence: true, ActorNumber: actorNumber}
		}
		return CommandArgs{Port: defaultPort, UsePersistence: true, ActorNumber: actorNumber}
	case 2:
		_, e := strconv.Atoi(args[0])
		if e == nil && args[1] == noDb {
			return CommandArgs{Port: args[0], UsePersistence: false, ActorNumber: actorNumber}
		}
		return CommandArgs{Port: defaultPort, UsePersistence: true, ActorNumber: actorNumber}
	case 3:
		_, e := strconv.Atoi(args[0])
		n, e2 := strconv.Atoi(args[2])
		if e == nil && e2 == nil && args[1] == noDb {
			return CommandArgs{Port: args[0], UsePersistence: false, ActorNumber: n}
		}
		return CommandArgs{Port: defaultPort, UsePersistence: true, ActorNumber: actorNumber}
		case 4:
		_, e := strconv.Atoi(args[0])
		n, e2 := strconv.Atoi(args[2])
		if e == nil && e2 == nil && args[1] == noDb && args[3] == "remote" {
			return CommandArgs{Port: args[0], UsePersistence: false, ActorNumber: n, IsRemote: true}
		}
		return CommandArgs{Port: defaultPort, UsePersistence: true, ActorNumber: actorNumber}
	default:
		return CommandArgs{Port: defaultPort, UsePersistence: true, ActorNumber: actorNumber}
	}
}
