package contracts

// StringCacheValueContract is used to serialize string cache entry via API.
type StringCacheValueContract struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewStringCacheValueContract is used to add new string cache entry using API.
type NewStringCacheValueContract struct {
	Key   string `form:"key" json:"key" binding:"required"`
	Value string `form:"value" json:"value" binding:"required"`
	TTL   string `form:"ttl" json:"ttl"`
}

// UpdateStringCacheValueContract is used to update string cache entry using API.
type UpdateStringCacheValueContract struct {
	NewValue      string `form:"value" json:"value" binding:"required"`
	OriginalValue string `form:"original" json:"original" binding:"required"`
}
