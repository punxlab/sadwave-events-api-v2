package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gadelkareem/cachita"
	"github.com/punxlab/sadwave-events-api-v2/internal/api"
)

func main() {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panicln(err)
	}
	path = filepath.Join(path, "tmp/file-cache")
	cache, err := cachita.NewFileCache(path, -1, 0)

	i := api.New(cache)
	i.Run()
}
