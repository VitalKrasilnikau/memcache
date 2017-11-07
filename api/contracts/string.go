package contracts

type StringCacheValueContract struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StringCacheKeysContract struct {
	Keys []string `json:"keys"`
}

type ErrorContract struct {
	Status string `json:"status"`
}

type NewStringCacheValueContract struct {
	Key   string `form:"key" json:"key" binding:"required"`
	Value string `form:"value" json:"value" binding:"required"`
	TTL   string `form:"ttl" json:"ttl"`
}

type UpdateStringCacheValueContract struct {
	NewValue      string `form:"value" json:"value" binding:"required"`
	OriginalValue string `form:"original" json:"original" binding:"required"`
}
