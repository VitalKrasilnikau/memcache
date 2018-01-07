package main

import (
	"flag"
	"fmt"
	"log"
	_ "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/VitalKrasilnikau/memcache/core/actors"
)

var (
	port = flag.Int("port", 1, "port to start the node at")
	nodeType = flag.String("type", "", "node type: string, list or dictionary")
	nodeIndex = flag.Int("index", 1, "node index in cluster")
	usePersistence = flag.Bool("persist", true, "use DB persistence")
)

func main() {
	flag.Parse()
	p := fmt.Sprintf("127.0.0.1:%d", *port)
	started := false
	switch *nodeType {
	case "string":
		remote.Start(p)
		act.NewStringCacheActor("memcache", *nodeIndex, *usePersistence)
		started = true
		break
	case "list":
		remote.Start(p)
		act.NewListCacheActor("memcache", *nodeIndex, *usePersistence)
		started = true
		break
	case "dictionary":
		remote.Start(p)
		act.NewDictionaryCacheActor("memcache", *nodeIndex, *usePersistence)
		started = true
		break
	}
	if started {
		log.Printf("Started %s%d node on port %s\n", *nodeType, *nodeIndex, p)
		for {} // TODO: Support graceful shutdown
		remote.Shutdown(true)
		log.Printf("Stopped %s%d node on port %s\n", *nodeType, *nodeIndex, p)
	}
}