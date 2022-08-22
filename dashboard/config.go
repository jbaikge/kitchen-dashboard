package dashboard

import (
	"github.com/BurntSushi/toml"
)

type HomeAssistantConfig struct {
	Endpoint string `toml:"endpoint"`
	Token    string `toml:"token"`
}

type SchoolBusConfig struct {
	Endpoint string `toml:"endpoint"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type CalendarConfig struct {
	Title string `toml:"title"`
	Key   string `toml:"key"`
}

type LockConfig struct {
	Title string `toml:"title"`
	Key   string `toml:"key"`
}

type SunConfig struct {
	Key string `toml:"key"`
}

type WeatherConfig struct {
	Key string `toml:"key"`
}

type Config struct {
	HomeAssistant HomeAssistantConfig `toml:"home-assistant"`
	SchoolBus     SchoolBusConfig     `toml:"school-bus"`
	Calendars     []CalendarConfig    `toml:"calendars"`
	Locks         []LockConfig        `toml:"locks"`
	Sun           SunConfig           `toml:"sun"`
	Weather       WeatherConfig       `toml:"weather"`
}

func ParseConfig(filename string) (config Config, err error) {
	_, err = toml.DecodeFile(filename, &config)
	return
}
