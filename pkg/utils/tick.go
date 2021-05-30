/*
Copyright © 2021 Bhargav SNV bhargavsnv100@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"context"
	"time"
)

// LoopTick, runs a given action in a loop in periods of 't' duration. It exits
// when the context is cancelled
func LoopTick(ctx context.Context, t time.Duration, action func() error) error {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		// Run action function
		// err := action()
		// if err != nil {
		// 	return err
		// }
		go action()

		select {
		// Return if context is cancelled
		case <-ctx.Done():
			return ctx.Err()
		// Break select every tick
		case <-ticker.C:
		}
	}
}
