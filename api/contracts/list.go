package contracts

// ListCacheValueContract is used to serialize list cache entry via API.
type ListCacheValueContract struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// NewListCacheValuesContract is used to add new list cache entry using API.
type NewListCacheValuesContract struct {
	Key    string   `form:"key" json:"key" binding:"required"`
	Values []string `form:"values" json:"values" binding:"required"`
	TTL    string   `form:"ttl" json:"ttl"`
}

// UpdateListCacheValueContract is used to update list cache entry using API.
type UpdateListCacheValueContract struct {
	Value string `form:"value" json:"value" binding:"required"`
}
