# Memory cache based on Go, Protoactor, MongoDB persistence, Gin-based REST API with swagger

## Core

The cache is based on distributed actor system with consistent hashing router of **Protoactor**. There are three hash groups: one for string cache, one for list cache and one for dictionary cache. These groups are created by `NewStringCacheActorCluster`, `NewListCacheActorCluster` and `NewDictionaryCacheActorCluster` functions. There are 10 actors in each group by default and they all run on the local machine. If the user requests all keys stored in e.g. string cache, all the actors are asked for their keys and all the results are merged before returning to the end user. See `BroadcastStringKeysGroup` for details.

## Persistence

When the API is stopped `actor.Stop` message is sent to each actor to persist the memory cache to MongoDB collection. There is no correspondence between hash keys and MongoDB collection names in this version, so if the cluster is rebalanced previous state can't be restored from MongoDB. **Labix** driver is used for that purpose. No data replication in this version.

When some data was successfully persisted the following message is printed in the log:

`[StringCacheDBEntry] Persisted data to memcache.strings4 successfully.`

When some data was successfully restored the following message is printed in the log:

`[ListCacheDBEntry] Read snapshot from memcache.lists6 successfully.`

## Web API

Web REST API is available for the following cache operations. The API is based on **Gin** and **Gin-swagger**.

![swagger](https://raw.githubusercontent.com/VitalKrasilnikau/memcache/master/swagger.png)

Swagger endpoint is available on `/swagger/index.html` route.
At this moment not all endpoints can be tested from swagger UI due to some issues in Gin-swagger, please use `curl` or Go API cache client.

## Build the project

`$ go build api/main.go`

## Run the project

Run the following line to start memory cache server on port 8080 and MongoDB persistence:

1. `$ sudo mongod &`
1. `$ ./main`

Run the following line to start memory cache server custom port and MongoDB persistence:

1. `$ sudo mongod &`
1. `$ ./main %PORT%`

Run the following line to start memory cache server custom port and without MongoDB persistence:

`$ ./main %PORT% no-db`

Run the following line to start memory cache server on port 8080 and without MongoDB persistence:

`$ ./main no-db`

Press `CTRL-C` to stop the server.

## API cache client