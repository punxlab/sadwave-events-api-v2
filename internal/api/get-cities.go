package api

import (
	"encoding/json"
	"fmt"

	"github.com/punxlab/sadwave-events-api-v2/internal/parser"
)

func (i *Implementation) getCities() (int, string) {
	allEvents, err := i.getEventsFromCache()
	if err != nil {
		return 500, fmt.Sprintf("get cached events: %s", err)
	}

	cities := make([]*parser.City, 0, len(allEvents))
	for _, cityEvents := range allEvents {
		cities = append(cities, cityEvents.City)
	}

	data, err := json.Marshal(cities)
	if err != nil {
		return 500, fmt.Sprintf("marshal cities: %s", err)
	}

	return 200, string(data)
}
