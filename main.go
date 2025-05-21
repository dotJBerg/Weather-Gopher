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
	flag.Parse()

	if *location == "" {
		fmt.Println("Please provide a location using the -location flag")
		fmt.Println("Example: ./weather-gopher -location \"New York\"")
		os.Exit(1)
	}

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

func displayWeather(weather *WeatherData) {
	fmt.Println("\n=================================")
	fmt.Printf("  Weather for %s\n", weather.Location)
	fmt.Println("=================================")
	
	fmt.Printf("Temperature: %.1f°F\n", weather.Temp)
	fmt.Printf("Condition:   %s\n", weather.Condition)
	fmt.Printf("Humidity:    %d%%\n", weather.Humidity)
	fmt.Printf("Wind Speed:  %.1f mph\n", weather.WindSpeed)
	fmt.Println("=================================")
}
