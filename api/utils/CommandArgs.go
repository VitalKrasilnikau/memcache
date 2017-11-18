package api

import (
	"os"
	"strconv"
)

const (
	defaultPort = "8080"
	noDb        = "no-db"
)

// CommandArgs is a structure holding parameters from console.
type CommandArgs struct {
	UsePersistence bool
	Port           string
}

// NewCommandArgs parses the console parameters.
func NewCommandArgs() CommandArgs {
	args := os.Args[1:]
	switch len(args) {
	case 1:
		if args[0] == noDb {
			return CommandArgs{Port: defaultPort, UsePersistence: false}
		}
		_, e := strconv.Atoi(args[0])
		if e == nil {
			return CommandArgs{Port: args[0], UsePersistence: true}
		}
		return CommandArgs{Port: defaultPort, UsePersistence: true}
	case 2:
		_, e := strconv.Atoi(args[0])
		if e == nil && args[1] == noDb {
			return CommandArgs{Port: args[0], UsePersistence: false}
		}
		return CommandArgs{Port: defaultPort, UsePersistence: true}
	default:
		return CommandArgs{Port: defaultPort, UsePersistence: true}
	}
}
