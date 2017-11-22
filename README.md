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

The following endpoints support TTL options in the payload:

1. `POST /api/string`
1. `POST /api/list`
1. `POST /api/dictionary`

TTL field is a string which can have values in two possible formats: `hh:mm` or `hh:mm:ss`. If the field is not specified or malformed, TTL is set to zero which means no expiration.

At this moment not all endpoints can be tested from swagger UI due to some issues in Gin-swagger, please use `curl` or Go API cache client.

## Build the project

1. `$ go get github.com/VitalKrasilnikau/memcache`
1. `$ cd src/github.com/VitalKrasilnikau/memcache/`
1. `$ go get ./...`
1. `$ go build api/main.go`

## Run the project

Run the following line to start memory cache server on port 8080, with 10 actors per cache and MongoDB persistence:

1. `$ sudo mongod &`
1. `$ ./main`

Run the following line to start memory cache server custom port and MongoDB persistence:

1. `$ sudo mongod &`
1. `$ ./main %PORT%`

Run the following line to start memory cache server on custom port, with custom number of actors per cache and without MongoDB persistence:

`$ ./main %PORT% no-db %ACTORS_NUMBER%`

Run the following line to start memory cache server on port 8080 and without MongoDB persistence:

`$ ./main no-db`

Press `CTRL-C` to stop the server.

## API cache client

ApiClient is a go client for the cache REST API. It is implemented using **resty** lib with HTTP mode switched on to allow redirects.

Example:

```go
apiClient := apiclient.APIClient{Host:"http://localhost", Port:8080}
keys, err := apiClient.GetStringKeys()
```

See details in `/client` folder.

## API cache load test

Run the server with custom number of actors and GOMAXPROCS to have good performance during load test.

`$ GOMAXPROCS=1000 ./main 8080 no-db 1000`

See details in `/loadtest` folder.