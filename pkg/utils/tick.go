package utils

import (
	"context"
	"time"
)

func LoopTick(ctx context.Context, t time.Duration, action func() error) error {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		// Run action function
		err := action()
		if err != nil {
			return err
		}

		select {
		// Return if context is cancelled
		case <-ctx.Done():
			return ctx.Err()
		// Break select every tick
		case <-ticker.C:
		}
	}
}
