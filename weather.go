package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WeatherData struct {
	Location  string
	Temp      int
	Condition string
	Humidity  int
	WindSpeed int
}

func GetWeather(location string) (*WeatherData, error) {
	if cachedData, found := getWeatherFromCache(location); found {
		fmt.Println("Using cached weather data")
		return cachedData, nil
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY environment variable not set")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	baseURL := "https://api.openweathermap.org/data/2.5/weather"
	params := url.Values{}
	params.Add("q", location)
	params.Add("appid", apiKey)
	params.Add("units", "imperial")

	resp, err := client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned a non 200 status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	weather := &WeatherData{
		Location: location,
	}

	if main, ok := result["main"].(map[string]interface{}); ok {
		if temp, ok := main["temp"].(float64); ok {
			weather.Temp = int(math.Round(temp))
		}
		if humidity, ok := main["humidity"].(float64); ok {
			weather.Humidity = int(humidity)
		}
	}

	if weatherArray, ok := result["weather"].([]interface{}); ok && len(weatherArray) > 0 {
		if weatherMap, ok := weatherArray[0].(map[string]interface{}); ok {
			if description, ok := weatherMap["description"].(string); ok {
				weather.Condition = description
			}
		}
	}

	if wind, ok := result["wind"].(map[string]interface{}); ok {
		if speed, ok := wind["speed"].(float64); ok {
			weather.WindSpeed = int(math.Round(speed))
		}
	}

	if err := saveWeatherToCache(location, weather); err != nil {
		fmt.Printf("Warning: Failed to cach weather data: %v\n", err)
	}

	return weather, nil
}
