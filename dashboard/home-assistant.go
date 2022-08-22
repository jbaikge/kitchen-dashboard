package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type HomeAssistant struct {
	config Config
}

func NewHomeAssistant(config Config) HomeAssistant {
	return HomeAssistant{
		config: config,
	}
}

// Fetches a state by the provided key and decodes it into the state argument
func (ha HomeAssistant) GetState(key string, state interface{}) (err error) {
	target, err := url.JoinPath(ha.config.HomeAssistant.Endpoint, "states", key)
	fmt.Println(target)
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
