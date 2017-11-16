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

// APIClient is a go client lib for accessing memory cache.
type APIClient struct{
	Host string
	Port int32
}

// GetStringKeys returns all string keys in the cache.
func (c APIClient) GetStringKeys() ([]string, error) {
	return c.getKeys("string");
}

// GetListKeys returns all list keys in the cache.
func (c APIClient) GetListKeys() ([]string, error) {
	return c.getKeys("list");
}

// GetDictionaryKeys returns all dictionary keys in the cache.
func (c APIClient) GetDictionaryKeys() ([]string, error) {
	return c.getKeys("dictionary");
}

// GetDictionaryKey returns dictionary value by key from the cache.
func (c APIClient) GetDictionaryKey(key string) (bool, contracts.DictionaryCacheValueContract, error) {
	resp, err := c.getKey(fmt.Sprintf("dictionary/%s", key))
	if err != nil {
		return false, contracts.DictionaryCacheValueContract{}, err
	}
	var reply contracts.DictionaryCacheValueContract
	if err = json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return false, contracts.DictionaryCacheValueContract{}, err
	}
	return resp.StatusCode() == 200, reply, nil;
}

// GetListKey returns list value by key from the cache.
func (c APIClient) GetListKey(key string) (bool, contracts.ListCacheValueContract, error) {
	resp, err := c.getKey(fmt.Sprintf("list/%s", key))
	if err != nil {
		return false, contracts.ListCacheValueContract{}, err
	}
	var reply contracts.ListCacheValueContract
	if err = json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return false, contracts.ListCacheValueContract{}, err
	}
	return resp.StatusCode() == 200, reply, nil;
}

// GetStringKey returns string value by key from the cache.
func (c APIClient) GetStringKey(key string) (bool, contracts.StringCacheValueContract, error) {
	resp, err := c.getKey(fmt.Sprintf("string/%s", key))
	if err != nil {
		return false, contracts.StringCacheValueContract{}, err
	}
	var reply contracts.StringCacheValueContract
	if err = json.Unmarshal(resp.Body(), &reply); err != nil {
		log.Fatal("unmarshal failed: " + err.Error())
		return false, contracts.StringCacheValueContract{}, err
	}
	return resp.StatusCode() == 200, reply, nil;
}

// PostStringKey adds new string key and value to the cache.
func (c APIClient) PostStringKey(key string, value string, ttl time.Duration) (bool, contracts.ErrorContract, error) {
	req := contracts.NewStringCacheValueContract{Key: key, Value: value, TTL: api.DurationToString(ttl)}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL("string/"));
	return c.processResponse(resp, err, 201)
}

// PostListKey adds new list key and value to the cache.
func (c APIClient) PostListKey(key string, values []string, ttl time.Duration) (bool, contracts.ErrorContract, error) {
	req := contracts.NewListCacheValuesContract{Key: key, Values: values, TTL: api.DurationToString(ttl)}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL("list/"));
	return c.processResponse(resp, err, 201)
}

// PostDictionaryKey adds new dictionary key and value to the cache.
func (c APIClient) PostDictionaryKey(
		key string,
		values []contracts.DictionaryKeyValueContract,
		ttl time.Duration) (bool, contracts.ErrorContract, error) {
	req := contracts.NewDictionaryCacheValuesContract{Key: key, Values: values, TTL: api.DurationToString(ttl)}
	resp, err := resty.SetHTTPMode().R().
		SetBody(req).
		Post(c.buildURL("dictionary/"));
	return c.processResponse(resp, err, 201)
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