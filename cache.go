package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CachedWeather struct {
	Data      *WeatherData
	Timestamp time.Time
}

func getCachePath(location string) string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// Note: this is supposed to use some temp directory. in case there is an error with the cache file there should be a backup.
		cacheDir = os.TempDir()
	}

	weatherDir := filepath.Join(cacheDir, "weather-gopher")
	os.MkdirAll(weatherDir, 0755)

	filename := fmt.Sprintf("%s.json", location)
	return filepath.Join(weatherDir, filename)
}

func getWeatherFromCache(location string) (*WeatherData, bool) {
	cachePath := getCachePath(location)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, false
	}

	var cached CachedWeather
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, false
	}

	if time.Since(cached.Timestamp) > 30*time.Minute {
		return nil, false
	}

	return cached.Data, true
}

func saveWeatherToCache(location string, data *WeatherData) error {
	cachePath := getCachePath(location)

	cached := CachedWeather{
		Data:      data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, jsonData, 0644)
}
