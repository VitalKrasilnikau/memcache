package contracts

// CacheKeysContract is a data contract used for the list of all cache keys.
type CacheKeysContract struct {
	Keys []string `json:"keys"`
}

// ErrorContract is a data contract used for custom error message.
type ErrorContract struct {
	Status string `json:"status"`
}
