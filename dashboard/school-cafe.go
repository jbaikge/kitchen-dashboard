package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	SchoolCafeEntree    = "ENTREE"
	SchoolCafeDaily     = "DAILY OFFERINGS"
	SchoolCafeGrain     = "GRAIN"
	SchoolCafeVegetable = "VEGETABLE"
	SchoolCafeFruit     = "FRUIT"
	SchoolCafeMilk      = "MILK"
)

const (
	SchoolCafeMealTypeLunch = "Lunch"
)

type SchoolCafeItem struct {
	Description string `json:"MenuItemDescription"`
}

type SchoolCafe struct {
	config Config
}

func NewSchoolCafe(config Config) SchoolCafe {
	return SchoolCafe{
		config: config,
	}
}

func (sc SchoolCafe) GetMenuItems(date time.Time, mealType string) (items map[string][]SchoolCafeItem, err error) {
	base, err := url.Parse(sc.config.SchoolLunch.Endpoint)
	if err != nil {
		err = fmt.Errorf("parsing base URL: %w", err)
		return
	}

	target := base.JoinPath("CalendarView", "GetDailyMenuitemsByGrade")

	query := url.Values{}
	query.Set("SchoolId", sc.config.SchoolLunch.SchoolId)
	query.Set("ServingDate", date.Format("01/02/2006"))
	query.Set("ServingLine", "Main Line")
	query.Set("MealType", mealType)
	query.Set("Grade", fmt.Sprintf("%02d", sc.config.SchoolLunch.Grade))
	query.Set("PersonId", "null")
	target.RawQuery = query.Encode()

	log.Printf("%s", target.String())
	request, err := http.NewRequest(http.MethodGet, target.String(), nil)
	if err != nil {
		err = fmt.Errorf("creating request: %w", err)
		return
	}

	// This may or may not fix the timeout issues when contacting the API
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	http.DefaultClient.Timeout = 10 * time.Second
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		err = fmt.Errorf("performing request: %w", err)
		return
	}
	defer response.Body.Close()

	items = make(map[string][]SchoolCafeItem)
	if err = json.NewDecoder(response.Body).Decode(&items); err != nil {
		return
	}

	return
}
