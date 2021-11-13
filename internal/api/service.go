package api

import (
	"fmt"

	"github.com/ReneKroon/ttlcache"
	"github.com/go-martini/martini"
	"github.com/punxlab/sadwave-events-api-v2/internal/parser"
)

const (
	paramCity = "city"

	cacheKeyEvents = "events-key"
)

type Implementation struct {
	cache *ttlcache.Cache
}

func New(cache *ttlcache.Cache) *Implementation {
	return &Implementation{
		cache: cache,
	}
}

func (i *Implementation) Run() {
	m := martini.Classic()

	m.Get(fmt.Sprintf("/events/:%s", paramCity), i.getEvents)
	m.Get("/cities", i.getCities)

	m.Run()
}

func (i *Implementation) getCachedEvents() (map[parser.CityCode]*parser.CityEvents, error) {
	cache, ok := i.cache.Get(cacheKeyEvents)
	if ok {
		events, ok := cache.(map[parser.CityCode]*parser.CityEvents)
		if ok {
			return events, nil

		}

		i.cache.Purge()
	}

	events, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse events: %s", err)
	}

	i.cache.Set(cacheKeyEvents, events)

	return events, nil
}
