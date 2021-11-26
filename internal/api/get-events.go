package api

import (
	"encoding/json"
	"fmt"

	"github.com/go-martini/martini"
	"github.com/punxlab/sadwave-events-api-v2/internal/parser"
)

func (i *Implementation) getEvents(params martini.Params) (int, string) {
	city := params[paramCity]

	allEvents, err := i.getEventsFromCache()
	if err != nil {
		return 500, fmt.Sprintf("get cached events: %s", err)
	}

	code := parser.CityCode(city)
	events, ok := allEvents[code]
	if !ok {
		return 404, err.Error()
	}

	data, err := json.Marshal(events.Events)
	if err != nil {
		return 500, fmt.Sprintf("marshal events: %s", err)
	}

	return 200, string(data)
}
