package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WeatherData struct {
	Location  string
	Temp      float64
	Condition string
	Humidity  int
	WindSpeed float64
}

func GetWeather(location string) (*WeatherData, error) {
	apiKey:= os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY environment variable not set")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	baseURL:= "https://api.openweathermap.org/data/2.5/weather"
	params := url.Values{}
	params.Add("q", location)
	params.Add("appid", apiKey)
	params.Add("units", "metric")

	resp, err := client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned a non 200 status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err:= json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	weather := &WeatherData{
		Location: location,
	}

	if main, ok := result["main"].(map[string]interface{}); ok {
		if temp, ok := main["temp"].(float64); ok {
			weather.Temp = temp
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
			weather.WindSpeed = speed
		}
	}
	return weather, nil
}
