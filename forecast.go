package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
	"math"
)

type ForecastData struct {
	Date				time.Time
	TempMin			int	
	TempMax			int	
	Condition		string
	Humidity		int
	WindSpeed		int	
	Description	string
}

func GetForecast(location string) ([]ForecastData, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY environment variable not set")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	baseURL := "https://api.openweathermap.org/data/2.5/forecast"
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
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	forecastList, ok := result["list"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	forecasts := []ForecastData{}
	currentDay := -1

	for _, item := range forecastList {
		forecast, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		timestamp, ok := forecast["dt"].(float64)
		if !ok {
			continue
		}

		forecastTime := time.Unix(int64(timestamp), 0)

		if currentDay == forecastTime.Day() {
			continue
		}
		currentDay = forecastTime.Day()

		forecastData := ForecastData{
			Date: forecastTime,
		}

		if main, ok := forecast["main"].(map[string]interface{}); ok {
			if temp, ok := main["temp_min"].(float64); ok {
				forecastData.TempMin =	int(math.Round(temp)) 
			}
			if temp, ok := main["temp_max"].(float64); ok {
				forecastData.TempMax = int(math.Round(temp))
			}
			if humidity, ok := main["humidity"].(float64); ok {
				forecastData.Humidity = int(humidity)
			}
		}

		if weatherArray, ok := forecast["weather"].([]interface{}); ok && len(weatherArray) > 0 {
			if weatherMap, ok := weatherArray[0].(map[string]interface{}); ok {
				if main, ok := weatherMap["main"].(string); ok {
					forecastData.Condition = main
				}
				if description, ok := weatherMap["desceiption"].(string); ok {
					forecastData.Description = description	
				}
			}
		}
	
		if wind, ok := forecast ["wind"].(map[string]interface{}); ok {
			if speed, ok := wind["speed"].(float64); ok {
				forecastData.WindSpeed = int(math.Round(speed)) 
			}
		}

		forecasts = append(forecasts, forecastData)

		if len(forecasts) >= 5 {
			break
		}
	}
	
		return forecasts, nil
}
