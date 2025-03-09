package log

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type spinner struct {
	frames []string
	speed  time.Duration
}

var (
	hourglass spinner = spinner{
		frames: []string{"⏳", "⌛"},
		speed:  500 * time.Millisecond,
	}
	dots spinner = spinner{
		frames: []string{"⠴", "⠦", "⠧", "⠇", "⠏", "⠋", "⠙", "⠹", "⠸", "⠼"},
		speed:  100 * time.Millisecond,
	}
	pulse spinner = spinner{
		frames: []string{"⬫", "⬨", "◊", "⬨"},
		speed:  200 * time.Millisecond,
	}
	runner spinner = spinner{
		frames: []string{"🏃", "🚶"},
		speed:  200 * time.Millisecond,
	}
	locking spinner = spinner{
		frames: []string{"🔓", "🔓", "🔓", "🔓", "🔒"},
		speed:  200 * time.Millisecond,
	}
	unlocking spinner = spinner{
		frames: []string{"🔒", "🔒", "🔒", "🔒", "🔓"},
		speed:  200 * time.Millisecond,
	}
	monkies spinner = spinner{
		frames: []string{"🙉", "🙈", "🙊"},
		speed:  500 * time.Millisecond,
	}
)

// WaitSpinner creates a new spinner, providing a wait group and spinner cancel function.
// Use it to provide a visual indication of a long running process.
//
//	ctx := context.Background()
//	// start a spinner
//	wg, cancel := WaitSpinner(, "waiting on foo", runner)
//	defer wg.Wait()
//	_, err := doSomething()
//	// stop the spinner with err values (nil means success)
//	cancel(err)
//	// then handle error as normal
//	if err != nil {
//		// handle error
//	}
func WaitSpinner(ctx context.Context, message string, s spinner) (*sync.WaitGroup, context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go frame(ctx, &wg, message, s)
	return &wg, cancel
}

// frame prints spinner frames until context is cancelled,
// then prints a final message depending on cancellation cause
func frame(ctx context.Context, wg *sync.WaitGroup, message string, s spinner) {
	// we require this wait group to be done so we can signal the main function to exit
	defer wg.Done()
	for {
		for _, frame := range s.frames {
			select {
			case <-ctx.Done():
				cause := context.Cause(ctx)
				// fmt.Println("cause:", cause)
				switch cause {
				case context.Canceled:
					fmt.Printf("\r%s %s ... done!\n", IconInfo, message)
				default:
					fmt.Printf("\r%s %s ... failed: %v\n", IconError, message, cause)
				}
				return
			// print the next spinner frame
			default:
				fmt.Printf("\r%s %s ... ", frame, message)
				time.Sleep(s.speed)
			}
		}
	}
}
