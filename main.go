package main

import (
	"log"
	"time"

	"github.com/gadelkareem/cachita"
	"github.com/punxlab/sadwave-events-api-v2/internal/api"
)

func main() {
	cache, err := cachita.NewFileCache("bin/cache", 30*24*time.Hour, 0)
	if err != nil {
		log.Panicln(err)
	}

	i := api.New(cache)
	i.Run()
}
