package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// AppCache es el caché global de la aplicación (in-memory)
var AppCache = cache.New(
	5*time.Minute,  // expiración por defecto
	10*time.Minute, // limpieza automática
)
