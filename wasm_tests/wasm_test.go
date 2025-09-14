//go:build js && wasm
// +build js,wasm

package wasm_test

import (
	"syscall/js"
	"testing"

	"github.com/cdvelop/tinytime"
)

// Test data for WASM tests (cannot access data_test.go due to build tags)
const (
	testUnixSeconds    = int64(1624397134)
	testUnixNano       = int64(1624397134000000000)
	expectedTimeString = "17:25:34"
)

// Tests for WebAssembly environment

func TestWasmUnixNano(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test that UnixNano returns a reasonable timestamp
	nano := tp.UnixNano()
	if nano <= 0 {
		t.Errorf("UnixNano() returned invalid timestamp: %d", nano)
	}

	// Test that multiple calls return increasing values (monotonic)
	nano2 := tp.UnixNano()
	if nano2 < nano {
		t.Errorf("UnixNano() not monotonic: %d < %d", nano2, nano)
	}

	t.Logf("UnixNano: %d", nano)
}

func TestWasmUnixSecondsToDate(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test with known timestamp - check that it returns a properly formatted string
	result := tp.UnixSecondsToDate(testUnixSeconds)
	// Just verify the format is correct: YYYY-MM-DD HH:MM
	if len(result) != 16 || result[10] != ' ' || result[4] != '-' || result[7] != '-' || result[13] != ':' {
		t.Errorf("UnixSecondsToDate(%d) returned invalid format: %s", testUnixSeconds, result)
	}

	t.Logf("UnixSecondsToDate(%d) = %s", testUnixSeconds, result)
}

func TestWasmUnixNanoToTime(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test with known timestamp
	result := tp.UnixNanoToTime(testUnixNano)
	if result != expectedTimeString {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", testUnixNano, result, expectedTimeString)
	}

	t.Logf("UnixNanoToTime test passed")
}

func TestWasmUnixNanoToTimeWithDifferentTypes(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test with different input types
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", testUnixNano},
		{"int", int(testUnixNano)},
		{"float64", float64(testUnixNano)},
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
		_ = tp.UnixNanoToTime(testUnixNano)
	}
}

func BenchmarkWasmUnixSecondsToDate(b *testing.B) {
	tp := tinytime.NewTimeProvider()

	for i := 0; i < b.N; i++ {
		_ = tp.UnixSecondsToDate(testUnixSeconds)
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
		})
	}
}

func TestWasmUnixNanoToTimeEdgeCases(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	testCases := []struct {
		name  string
		input any
		check func(string) bool
	}{
		{"zero_nanoseconds", int64(0), func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"negative_nanoseconds", int64(-1000000000), func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"zero_int", 0, func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"zero_float", 0.0, func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"zero_string", "0", func(s string) bool { return s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':' }},
		{"empty_string", "", func(s string) bool { return s == "" || (s != "" && len(s) == 8 && s[2] == ':' && s[5] == ':') }},
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

func TestWasmUnixNanoToTimeSingleDigitHour(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Timestamp for 01:25:34 (subtract 16 hours from testUnixNano)
	singleDigitNano := testUnixNano - 16*3600*1e9
	result := tp.UnixNanoToTime(singleDigitNano)
	expected := "01:25:34"

	if result != expected {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", singleDigitNano, result, expected)
	}

	t.Logf("Single digit hour test: %s", result)
}

func TestWasmTimeProviderConsistency(t *testing.T) {
	tp1 := tinytime.NewTimeProvider()
	tp2 := tinytime.NewTimeProvider()

	// Both should return valid timestamps
	nano1 := tp1.UnixNano()
	nano2 := tp2.UnixNano()

	if nano1 <= 0 || nano2 <= 0 {
		t.Error("TimeProvider instances returned invalid timestamps")
	}

	// Test that both can format the same timestamp consistently
	testTime := tp1.UnixNanoToTime(testUnixNano)
	expected := tp2.UnixNanoToTime(testUnixNano)

	if testTime != expected {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testTime, expected)
	}

	// Test that both format dates consistently
	testDate := tp1.UnixSecondsToDate(testUnixSeconds)
	expectedDate := tp2.UnixSecondsToDate(testUnixSeconds)

	if testDate != expectedDate {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testDate, expectedDate)
	}
}
