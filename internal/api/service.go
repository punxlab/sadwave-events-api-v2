package api

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gadelkareem/cachita"
	"github.com/go-martini/martini"
	"github.com/punxlab/sadwave-events-api-v2/internal/parser"
)

const (
	paramCity = "city"

	cacheKeyEvents = "events-key"
)

type Implementation struct {
	mu    sync.RWMutex
	cache cachita.Cache
}

func New(cache cachita.Cache) *Implementation {
	return &Implementation{
		cache: cache,
	}
}

func (i *Implementation) Run() {
	m := martini.Classic()

	i.watchEvents()

	m.Get(fmt.Sprintf("/events/:%s", paramCity), i.getEvents)
	m.Get("/cities", i.getCities)

	m.RunOnAddr(":80")
}

func (i *Implementation) getEventsFromCache() (map[parser.CityCode]*parser.CityEvents, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	events := make(map[parser.CityCode]*parser.CityEvents, 0)
	err := i.cache.Get(cacheKeyEvents, &events)
	if err != nil && err != cachita.ErrNotFound {
		return nil, fmt.Errorf("get cache: %s", err)
	}

	return events, nil
}

func (i *Implementation) setEventsToCache(events map[parser.CityCode]*parser.CityEvents) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	err := i.cache.Put(cacheKeyEvents, events, -1)
	if err != nil && err != cachita.ErrNotFound {
		return fmt.Errorf("put cache: %s", err)
	}

	return nil
}

func (i *Implementation) watchEvents() {
	i.updateEvents()
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			i.updateEvents()
		}
	}()
}

func (i *Implementation) updateEvents() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	events, err := parser.Parse()
	if err != nil {
		log.Println(err)
	}

	if len(events) > 0 {
		err = i.setEventsToCache(events)
		if err != nil {
			log.Println(err)
		}
	}
}
