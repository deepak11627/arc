package arc

// CacheService is interface for ARC
type CacheService interface {
	Get(key interface{}) (value interface{}, ok bool)
	Put(key, value interface{}) bool
	Traverse()
	Len() int
}
