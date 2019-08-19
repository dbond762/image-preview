package timer

import (
	"context"
	"log"
	"time"
)

func Start(timeout time.Duration, resetTimeout <-chan struct{}, timeIsOver chan<- struct{}) {
	stopTimer := make(chan struct{})

	for {
		timerCtx, cancel := context.WithTimeout(context.Background(), timeout)

		go func(ctx context.Context, stopTimer chan<- struct{}) {
			<-ctx.Done()
			stopTimer <- struct{}{}
		}(timerCtx, stopTimer)

		select {
		case <-resetTimeout:
			log.Print("reset")
			cancel()
			continue
		case <-stopTimer:
			timeIsOver <- struct{}{}
			break
		}
	}
}
