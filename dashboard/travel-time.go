package dashboard

import (
	"fmt"
	"log"
)

type TravelTime struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type TravelState struct {
	State      string `json:"state"`
	Attributes struct {
		Duration float64 `json:"duration"`
		Distance float64 `json:"distance"`
		Unit     string  `json:"unit_of_measurement"`
		Name     string  `json:"friendly_name"`
	} `json:"attributes"`
}

func GetTravelTimes(config Config) (travelTimes []TravelTime, err error) {
	travelTimes = make([]TravelTime, len(config.Travel))
	ha := NewHomeAssistant(config)

	for i, travelConfig := range config.Travel {
		key := fmt.Sprintf("sensor.%s", travelConfig.Key)
		state := new(TravelState)
		if err = ha.GetState(key, state); err != nil {
			return
		}

		log.Printf("%+v", state)
		travelTimes[i].Name = state.Attributes.Name
		travelTimes[i].Duration = fmt.Sprintf("%s%s", state.State, state.Attributes.Unit)
	}
	return
}
