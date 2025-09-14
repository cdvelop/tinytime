package tinytime

import (
	"fmt"
	"syscall/js"
	"testing"
	"time"

	"github.com/cdvelop/tinystring"
)

// This test will be executed in a browser environment by wasmbrowsertest.
// The `fmt.Println` statements will be logged to the browser's console.
func TestWasm(t *testing.T) {
	// We need to keep the test running until the async operations are complete.
	// We will use a channel for this.
	done := make(chan struct{})

	// The test logic will run in a goroutine to avoid blocking the main thread.
	go func() {
		defer close(done)

		provider := NewTimeProvider()

		// Test UnixNano
		nano := provider.UnixNano()
		if nano <= 0 {
			t.Errorf("UnixNano() returned a non-positive value: %d", nano)
		}
		fmt.Println("UnixNano:", nano)

		// Test UnixSecondsToDate
		// January 1, 2025 00:00:00 UTC
		seconds := int64(1735689600)
		expectedDate := "2025-01-01 00:00"
		dateStr := provider.UnixSecondsToDate(seconds)
		if dateStr != expectedDate {
			t.Errorf("UnixSecondsToDate() = %q, want %q", dateStr, expectedDate)
		}
		fmt.Println("UnixSecondsToDate:", dateStr)

		// Test UnixNanoToTime
		timeStr := provider.UnixNanoToTime(nano)
		fmt.Println("UnixNanoToTime:", timeStr)

		// We need to parse the time string manually because the `tinystring.Fmt`
		// function does not format with leading zeros correctly.
		parts := tinystring.Convert(timeStr).Split(":")
		if len(parts) != 3 {
			t.Fatalf("UnixNanoToTime() returned an invalid time string format: %q", timeStr)
		}

		hour, err := tinystring.Convert(parts[0]).TrimSpace().Int()
		if err != nil {
			t.Fatalf("Failed to parse hour: %v", err)
		}

		minute, err := tinystring.Convert(parts[1]).TrimSpace().Int()
		if err != nil {
			t.Fatalf("Failed to parse minute: %v", err)
		}

		second, err := tinystring.Convert(parts[2]).TrimSpace().Int()
		if err != nil {
			t.Fatalf("Failed to parse second: %v", err)
		}

		// Format the time with leading zeros for parsing.
		formattedTime := fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)

		_, err = time.Parse("15:04:05", formattedTime)
		if err != nil {
			t.Errorf("UnixNanoToTime() returned an invalid time string: %q, formatted: %q, error: %v", timeStr, formattedTime, err)
		}

		fmt.Println("Test finished successfully.")
	}()

	// Wait for the test to complete or timeout.
	select {
	case <-done:
		// Test finished.
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out")
	}
}

// The wasmbrowsertest runner requires a main function to be present.
// It will not be executed, but it is required for compilation.
func main() {
	// This line is required to prevent the compiler from optimizing away the syscall/js import.
	_ = js.Global()
}
