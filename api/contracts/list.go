package contracts

type ListCacheValueContract struct {
	Key    string `json:"key"`
	Values []string `json:"values"`
}

type NewListCacheValuesContract struct {
	Key    string `form:"key" json:"key" binding:"required"`
	Values []string `form:"values" json:"values" binding:"required"`
	TTL    string `form:"ttl" json:"ttl"`
}

type UpdateListCacheValueContract struct {
	Value      string `form:"value" json:"value" binding:"required"`
}
