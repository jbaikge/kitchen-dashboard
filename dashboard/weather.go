package dashboard

import (
	"fmt"
	"time"
)

type CurrentWeather struct {
	Condition   string
	Temperature int
	Humidity    int
}

type ForecastWeather struct {
	Condition     string
	Low           int
	High          int
	Date          time.Time
	Precipitation float64
}

type Weather struct {
	Current  CurrentWeather
	Forecast []ForecastWeather
}

type WeatherState struct {
	State      string
	Attributes struct {
		Temperature int
		Humidity    int
		Forecast    []struct {
			Condition     string
			Low           int `json:"templow"`
			High          int `json:"temperature"`
			Precipitation float64
			Date          time.Time `json:"datetime"`
		}
	}
}

func GetWeather(config Config) (weather Weather, err error) {
	ha := NewHomeAssistant(config)
	key := fmt.Sprintf("weather.%s", config.Weather.Key)
	state := new(WeatherState)
	if err = ha.GetState(key, state); err != nil {
		err = fmt.Errorf("problem getting weather state: %w", err)
		return
	}

	weather.Current.Condition = state.State
	weather.Current.Temperature = state.Attributes.Temperature
	weather.Current.Humidity = state.Attributes.Humidity
	weather.Forecast = make([]ForecastWeather, len(state.Attributes.Forecast))
	for i, f := range state.Attributes.Forecast {
		weather.Forecast[i].Condition = f.Condition
		weather.Forecast[i].High = f.High
		weather.Forecast[i].Low = f.Low
		weather.Forecast[i].Precipitation = f.Precipitation
		weather.Forecast[i].Date = f.Date
	}

	return
}
