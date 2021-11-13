package main

import (
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/punxlab/sadwave-events-api-v2/internal/api"
)

func main() {
	cache := ttlcache.NewCache()
	cache.SetTTL(5 * time.Minute)

	i := api.New(cache)
	i.Run()
}
