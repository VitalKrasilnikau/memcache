package main

import (
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