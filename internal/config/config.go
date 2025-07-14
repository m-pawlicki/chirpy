package config

import "sync/atomic"

type APIConfig struct {
	FileserverHits atomic.Int32
}
