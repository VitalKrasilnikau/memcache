package main

import (
	"time"
	"fmt"
	"github.com/VitalKrasilnikau/memcache/client/apiclient"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
)

func main() {
	apiClient := apiclient.APIClient{Host:"http://localhost", Port:8080}
	{
		stringKey := "test123"

		// Create new string value in the cache
		ok, e, err := apiClient.PostStringKey(stringKey, "val", 0 * time.Second)
		if !ok {
			fmt.Println(e.Status)
		} else {
			fmt.Printf("'%s' was created\n", stringKey)
		}
		if err != nil {
			fmt.Println(err.Error())
		}

		// Make sure it was created
		b, skey, err := apiClient.GetStringKey(stringKey)
		if b {
			printJSON(skey, err)
		}

		// Modify the string value specifying the latest known value
		b, res, err := apiClient.PutStringKey(stringKey, "val2", "val")
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", stringKey)
		}

		// Make sure it was updated
		b, skey, err = apiClient.GetStringKey(stringKey)
		if b {
			printJSON(skey, err)
		}

		// List all string keys
		keys, err := apiClient.GetStringKeys()
		printKeys(keys, err, "string keys:")

		// Delete the string value by the key
		b, res, err = apiClient.DeleteStringKey(stringKey)
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was deleted\n", stringKey)
		}

		// Make sure it was deleted
		b, skey, err = apiClient.GetStringKey(stringKey)
		if b {
			printJSON(skey, err)
		}
	}
	{
		listKey := "list123"

		// Create new list in the cache
		ok, e, err := apiClient.PostListKey(listKey, []string { "val1" }, 0 * time.Second)
		if !ok {
			fmt.Println(e.Status)
		} else {
			fmt.Printf("'%s' was created\n", listKey)
		}
		if err != nil {
			fmt.Println(err.Error())
		}

		// Make sure it was created
		b, lkey, err := apiClient.GetListKey(listKey)
		if b {
			printJSON(lkey, err)
		}

		// Add new value to the list
		b, res, err := apiClient.PostListValue(listKey, "val2")
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", listKey)
		}

		// Modify first value to the list
		b, res, err = apiClient.PutListValue(listKey, "val3", "val1")
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", listKey)
		}

		// Make sure it was updated
		b, lkey, err = apiClient.GetListKey(listKey)
		if b {
			printJSON(lkey, err)
		}

		// Delete last value in the list
		b, res, err = apiClient.DeleteListValue(listKey, "val2")
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", listKey)
		}

		// Make sure it was updated
		b, lkey, err = apiClient.GetListKey(listKey)
		if b {
			printJSON(lkey, err)
		}

		// List all list keys
		keys, err := apiClient.GetListKeys()
		printKeys(keys, err, "list keys:")

		// Delete the list value by the key
		b, res, err = apiClient.DeleteListKey(listKey)
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was deleted\n", listKey)
		}

		// Make sure it was deleted
		b, lkey, err = apiClient.GetListKey(listKey)
		if b {
			printJSON(lkey, err)
		}
	}
	{
		dictionaryKey := "dictionary123"

		// Create new dictionary in the cache
		ok, e, err := apiClient.PostDictionaryKey(
			dictionaryKey,
			[]contracts.DictionaryKeyValueContract{contracts.DictionaryKeyValueContract{Key:"key1", Value:"val1"}},
			0 * time.Second)
		if !ok {
			fmt.Println(e.Status)
		} else {
			fmt.Printf("'%s' was created\n", dictionaryKey)
		}
		if err != nil {
			fmt.Println(err.Error())
		}

		// Make sure it was created
		b, dkey, err := apiClient.GetDictionaryKey(dictionaryKey)
		if b {
			printJSON(dkey, err)
		}

		// Add new value to the dictionary
		b, res, err := apiClient.PostDictionaryValue(dictionaryKey, contracts.DictionaryKeyValueContract{ Key: "key2", Value:"val2" })
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", dictionaryKey)
		}

		// Modify first value in the dictionary
		b, res, err = apiClient.PutDictionaryValue(dictionaryKey, "key1", contracts.UpdateDictionaryCacheValueContract{ Original:"val1", Value:"val3"})
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", dictionaryKey)
		}

		// Make sure it was updated
		b, dkey, err = apiClient.GetDictionaryKey(dictionaryKey)
		if b {
			printJSON(dkey, err)
		}

		// Delete last value in the dictionary
		b, res, err = apiClient.DeleteDictionaryValue(dictionaryKey, "key2")
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was updated\n", dictionaryKey)
		}

		// Make sure it was updated
		b, dkey, err = apiClient.GetDictionaryKey(dictionaryKey)
		if b {
			printJSON(dkey, err)
		}

		// List all dictionary keys
		keys, err := apiClient.GetDictionaryKeys()
		printKeys(keys, err, "dictionary keys:")

		// Delete the dictionary value by the key
		b, res, err = apiClient.DeleteDictionaryKey(dictionaryKey)
		if !b {
			printJSON(res, err)
		} else {
			fmt.Printf("'%s' was deleted\n", dictionaryKey)
		}

		// Make sure it was deleted
		b, dkey, err = apiClient.GetDictionaryKey(dictionaryKey)
		if b {
			printJSON(dkey, err)
		}
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