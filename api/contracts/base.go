package contracts

type CacheKeysContract struct {
	Keys []string `json:"keys"`
}

type ErrorContract struct {
	Status string `json:"status"`
}
