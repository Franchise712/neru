package main

import (
	"testing"
	"time"
)

func TestWaitForAccessibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		trueAfter  int // check reports granted from call trueAfter+1 on; -1 never
		want       bool
		wantChecks int
	}{
		{name: "granted immediately", trueAfter: 0, want: true, wantChecks: 1},
		{name: "transient denial at login", trueAfter: 3, want: true, wantChecks: 4},
		{name: "really not granted", trueAfter: -1, want: false, wantChecks: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			calls := 0
			check := func() bool {
				calls++

				return tt.trueAfter >= 0 && calls > tt.trueAfter
			}

			var slept time.Duration

			got := waitForAccessibility(check, time.Second, 250*time.Millisecond,
				func(d time.Duration) { slept += d })

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}

			if calls != tt.wantChecks {
				t.Fatalf("checked %d times, want %d", calls, tt.wantChecks)
			}

			if slept > time.Second {
				t.Fatalf("slept %v, past the 1s grace", slept)
			}
		})
	}
}
