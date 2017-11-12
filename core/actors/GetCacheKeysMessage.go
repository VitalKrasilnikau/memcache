package act

// GetCacheKeysMessage is used to request all cache keys.
type GetCacheKeysMessage struct{}
// GetCacheKeysReply is a reply message for GetCacheKeysMessage.
type GetCacheKeysReply struct {
	Keys []string
}