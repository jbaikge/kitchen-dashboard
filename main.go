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
	var err error
	for {
		if data.Lunch, err = dashboard.GetSchoolLunch(config); err != nil {
			log.Printf("error getting school lunch: %v", err)
		}

		now := time.Now()
		nextFetch := time.Date(now.Year(), now.Month(), now.Day(), 8, 55, 0, 0, now.Location())
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

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/data", dataHandler)
	http.Handle("/icons/", http.FileServer(http.Dir(assetDir)))
	http.ListenAndServe(config.Global.Listen, nil)
}
