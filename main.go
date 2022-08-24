package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jbaikge/kitchen-dashboard/dashboard"
)

var assetDir = "./assets"
var config dashboard.Config
var configFile = "config.toml"
var data = struct {
	Calendars []dashboard.Calendar `json:"calendars"`
	Locks     []dashboard.Lock     `json:"locks"`
	Lunch     dashboard.SchoolMenu `json:"lunch"`
	Weather   dashboard.Weather    `json:"weather"`
	Sun       dashboard.Sun        `json:"sun"`
	Trash     dashboard.Trash      `json:"trash"`
}{}

func init() {
	flag.StringVar(&assetDir, "assets", assetDir, "Path to assets directory")
	flag.StringVar(&configFile, "config", configFile, "Config file location (TOML)")
}

func updateCalendars() {
	ticker := time.Tick(5 * time.Minute)
	var err error
	for {
		if data.Calendars, err = dashboard.GetCalendars(config); err != nil {
			log.Printf("error getting calendar information: %v", err)
		}
		<-ticker
	}
}

func updateLocks() {
	ticker := time.Tick(time.Minute)
	var err error
	for {
		if data.Locks, err = dashboard.GetLocks(config); err != nil {
			log.Printf("error getting lock information: %v", err)
		}
		<-ticker
	}
}

func updateSchoolLunch() {
	hour := 9
	var err error
	for {
		date := time.Now()
		if date.After(time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())) {
			date = date.Add(24 * time.Hour)
		}

		if data.Lunch, err = dashboard.GetSchoolLunch(config, date); err != nil {
			log.Printf("error getting school lunch: %v", err)
		}

		now := time.Now()
		nextFetch := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
		if now.After(nextFetch) {
			nextFetch = nextFetch.Add(24 * time.Hour)
		}
		duration := nextFetch.Sub(now).Abs()
		<-time.After(duration)
	}
}

func updateSun() {
	ticker := time.Tick(15 * time.Minute)
	var err error
	for {
		if data.Sun, err = dashboard.GetSun(config); err != nil {
			log.Printf("error getting sun information: %v", err)
		}
		<-ticker
	}
}

func updateTrash() {
	ticker := time.Tick(time.Minute)
	for {
		data.Trash = dashboard.GetTrash(config)
		<-ticker
	}
}

func updateWeather() {
	ticker := time.Tick(15 * time.Minute)
	var err error
	for {
		if data.Weather, err = dashboard.GetWeather(config); err != nil {
			log.Printf("error getting weather information: %v", err)
		}
		<-ticker
	}
}

func main() {
	flag.Parse()

	var err error
	config, err = dashboard.ParseConfig(configFile)
	if err != nil {
		log.Fatalf("error parsing config: %s", err)
	}

	go updateCalendars()
	go updateLocks()
	go updateTrash()
	go updateSun()
	go updateWeather()
	go updateSchoolLunch()

	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		path := filepath.Join(assetDir, "index.html")
		f, err := os.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(w, f)
	}

	dataHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}

	lockHandler := func(w http.ResponseWriter, req *http.Request) {
		status := struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		}{
			Success: true,
		}
		defer json.NewEncoder(w).Encode(&status)

		data := struct {
			Key string `json:"key"`
		}{}
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			status.Success = false
			status.Message = err.Error()
			return
		}

		if err := dashboard.ToggleLock(config, data.Key); err != nil {
			status.Success = false
			status.Message = err.Error()
			return
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/toggle-lock", lockHandler)
	http.Handle("/icons/", http.FileServer(http.Dir(assetDir)))
	http.ListenAndServe(config.Global.Listen, nil)
}
