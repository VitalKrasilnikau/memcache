package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/goconsole"
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
	remote.Start(fmt.Sprintf("127.0.0.1:%d", *port))
	switch *nodeType {
	case "string":
		act.NewStringCacheActor("memcache", *nodeIndex, *usePersistence)
		console.ReadLine()
	case "list":
		act.NewListCacheActor("memcache", *nodeIndex, *usePersistence)
		console.ReadLine()
	case "dictionary":
		act.NewDictionaryCacheActor("memcache", *nodeIndex, *usePersistence)
		console.ReadLine()
	}
}