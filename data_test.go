package tinytime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cdvelop/tinytime"
)

// Shared test data for tinytime tests (used by both native and wasm tests)
var (
	// A known Unix timestamp in seconds (2021-06-22 15:32:14 UTC)
	TestUnixSeconds int64 = 1624397134

	// The corresponding time.Time value (UTC)
	TestTime = time.Unix(TestUnixSeconds, 0)

	// The expected formatted time string for UnixNanoToTime
	ExpectedTimeString = TestTime.Format("15:04:05")

	// A helper to get a nano timestamp (int64)
	TestUnixNano int64 = TestUnixSeconds * 1000000000

	// Test case types for UnixNanoToTimeWithDifferentTypes
	TestCaseTypes = []string{"int64", "int", "float64", "string"}
)

// Shared test functions that can be used by both native and WASM test files

// UnixNanoToTimeShared tests UnixNanoToTime with shared validation logic
func UnixNanoToTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test con timestamp conocido
	result := tp.UnixNanoToTime(TestUnixNano)
	if result != ExpectedTimeString {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", TestUnixNano, result, ExpectedTimeString)
	}

	// Test con string
	result = tp.UnixNanoToTime(fmt.Sprintf("%d", TestUnixNano))
	if result != ExpectedTimeString {
		t.Errorf("UnixNanoToTime(string) = %s; want %s", result, ExpectedTimeString)
	}

	// Test con timestamps secuenciales para verificar orden
	now := time.Now()
	baseNano := now.UnixNano()

	var results []string
	for i := 0; i < 3; i++ {
		nano := baseNano + int64(i)*int64(time.Second) // Incrementar 1 segundo
		timeStr := tp.UnixNanoToTime(nano)
		results = append(results, timeStr)
		t.Logf("Nano: %d -> Time: %s", nano, timeStr)
	}

	// Verificar que los tiempos están en orden
	for i := 1; i < len(results); i++ {
		if results[i] <= results[i-1] {
			t.Errorf("Los timestamps no están en orden: %s <= %s", results[i], results[i-1])
		}
	}
}

// UnixNanoToTimeWithDifferentTypesShared tests UnixNanoToTime with different input types
func UnixNanoToTimeWithDifferentTypesShared(t *testing.T, tp tinytime.TimeProvider) {
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", TestUnixNano},
		{"int", int(TestUnixNano)},
		{"float64", float64(TestUnixNano)},
		{"string", fmt.Sprintf("%d", TestUnixNano)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tp.UnixNanoToTime(tc.input)

			if result == "" {
				t.Errorf("UnixNanoToTime devolvió string vacío para tipo %s", tc.name)
			}

			t.Logf("Tipo %s: %v -> %s", tc.name, tc.input, result)
		})
	}

	// Test con tipo no soportado
	invalidResult := tp.UnixNanoToTime(make(chan int))
	if invalidResult != "" {
		t.Error("UnixNanoToTime debería devolver string vacío para tipos no soportados")
	}
}

// UnixNanoShared tests UnixNano function
func UnixNanoShared(t *testing.T, tp tinytime.TimeProvider) {
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

	// Test that timestamp is recent (within last hour)
	now := time.Now().UnixNano()
	if nano > now || nano < now-3600000000000 { // 1 hour in nanoseconds
		t.Errorf("UnixNano() returned timestamp not within reasonable range: %d", nano)
	}

	t.Logf("UnixNano: %d", nano)
}

// UnixSecondsToDateShared tests UnixSecondsToDate function
func UnixSecondsToDateShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test with known timestamp - check that it returns a properly formatted string
	result := tp.UnixSecondsToDate(TestUnixSeconds)
	// Just verify the format is correct: YYYY-MM-DD HH:MM
	if len(result) != 16 || result[10] != ' ' || result[4] != '-' || result[7] != '-' || result[13] != ':' {
		t.Errorf("UnixSecondsToDate(%d) returned invalid format: %s", TestUnixSeconds, result)
	}

	// Test with zero timestamp
	result = tp.UnixSecondsToDate(0)
	if result == "" || len(result) != 16 || result[10] != ' ' {
		t.Errorf("UnixSecondsToDate(0) returned invalid format: %s", result)
	}

	// Test with negative timestamp
	result = tp.UnixSecondsToDate(-1)
	if result == "" || len(result) != 16 || result[10] != ' ' {
		t.Errorf("UnixSecondsToDate(-1) returned invalid format: %s", result)
	}

	// Test with current time (should not be empty)
	currentSeconds := time.Now().Unix()
	result = tp.UnixSecondsToDate(currentSeconds)
	if result == "" {
		t.Error("UnixSecondsToDate returned empty string for current time")
	}

	t.Logf("UnixSecondsToDate(%d) = %s", TestUnixSeconds, result)
}

// UnixSecondsToDateEdgeCasesShared tests UnixSecondsToDate with edge cases
func UnixSecondsToDateEdgeCasesShared(t *testing.T, tp tinytime.TimeProvider) {
	testCases := []struct {
		name  string
		input int64
	}{
		{"epoch", 0},
		{"negative", -1},
		{"year_2000", 946684800},
		{"future", 2147483647},        // 32-bit signed int max
		{"distant_past", -2147483648}, // 32-bit signed int min
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

// UnixNanoToTimeEdgeCasesShared tests UnixNanoToTime with edge cases
func UnixNanoToTimeEdgeCasesShared(t *testing.T, tp tinytime.TimeProvider) {
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

// TimeProviderConsistencyShared tests consistency between TimeProvider instances
func TimeProviderConsistencyShared(t *testing.T, tp1, tp2 tinytime.TimeProvider) {
	// Both should return valid timestamps
	nano1 := tp1.UnixNano()
	nano2 := tp2.UnixNano()

	if nano1 <= 0 || nano2 <= 0 {
		t.Error("TimeProvider instances returned invalid timestamps")
	}

	// Test that both can format the same timestamp consistently
	testTime := tp1.UnixNanoToTime(TestUnixNano)
	expected := tp2.UnixNanoToTime(TestUnixNano)

	if testTime != expected {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testTime, expected)
	}

	// Test that both format dates consistently
	testDate := tp1.UnixSecondsToDate(TestUnixSeconds)
	expectedDate := tp2.UnixSecondsToDate(TestUnixSeconds)

	if testDate != expectedDate {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testDate, expectedDate)
	}
}
