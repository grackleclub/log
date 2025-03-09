package log

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestFrame(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	go frame(ctx, &wg, "waiting on foo", runner)

	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()
}

func TestWaitSpinner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wg, cancel := WaitSpinner(context.Background(), "waiting on foo", runner)
		defer wg.Wait()
		time.Sleep(1 * time.Nanosecond)
		cancel(nil)
	})
	t.Run("failure", func(t *testing.T) {
		wg, cancel := WaitSpinner(context.Background(), "waiting on foo", runner)
		defer wg.Wait()
		time.Sleep(1 * time.Second)
		cancel(errors.New("tired of waiting"))
	})
}

func TestShowExample(t *testing.T) {
	ctx := context.Background()
	wg, cancel := WaitSpinner(ctx, "reticulating splines", runner)
	defer wg.Wait()
	time.Sleep(2 * time.Second)
	cancel(nil)

	wg, cancel = WaitSpinner(ctx, "fleebing florbs", runner)
	defer wg.Wait()
	time.Sleep(2 * time.Second)
	cancel(nil)

	wg, cancel = WaitSpinner(ctx, "dismantling capitalism", runner)
	defer wg.Wait()
	time.Sleep(4 * time.Second)
	cancel(errors.New("too entrenched, build community and try again"))

	time.Sleep(2 * time.Second)
}
