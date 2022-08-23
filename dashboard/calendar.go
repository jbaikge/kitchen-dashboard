package dashboard

import (
	"fmt"
	"time"
)

type CalendarEvent struct {
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	PrettyDate  string    `json:"date"`
	PrettyTime  string    `json:"time"`
}

type Calendar struct {
	Title  string          `json:"title"`
	Key    string          `json:"key"`
	Events []CalendarEvent `json:"events"`
}

type CalendarEntry struct {
	Summary     string
	Description string
	Location    string
	Start       struct {
		Date     string
		DateTime string
	}
	End struct {
		Date     string
		DateTime string
	}
}

func (entry CalendarEntry) GetStart(loc *time.Location) (t time.Time) {
	if value := entry.Start.Date; value != "" {
		t, _ = time.ParseInLocation("2006-01-02", value, loc)
		return
	}
	if value := entry.Start.DateTime; value != "" {
		t, _ = time.ParseInLocation(time.RFC3339, value, loc)
	}
	return
}

func (entry CalendarEntry) GetEnd(loc *time.Location) (t time.Time) {
	if value := entry.End.Date; value != "" {
		t, _ = time.ParseInLocation("2006-01-02", value, loc)
		return
	}
	if value := entry.End.DateTime; value != "" {
		t, _ = time.ParseInLocation(time.RFC3339, value, loc)
	}
	return
}

func GetCalendars(config Config) (calendars []Calendar, err error) {
	calendars = make([]Calendar, len(config.Calendars))

	ha := NewHomeAssistant(config)
	tz := config.TimeZone()
	for i, calConfig := range config.Calendars {
		key := fmt.Sprintf("calendar.%s", calConfig.Key)
		start := time.Now()
		end := start.Add(time.Duration(calConfig.Outlook*24) * time.Hour)
		entries := make([]CalendarEntry, 0, 30)
		if err = ha.GetCalendar(key, start, end, &entries); err != nil {
			err = fmt.Errorf("fetching calendar [%s]: %w", calConfig.Key, err)
			return
		}

		events := make([]CalendarEvent, 0, len(entries))
		for _, entry := range entries {
			event := CalendarEvent{
				Summary:     entry.Summary,
				Description: entry.Description,
				Location:    entry.Location,
				Start:       entry.GetStart(tz),
				End:         entry.GetEnd(tz),
			}
			event.PrettyDate = event.Start.Format("Wed, Jan 02")
			if event.End.Sub(event.Start) < 24*time.Hour {
				event.PrettyTime = event.Start.Format("3:04pm - ") + event.End.Format("3:04pm")
			}
			events = append(events, event)
		}

		calendars[i].Title = calConfig.Title
		calendars[i].Key = calConfig.Key
		calendars[i].Events = events
	}

	return
}
