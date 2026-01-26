package utils

import (
	"log"
	"time"
)

// GetFromCache intenta obtener un valor del cache y registra hit/miss
func GetFromCache(key string) (interface{}, bool) {
	value, found := AppCache.Get(key)

	if found {
		log.Printf("🟢 CACHE HIT  [%s]", key)
	} else {
		log.Printf("🔴 CACHE MISS [%s]", key)
	}

	return value, found
}

// SetToCache guarda un valor en cache
func SetToCache(key string, value interface{}, ttl time.Duration) {
	AppCache.Set(key, value, ttl)
	log.Printf("🟡 CACHE SET  [%s] ttl=%s", key, ttl.String())
}

// DeleteCache elimina una clave del cache
func DeleteCache(key string) {
	AppCache.Delete(key)
	log.Printf("⚪ CACHE DEL  [%s]", key)
}
