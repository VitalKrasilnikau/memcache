package main

import (
	"fmt"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/client/apiclient"
	"time"
	"sync"
)

var (
	apiClient = apiclient.APIClient{Host: "http://localhost", Port: 8080}
)

func main() {
	numberOfClients := 300
	numberOfSessions := 10000
	guard := make(chan struct{}, numberOfClients)
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(numberOfSessions)
	for i := 0; i < numberOfSessions; i++ {
		guard <- struct{}{}
		go func(n int) {
			defer wg.Done()
			simulate(n)
			<-guard
		}(i)
	}
	wg.Wait()
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Total time: %v", elapsed)
}

func simulate(index int) {
	dictionaryKey := fmt.Sprintf("aaa%d", index)

	// Create new dictionary in the cache
	ok, e, err := apiClient.PostDictionaryKey(
		dictionaryKey,
		[]contracts.DictionaryKeyValueContract{contracts.DictionaryKeyValueContract{Key: "key1", Value: "val1"}},
		0*time.Second)
	if !ok {
		fmt.Println(e.Status)
	}
	if err != nil {
		fmt.Println(err.Error())
	}

	// Make sure it was created
	b, dkey, err := apiClient.GetDictionaryKey(dictionaryKey)
	if !b {
		printJSON(dkey, err)
	}

	// Add new value to the dictionary
	b, res, err := apiClient.PostDictionaryValue(dictionaryKey, contracts.DictionaryKeyValueContract{Key: "key2", Value: "val2"})
	if !b {
		printJSON(res, err)
	}

	// Modify first value in the dictionary
	b, res, err = apiClient.PutDictionaryValue(dictionaryKey, "key1", contracts.UpdateDictionaryCacheValueContract{Original: "val1", Value: "val3"})
	if !b {
		printJSON(res, err)
	}

	// Make sure it was updated
	b, dkey, err = apiClient.GetDictionaryKey(dictionaryKey)
	if !b {
		printJSON(dkey, err)
	}

	// Delete last value in the dictionary
	b, res, err = apiClient.DeleteDictionaryValue(dictionaryKey, "key2")
	if !b {
		printJSON(res, err)
	}

	// Make sure it was updated
	b, dkey, err = apiClient.GetDictionaryKey(dictionaryKey)
	if !b {
		printJSON(dkey, err)
	}

	// TODO: switch to Await
	// List all dictionary keys
	/*keys, err := apiClient.GetDictionaryKeys()
	if err != nil {
		printKeys(keys, err, "dictionary keys:")
	}*/

	// Delete the dictionary value by the key
	b, res, err = apiClient.DeleteDictionaryKey(dictionaryKey)
	if !b {
		printJSON(res, err)
	}

	// Make sure it was deleted
	b, dkey, err = apiClient.GetDictionaryKey(dictionaryKey)
	if b {
		printJSON(dkey, err)
	}
}

func printJSON(keyObj interface{}, err error) {
	if err == nil {
		fmt.Printf("\t%+v\n", keyObj)
	} else {
		fmt.Println(err.Error())
	}
}

func printKeys(keys []string, err error, message string) {
	if err == nil {
		fmt.Println(message)
		for _, k := range keys {
			fmt.Printf("\t%s\n", k)
		}
	} else {
		fmt.Println(err.Error())
	}
}
