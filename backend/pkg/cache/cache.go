package cache

var data map[string]any = make(map[string]any)

// Get from cache get
func Get(key string) any {
	return data[key]
}

// Set set to cache
func Set(key string, value any) {
	data[key] = value
}
