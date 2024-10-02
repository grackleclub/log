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

func frame(ctx context.Context, wg *sync.WaitGroup, message string, s spinner) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\r%s %s done!\n", iconInfo, message)
			return
		default:
			for _, frame := range s.frames {
				fmt.Printf("\r%s %s ...", frame, message)
				time.Sleep(s.speed)
			}
		}
	}
}
