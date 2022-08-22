package dashboard

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type GlobalConfig struct {
	TimeZone string `toml:"timezone"`
}

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
	Global        GlobalConfig        `toml:"global"`
	HomeAssistant HomeAssistantConfig `toml:"home-assistant"`
	SchoolBus     SchoolBusConfig     `toml:"school-bus"`
	Calendars     []CalendarConfig    `toml:"calendars"`
	Locks         []LockConfig        `toml:"locks"`
	Sun           SunConfig           `toml:"sun"`
	Weather       WeatherConfig       `toml:"weather"`
}

func (c Config) TimeZone() (tz *time.Location) {
	tz, err := time.LoadLocation(c.Global.TimeZone)
	if err != nil {
		log.Printf("invalid timezone value, `%s`, using local", c.Global.TimeZone)
		tz = time.Local
	}
	return
}

func ParseConfig(filename string) (config Config, err error) {
	config.Global.TimeZone = "Local"
	config.Sun.Key = "sun"

	_, err = toml.DecodeFile(filename, &config)
	return
}
