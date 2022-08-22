package dashboard

import (
	"fmt"
	"time"
)

type Sun struct {
	Dawn     time.Time
	Dusk     time.Time
	Midnight time.Time
	Noon     time.Time
	Sunrise  time.Time
	Sunset   time.Time
}

type SunState struct {
	State      string
	Attributes struct {
		NextDawn     time.Time `json:"next_dawn"`
		NextDusk     time.Time `json:"next_dusk"`
		NextMidnight time.Time `json:"next_midnight"`
		NextNoon     time.Time `json:"next_noon"`
		NextRising   time.Time `json:"next_rising"`
		NextSetting  time.Time `json:"next_setting"`
		Rising       bool      `json:"rising"`
	}
}

func GetSun(config Config) (sun Sun, err error) {
	ha := NewHomeAssistant(config)
	key := fmt.Sprintf("sun.%s", config.Sun.Key)
	state := new(SunState)
	if err = ha.GetState(key, state); err != nil {
		return
	}

	tz := config.TimeZone()
	sun.Dawn = state.Attributes.NextDawn.In(tz)
	sun.Dusk = state.Attributes.NextDusk.In(tz)
	sun.Midnight = state.Attributes.NextMidnight.In(tz)
	sun.Noon = state.Attributes.NextNoon.In(tz)
	sun.Sunrise = state.Attributes.NextRising.In(tz)
	sun.Sunset = state.Attributes.NextSetting.In(tz)

	return
}
