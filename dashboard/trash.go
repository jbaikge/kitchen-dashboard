package dashboard

import (
	"time"
)

type Trash struct {
	PutCansOut bool `json:"putCansOut"`
}

func GetTrash(config Config) (trash Trash) {
	now := time.Now()

	// "After noon today, notify of trash pickup tomorrow"
	if now.Add(12*time.Hour).Weekday() == config.Trash.Day {
		trash.PutCansOut = true
		return
	}

	// "Before noon today, trash pickup is today"
	if now.Weekday() == config.Trash.Day {
		noonToday := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
		trash.PutCansOut = now.Before(noonToday)
	}

	return
}
