package main

import (
	"context"
	"fmt"
	"github.com/anna-fry/reinforcer/example/client"
	"github.com/anna-fry/reinforcer/example/client/reinforced"
	"github.com/anna-fry/reinforcer/pkg/runner"
	"github.com/slok/goresilience/retry"
	"github.com/slok/goresilience/timeout"
	"time"
)

func main() {
	cl := client.NewClient()
	f := runner.NewFactory(
		timeout.NewMiddleware(timeout.Config{Timeout: 100 * time.Millisecond}),
		retry.NewMiddleware(retry.Config{
			Times: 10,
		}),
	)
	rCl := reinforced.NewClient(cl, f, reinforced.WithRetryableErrorPredicate(func(s string, err error) bool {
		// Always retry SayHello, don't retry any other error
		return s == reinforced.ClientMethods.SayHello
	}))
	for i := 0; i < 100; i++ {
		err := rCl.SayHello(context.Background(), "Christian")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
