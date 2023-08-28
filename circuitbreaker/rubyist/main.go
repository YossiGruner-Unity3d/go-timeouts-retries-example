package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go"
	circuit "github.com/rubyist/circuitbreaker"
)

func main() {
	// Create a new circuit breaker with default settings
	cb := circuit.NewThresholdBreaker(3)

	// Define the maximum timeout duration
	timeout := 5 * time.Second

	// Create a context with the defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Define retry options
	retryOptions := []retry.Option{
		retry.Attempts(3),
		retry.Delay(2 * time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(logRetry), // Add log for retries
		retry.LastErrorOnly(true),
	}

	// Simulate invoking a function with circuit breaker, timeout, and retry
	err := cb.CallContext(ctx, func() error {
		return retry.Do(
			func() error {
				return performTask()
			},
			retryOptions...,
		)
	}, timeout)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success")
	}
}

func performTask() error {
	// Simulate a task with a chance of failing
	if time.Now().UnixNano()%2 == 0 {
		return fmt.Errorf("task failed")
	}
	return nil
}

func logRetry(n uint, err error) {
	log.Printf("Retry attempt %d failed: %v", n, err)
}
