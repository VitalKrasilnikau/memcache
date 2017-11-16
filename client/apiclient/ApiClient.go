package apiclient

import (
	"gopkg.in/resty.v1"
	"github.com/VitalKrasilnikau/memcache/api/contracts"
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