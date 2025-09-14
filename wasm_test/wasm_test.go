//go:build js && wasm
// +build js,wasm

package wasm_test

import (
	"fmt"
	"syscall/js"
	"testing"

	"github.com/cdvelop/tinytime"
)

func TestWasmUnixNano(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test that UnixNano returns a reasonable timestamp
	nano := tp.UnixNano()
	if nano <= 0 {
		t.Errorf("UnixNano() returned invalid timestamp: %d", nano)
	}

	// Test that multiple calls return increasing values
	nano2 := tp.UnixNano()
	if nano2 < nano {
		t.Errorf("UnixNano() not monotonic: %d < %d", nano2, nano)
	}

	t.Logf("UnixNano: %d", nano)
}

func TestWasmUnixSecondsToDate(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test with known timestamp
	result := tp.UnixSecondsToDate(1624397134)
	expected := "2021-06-22 21:25" // Correct date for timestamp 1624397134
	if result != expected {
		t.Errorf("UnixSecondsToDate(%d) = %s; want %s", 1624397134, result, expected)
	}

	// Test with current time (should not be empty)
	currentSeconds := int64(1624397134 + 1000) // Add some seconds
	result = tp.UnixSecondsToDate(currentSeconds)
	if result == "" {
		t.Error("UnixSecondsToDate returned empty string")
	}

	t.Logf("UnixSecondsToDate(%d) = %s", 1624397134, result)
}

func TestWasmUnixNanoToTime(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test with known timestamp
	nanoTimestamp := int64(1624397134000000000)
	expectedTime := "17:25:34"
	result := tp.UnixNanoToTime(nanoTimestamp)
	if result != expectedTime {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", nanoTimestamp, result, expectedTime)
	}

	// Test with string input
	result = tp.UnixNanoToTime(fmt.Sprintf("%d", nanoTimestamp))
	if result != expectedTime {
		t.Errorf("UnixNanoToTime(string) = %s; want %s", result, expectedTime)
	}

	// Test with sequential timestamps to verify order
	baseNano := nanoTimestamp
	var results []string
	for i := 0; i < 3; i++ {
		nano := baseNano + int64(i)*int64(1000000000) // Add 1 second in nanoseconds
		timeStr := tp.UnixNanoToTime(nano)
		results = append(results, timeStr)
		t.Logf("Nano: %d -> Time: %s", nano, timeStr)
	}

	// Verify that times are in order
	for i := 1; i < len(results); i++ {
		if results[i] <= results[i-1] {
			t.Errorf("Times not in order: %s <= %s", results[i], results[i-1])
		}
	}
}

func TestWasmUnixNanoToTimeWithDifferentTypes(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test with different input types
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", int64(1624397134000000000)},
		{"int", int(1624397134000000000)},
		{"float64", float64(1624397134000000000)},
		{"string", "1624397134000000000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tp.UnixNanoToTime(tc.input)

			if result == "" {
				t.Errorf("UnixNanoToTime returned empty string for type %s", tc.name)
			}

			t.Logf("Type %s: %v -> %s", tc.name, tc.input, result)
		})
	}

	// Test with unsupported type
	unsupportedResult := tp.UnixNanoToTime(make(chan int))
	if unsupportedResult != "" {
		t.Error("UnixNanoToTime should return empty string for unsupported types")
	}
}

func TestWasmEnvironmentDetection(t *testing.T) {
	// Test that we're running in WASM environment
	global := js.Global()

	// Check if we have access to JavaScript globals
	if global.IsUndefined() {
		t.Error("JavaScript global object should be available in WASM environment")
	}

	// Check if Date object is available (should be in browser or Node.js)
	dateObj := global.Get("Date")
	if dateObj.IsUndefined() {
		t.Error("Date object should be available in WASM environment")
	}

	t.Log("WASM environment detected successfully")
}

func BenchmarkWasmUnixNanoToTime(b *testing.B) {
	tp := tinytime.NewTimeProvider()

	for i := 0; i < b.N; i++ {
		_ = tp.UnixNanoToTime(int64(1624397134000000000))
	}
}

func BenchmarkWasmUnixSecondsToDate(b *testing.B) {
	tp := tinytime.NewTimeProvider()

	for i := 0; i < b.N; i++ {
		_ = tp.UnixSecondsToDate(1624397134)
	}
}

func TestWasmUnixSecondsToDateEdgeCases(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	testCases := []struct {
		name  string
		input int64
	}{
		{"epoch", 0},
		{"negative", -1},
		{"year_2000", 946684800},
		{"future", 2147483647},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tp.UnixSecondsToDate(tc.input)
			// Just check that it returns a valid formatted string
			if result == "" || len(result) != 16 || result[10] != ' ' {
				t.Errorf("UnixSecondsToDate(%d) returned invalid format: %s", tc.input, result)
			}
			// Check date format YYYY-MM-DD
			if len(result[:10]) != 10 || result[4] != '-' || result[7] != '-' {
				t.Errorf("UnixSecondsToDate(%d) date format invalid: %s", tc.input, result[:10])
			}
			// Check time format HH:MM
			if len(result[11:]) != 5 || result[13] != ':' {
				t.Errorf("UnixSecondsToDate(%d) time format invalid: %s", tc.input, result[11:])
			}
		})
	}
}

func TestWasmUnixNanoToTimeEdgeCases(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	testCases := []struct {
		name  string
		input any
		check func(string) bool // Custom check function since timezone affects results
	}{
		{"zero_nanoseconds", int64(0), func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"negative_nanoseconds", int64(-1000000000), func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"large_nanoseconds", int64(86400000000000), func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"zero_int", 0, func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"zero_float", 0.0, func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"zero_string", "0", func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"empty_string", "", func(s string) bool { return s == "" || (s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':') }}, // Allow time string since empty might convert to 0
		{"invalid_string", "abc", func(s string) bool { return s == "" }},
		{"nil_interface", nil, func(s string) bool { return s == "" }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tp.UnixNanoToTime(tc.input)
			if !tc.check(result) {
				t.Errorf("UnixNanoToTime(%v) = %q; check failed", tc.input, result)
			}
		})
	}
}

func TestWasmTimeProviderConsistency(t *testing.T) {
	// Test that the same TimeProvider instance behaves consistently
	tp1 := tinytime.NewTimeProvider()
	tp2 := tinytime.NewTimeProvider()

	// Both should return valid timestamps
	nano1 := tp1.UnixNano()
	nano2 := tp2.UnixNano()

	if nano1 <= 0 || nano2 <= 0 {
		t.Error("TimeProvider instances returned invalid timestamps")
	}

	// Test that both can format the same timestamp consistently
	testNano := int64(1624397134000000000) // 2021-06-22 21:25:34 UTC
	testTime := tp1.UnixNanoToTime(testNano)
	expected := tp2.UnixNanoToTime(testNano)

	if testTime != expected {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testTime, expected)
	}

	// Test that both format dates consistently
	testSeconds := int64(1624397134) // Same timestamp in seconds
	testDate := tp1.UnixSecondsToDate(testSeconds)
	expectedDate := tp2.UnixSecondsToDate(testSeconds)

	if testDate != expectedDate {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testDate, expectedDate)
	}
}
