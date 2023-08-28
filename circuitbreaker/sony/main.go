package main

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

func main() {
	// Define the settings for the circuit breaker
	settings := gobreaker.Settings{
		Name: "exampleCircuit",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open the circuit if there have been more than 3 consecutive failures
			return counts.ConsecutiveFailures > 3
		},
		Timeout: 5 * time.Second, // Time window to consider failures
	}

	// Create a new circuit breaker instance with the defined settings
	circuitBreaker := gobreaker.NewCircuitBreaker(settings)

	// Execute the function wrapped by the circuit breaker
	for i := 0; i < 10; i++ {
		result, err := circuitBreaker.Execute(func() (interface{}, error) {
			// Simulate a service call that may fail
			if i < 7 {
				return nil, fmt.Errorf("simulated error")
			}
			return "Success", nil
		})

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Result: %v\n", result)
		}
	}
}
