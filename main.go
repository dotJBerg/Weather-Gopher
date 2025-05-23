package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	location := flag.String("location", "", "Location to get weather for")
	forecast := flag.Bool("forecast", false, "Show 5-day forecast")
	simple := flag.Bool("simple", false, "Use simple CLI output instead of TUI")
	flag.Parse()

	if *location == "" {
		fmt.Println("Please provide a location using the -location flag")
		fmt.Println("Example: ./weather-gopher -location \"New York\"")
		os.Exit(1)
	}
	if *simple {
		if *forecast {
			fmt.Printf("Getting 5-day forecast for: %s\n", *location)
			forecastData, err := GetForecast(*location)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			displayForecast(forecastData)
		} else {
			fmt.Printf("Getting weather for: %s\n", *location)
			weather, err := GetWeather(*location)
			if err != nil {
				if strings.Contains(err.Error(), "OPENWEATHER_API_KEY") {
					fmt.Println("Error: API key not found. Please set the OPENWEATHER_API_KEY environment variable.")
					fmt.Println("You can get a free API key from https://openweathermap.org/")
				} else if strings.Contains(err.Error(), "no such host") {
					fmt.Println("Error: Could not connect to the weather service. Please check your internet connection.")
				} else {
					fmt.Printf("Error: %v\n", err)
				}

				os.Exit(1)
			}

			displayWeather(weather)
		}
	} else {
		if err := StartUI(*location, *forecast); err != nil {
			fmt.Printf("Error running UI: %v\n", err)
			os.Exit(1)
		}
	}
}

func displayWeather(weather *WeatherData) {
	fmt.Println("\n=================================")
	fmt.Printf("  Weather for %s\n", weather.Location)
	fmt.Println("=================================")

	fmt.Printf("Temperature: %d°F\n", weather.Temp)
	fmt.Printf("Condition:   %s\n", weather.Condition)
	fmt.Printf("Humidity:    %d%%\n", weather.Humidity)
	fmt.Printf("Wind Speed:  %dmph\n", weather.WindSpeed)
	fmt.Println("=================================")
}

func displayForecast(forecasts []ForecastData) {
	fmt.Println("\n=================================")
	fmt.Println("       5-DAY FORECAST")
	fmt.Println("=================================")

	for _, day := range forecasts {
		weekday := day.Date.Format("Monday")

		fmt.Printf("\n%s (%s):\n", weekday, day.Date.Format("Jan 2"))
		fmt.Printf("  Conditions: %s\n", day.Description)
		fmt.Printf("  Temperature: %d°F to %d°F\n", day.TempMin, day.TempMax)
		fmt.Printf("  Humidity: %d%%\n", day.Humidity)
		fmt.Printf("  Wind: %d mph\n", day.WindSpeed)
		fmt.Println("  ---------------------------------")
	}
}
