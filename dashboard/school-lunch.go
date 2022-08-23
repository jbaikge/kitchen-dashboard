package dashboard

import (
	"fmt"
	"time"
)

type SchoolMenuItem struct {
	Description string `json:"description"`
}

type SchoolMenu struct {
	Entrees []SchoolMenuItem `json:"entrees"`
}

func GetSchoolLunch(config Config, date time.Time) (menu SchoolMenu, err error) {
	cafe := NewSchoolCafe(config)
	entries, err := cafe.GetMenuItems(date, SchoolCafeMealTypeLunch)
	if err != nil {
		return
	}

	entrees, ok := entries[SchoolCafeEntree]
	if !ok {
		err = fmt.Errorf("menu items of type %s not found", SchoolCafeEntree)
	}

	menu.Entrees = make([]SchoolMenuItem, len(entrees))
	for i, entree := range entrees {
		menu.Entrees[i].Description = entree.Description
	}

	return
}
