package main

import (
	"flag"
	"fmt"
 	"os"

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
		os.Exit(1)
	}

	fmt.Printf("Getting weather for: %s\n", *location)
	weather, err := GetWeather(*location)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nWeather for %s:\n", weather.Location)
	fmt.Printf("Temperature: %.1f°C\n", weather.Temp)
	fmt.Printf("Condition: %s\n", weather.Condition)
	fmt.Printf("Humidity: %d%%\n", weather.Humidity)
	fmt.Printf("Wind Speed: %.1f m/s\n", weather.WindSpeed)
}
