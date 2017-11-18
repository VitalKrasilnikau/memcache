package contracts

// DictionaryKeyValueContract is a contract which represents key value mapping for dictionary API.
type DictionaryKeyValueContract struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DictionaryCacheValueContract is used to serialize dictionary cache entry via API.
type DictionaryCacheValueContract struct {
	Key    string                       `json:"key"`
	Values []DictionaryKeyValueContract `json:"values"`
}

// NewDictionaryCacheValuesContract is used to add new dictionary cache entry using API.
type NewDictionaryCacheValuesContract struct {
	Key    string                       `form:"key" json:"key" binding:"required"`
	Values []DictionaryKeyValueContract `form:"values" json:"values" binding:"required"`
	TTL    string                       `form:"ttl" json:"ttl"`
}

// UpdateDictionaryCacheValueContract is used to update dictionary cache entry value using API.
type UpdateDictionaryCacheValueContract struct {
	Value    string `form:"value" json:"value" binding:"required"`
	Original string `form:"original" json:"original" binding:"required"`
}

// AddDictionaryCacheValueContract is used to add new dictionary cache entry value using API.
type AddDictionaryCacheValueContract struct {
	Value DictionaryKeyValueContract `form:"value" json:"value" binding:"required"`
}
