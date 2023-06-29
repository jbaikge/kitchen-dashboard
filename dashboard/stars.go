package dashboard

import (
	"strconv"
)

type Star struct {
	Name  string `json:"name"`
	Stars int    `json:"stars"`
}

type StarAttributesResponse struct {
	Name string `json:"friendly_name"`
}

type StarResponse struct {
	State      string                 `json:"state"`
	Attributes StarAttributesResponse `json:"attributes"`
}

func GetStars(config Config) (stars []Star, err error) {
	ha := NewHomeAssistant(config)

	stars = make([]Star, len(config.Stars))
	for i, starConfig := range config.Stars {
		response := new(StarResponse)
		if err = ha.GetState(starConfig.Key, response); err != nil {
			return
		}

		stars[i].Name = response.Attributes.Name
		if stars[i].Stars, err = strconv.Atoi(response.State); err != nil {
			continue
		}
	}
	return
}
