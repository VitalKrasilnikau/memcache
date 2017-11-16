package main

import (
	"time"
	"fmt"
	"github.com/VitalKrasilnikau/memcache/client/apiclient"
)

func main() {
	apiClient := apiclient.APIClient{Host:"http://localhost", Port:8080}

	keys, err := apiClient.GetStringKeys()
	printKeys(keys, err, "string keys:")

	keys, err = apiClient.GetListKeys()
	printKeys(keys, err, "list keys:")

	keys, err = apiClient.GetDictionaryKeys()
	printKeys(keys, err, "dictionary keys:")

	b, dkey, err := apiClient.GetDictionaryKey("string")
	if b {
		printJSON(dkey, err)
	}

	ok, e, err := apiClient.PostStringKey("test123", "val", 0 * time.Second)
	if !ok {
		fmt.Println("sdf"+e.Status)
	}
	if err != nil {
		fmt.Println(err.Error())
	}

	b, skey, err := apiClient.GetStringKey("test123")
	if b {
		printJSON(skey, err)
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