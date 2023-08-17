package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Define the maximum timeout duration
	timeout := 25 * time.Second

	// Create a context with the defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Call the function with timeout and retries
	err := callWithTimeoutAndRetry(ctx, 15, 2*time.Second, performTask)
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

func callWithTimeoutAndRetry(ctx context.Context, maxRetries int, retryInterval time.Duration, fn func() error) error {
	for retries := 0; retries <= maxRetries; retries++ {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled:", ctx.Err())
			return ctx.Err()
		default:
			if err := fn(); err != nil {
				fmt.Printf("Retry %d: %s\n", retries, err)
				sleepDuration := retryInterval * time.Duration(1<<uint(retries))
				fmt.Printf("Sleeping for %s\n", sleepDuration)
				time.Sleep(sleepDuration)
				fmt.Println("Retrying...")
				continue
			}
			fmt.Println("Task completed successfully!")
			return nil
		}
	}
	fmt.Println("Exhausted all retries")
	return fmt.Errorf("exhausted all retries")
}
