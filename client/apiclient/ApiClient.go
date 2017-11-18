package apiclient

import (
	"time"
	"gopkg.in/resty.v1"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
	"github.com/VitalKrasilnikau/memcache/api/utils"
	"log"
	"fmt"
	"encoding/json"
)

const (
	stringEndpoint = "string/"
	listEndpoint = "list/"
	dictionaryEndpoint = "dictionary/"
)

// APIClient is a go client lib for accessing memory cache.
type APIClient struct{
	Host string
	Port int32
}

// GetStringKeys returns all string keys in the cache.
func (c APIClient) GetStringKeys() ([]string, error) {
	return c.getKeys(stringEndpoint)
}

// GetStringKey returns string value by key from the cache.
func (c APIClient) GetStringKey(key string) (bool, contracts.StringCacheValueContract, error) {
	resp, err := c.getKey(stringEndpoint + key)
	if err != nil {
		return false, contracts.StringCacheValueContract{}, err
	}
	var reply contracts.StringCacheValueContract
	if err = json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return false, contracts.StringCacheValueContract{}, err
	}
	return resp.StatusCode() == 200, reply, nil
}

// PostStringKey adds new string key and value to the cache.
func (c APIClient) PostStringKey(key string, value string, ttl time.Duration) (bool, contracts.ErrorContract, error) {
	req := contracts.NewStringCacheValueContract{Key: key, Value: value, TTL: api.DurationToString(ttl)}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL(stringEndpoint))
	return c.processResponse(resp, err, 201)
}

// PutStringKey updates string key with new value in the cache.
func (c APIClient) PutStringKey(key string, newValue string, originalValue string) (bool, contracts.ErrorContract, error) {
	req := contracts.UpdateStringCacheValueContract{NewValue: newValue, OriginalValue:originalValue}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Put(c.buildURL(stringEndpoint + key))
	return c.processResponse(resp, err, 204)
}

// DeleteStringKey removes string key from the cache.
func (c APIClient) DeleteStringKey(key string) (bool, contracts.ErrorContract, error) {
	return c.deleteKey(stringEndpoint + key)
}

// GetListKeys returns all list keys in the cache.
func (c APIClient) GetListKeys() ([]string, error) {
	return c.getKeys(listEndpoint)
}

// GetListKey returns list value by key from the cache.
func (c APIClient) GetListKey(key string) (bool, contracts.ListCacheValueContract, error) {
	resp, err := c.getKey(listEndpoint + key)
	if err != nil {
		return false, contracts.ListCacheValueContract{}, err
	}
	var reply contracts.ListCacheValueContract
	if err = json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return false, contracts.ListCacheValueContract{}, err
	}
	return resp.StatusCode() == 200, reply, nil
}

// DeleteListKey removes list from the cache.
func (c APIClient) DeleteListKey(key string) (bool, contracts.ErrorContract, error) {
	return c.deleteKey(listEndpoint + key)
}

// DeleteListValue removes list value from the cache.
func (c APIClient) DeleteListValue(key string, value string) (bool, contracts.ErrorContract, error) {
	return c.deleteKey(fmt.Sprintf("%s%s/%s", listEndpoint, key, value))
}

// PostListKey adds new list key and value to the cache.
func (c APIClient) PostListKey(key string, values []string, ttl time.Duration) (bool, contracts.ErrorContract, error) {
	req := contracts.NewListCacheValuesContract{Key: key, Values: values, TTL: api.DurationToString(ttl)}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL(listEndpoint))
	return c.processResponse(resp, err, 201)
}

// PostListValue adds new list value to the cache list entry.
func (c APIClient) PostListValue(key string, value string) (bool, contracts.ErrorContract, error) {
	req := contracts.UpdateListCacheValueContract{Value: value}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL(listEndpoint + key))
	return c.processResponse(resp, err, 201)
}

// PutListValue updates list value in the cache list entry.
func (c APIClient) PutListValue(key string, newValue string, originalValue string) (bool, contracts.ErrorContract, error) {
	req := contracts.UpdateListCacheValueContract{Value: newValue}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Put(c.buildURL(fmt.Sprintf("%s%s/%s", listEndpoint, key, originalValue)));
	return c.processResponse(resp, err, 204)
}

// GetDictionaryKeys returns all dictionary keys in the cache.
func (c APIClient) GetDictionaryKeys() ([]string, error) {
	return c.getKeys(dictionaryEndpoint)
}

// GetDictionaryKey returns dictionary value by key from the cache.
func (c APIClient) GetDictionaryKey(key string) (bool, contracts.DictionaryCacheValueContract, error) {
	resp, err := c.getKey(dictionaryEndpoint + key)
	if err != nil {
		return false, contracts.DictionaryCacheValueContract{}, err
	}
	var reply contracts.DictionaryCacheValueContract
	if err = json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return false, contracts.DictionaryCacheValueContract{}, err
	}
	return resp.StatusCode() == 200, reply, nil
}

// DeleteDictionaryKey removes dictionary from the cache.
func (c APIClient) DeleteDictionaryKey(key string) (bool, contracts.ErrorContract, error) {
	return c.deleteKey(dictionaryEndpoint + key)
}

// DeleteDictionaryValue removes dictionary value from the cache.
func (c APIClient) DeleteDictionaryValue(key string, value string) (bool, contracts.ErrorContract, error) {
	return c.deleteKey(fmt.Sprintf("%s%s/%s", dictionaryEndpoint, key, value))
}

// PostDictionaryKey adds new dictionary key and value to the cache.
func (c APIClient) PostDictionaryKey(
		key string,
		values []contracts.DictionaryKeyValueContract,
		ttl time.Duration) (bool, contracts.ErrorContract, error) {
	req := contracts.NewDictionaryCacheValuesContract{Key: key, Values: values, TTL: api.DurationToString(ttl)}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL(dictionaryEndpoint))
	return c.processResponse(resp, err, 201)
}

// PostDictionaryValue adds new value to the dictionary in the cache.
func (c APIClient) PostDictionaryValue(
		key string,
		value contracts.DictionaryKeyValueContract) (bool, contracts.ErrorContract, error) {
	resp, err := resty.SetHTTPMode().R().
		SetBody(contracts.AddDictionaryCacheValueContract{ Value: value}).
		Post(c.buildURL(dictionaryEndpoint + key))
	return c.processResponse(resp, err, 201)
}

// PutDictionaryValue updates dictionary value in the cache dictionary entry.
func (c APIClient) PutDictionaryValue(key string, subKey string, value contracts.UpdateDictionaryCacheValueContract) (bool, contracts.ErrorContract, error) {
	resp, err := resty.SetHTTPMode().R().
		SetBody(value).
		Put(c.buildURL(fmt.Sprintf("%s%s/%s", dictionaryEndpoint, key, subKey)));
	return c.processResponse(resp, err, 204)
}

func (c APIClient) processResponse(resp *resty.Response, err error, expectedCode int) (bool, contracts.ErrorContract, error) {
	if err != nil {
		log.Fatal("get failed: " + err.Error())
		return false, contracts.ErrorContract{}, err
	}
	if (resp.StatusCode() != expectedCode) {
		log.Printf(fmt.Sprintf("status: %d, code: %s", resp.StatusCode(), resp.Status()))
		var reply contracts.ErrorContract
		if err = json.Unmarshal(resp.Body(), &reply); err != nil {
			log.Fatal("unmarshal failed: " + err.Error())
			return false, contracts.ErrorContract{}, err
		}
		return false, reply, nil;
	}
	return true, contracts.ErrorContract{}, nil;
}

func (c APIClient) getKey(endpoint string) (*resty.Response, error) {
	var resp *resty.Response
	var err error
	if resp, err = resty.SetHTTPMode().R().Get(c.buildURL(endpoint)); err != nil {
		log.Fatal("get failed: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (c APIClient) getKeys(endpoint string) ([]string, error) {
	var resp *resty.Response
	var err error
	if resp, err = resty.SetHTTPMode().R().Get(c.buildURL(endpoint)); err != nil {
		log.Fatal("get failed: " + err.Error())
		return nil, err
	}
	var reply contracts.CacheKeysContract
	if err := json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return nil, err
	}
	return reply.Keys, nil;
}

func (c APIClient) buildURL(endpoint string) string {
	return fmt.Sprintf("%s:%d/api/%s", c.Host, c.Port, endpoint)
}

func (c APIClient) deleteKey(endpoint string) (bool, contracts.ErrorContract, error) {
	resp, err := resty.SetHTTPMode().R().Delete(c.buildURL(endpoint));
	return c.processResponse(resp, err, 204)
}