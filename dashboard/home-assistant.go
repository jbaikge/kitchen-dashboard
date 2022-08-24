package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HomeAssistant struct {
	config Config
}

func NewHomeAssistant(config Config) HomeAssistant {
	return HomeAssistant{
		config: config,
	}
}

// Fetches a calendar by the provided key and decodes the array into the entries argument
func (ha HomeAssistant) GetCalendar(key string, start time.Time, end time.Time, entries interface{}) (err error) {
	base, err := url.Parse(ha.config.HomeAssistant.Endpoint)
	if err != nil {
		return fmt.Errorf("parsing endpoint: %w", err)
	}

	target := base.JoinPath("calendars", key)

	query := url.Values{}
	query.Add("start", start.Format(time.RFC3339))
	query.Add("end", end.Format(time.RFC3339))
	target.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, target.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	request.Header.Add("Authorization", "Bearer "+ha.config.HomeAssistant.Token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad return: [%d] %s", response.StatusCode, response.Status)
	}

	return json.NewDecoder(response.Body).Decode(entries)
}

// Fetches a state by the provided key and decodes it into the state argument
func (ha HomeAssistant) GetState(key string, state interface{}) (err error) {
	target, err := url.JoinPath(ha.config.HomeAssistant.Endpoint, "states", key)
	if err != nil {
		return fmt.Errorf("creating URL: %w", err)
	}
	request, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	request.Header.Add("Authorization", "Bearer "+ha.config.HomeAssistant.Token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad return: [%d] %s", response.StatusCode, response.Status)
	}

	return json.NewDecoder(response.Body).Decode(state)
}

func (ha HomeAssistant) CallService(domain string, service string, payload interface{}) (err error) {
	base, err := url.Parse(ha.config.HomeAssistant.Endpoint)
	if err != nil {
		return fmt.Errorf("parsing endpoint: %w", err)
	}

	target := base.JoinPath("services", domain, service)

	buffer := new(bytes.Buffer)
	if err = json.NewEncoder(buffer).Encode(payload); err != nil {
		return fmt.Errorf("encoding payload: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, target.String(), buffer)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+ha.config.HomeAssistant.Token)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad return: [%d] %s", response.StatusCode, response.Status)
	}

	return
}
