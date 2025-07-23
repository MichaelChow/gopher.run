// Ftoc prints two Fahrenheit-to-Celsius conversions.
// See page 29.
package main

import "fmt"

func main() {
	const freezingF, boilingF = 32.0, 212.0                 // 熔点和沸点
	fmt.Printf("%g°F = %g°C\n", freezingF, fToC(freezingF)) // "32°F = 0°C"
	fmt.Printf("%g°F = %g°C\n", boilingF, fToC(boilingF))   // "212°F = 100°C"

	username := "xiaoming"
	username, password := "xiaoming", "123456"
	fmt.Printf("username: %s, password: %s\n", username, password)
}

func fToC(f float64) float64 {
	return (f - 32) * 5 / 9
}
