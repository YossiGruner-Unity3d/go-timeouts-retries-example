package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/avast/retry-go"
)

func main() {
	// Define the maximum timeout duration
	timeout := 5 * time.Second

	// Create a context with the defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Define retry options with jitter
	retryOptions := []retry.Option{
		retry.Attempts(3),
		retry.Delay(2 * time.Second),
		retry.LastErrorOnly(true),
		retry.MaxJitter(500 * time.Millisecond), // Maximum jitter duration
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("Retry %d: %s\n", n, err)
		}),
	}

	// Call the function with timeout and retries using retry-go
	err := retry.Do(
		func() error {
			return callWithRetry(ctx, performTask)
		},
		retryOptions...)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Task completed successfully!")
	}
}

func performTask() error {
	// Simulate a task with a chance of failing
	if rand.Intn(10) < 8 {
		return fmt.Errorf("task failed")
	}
	return nil
}

func callWithRetry(ctx context.Context, fn func() error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fn()
	}
}
