package contracts

type DictionaryKeyValueContract struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DictionaryCacheValueContract struct {
	Key    string                       `json:"key"`
	Values []DictionaryKeyValueContract `json:"values"`
}

type NewDictionaryCacheValuesContract struct {
	Key    string                       `form:"key" json:"key" binding:"required"`
	Values []DictionaryKeyValueContract `form:"values" json:"values" binding:"required"`
	TTL    string                       `form:"ttl" json:"ttl"`
}

type UpdateDictionaryCacheValueContract struct {
	Value    string `form:"value" json:"value" binding:"required"`
	Original string `form:"original" json:"original" binding:"required"`
}

type AddDictionaryCacheValueContract struct {
	Value DictionaryKeyValueContract `form:"value" json:"value" binding:"required"`
}
