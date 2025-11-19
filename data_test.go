package tinytime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cdvelop/tinytime"
)

var (
	TestUnixNano     int64 = 1705307400000000000 // 2024-01-15 08:30:00 UTC
	TestDateStr      string = "2024-01-15"
	TestTimeStr      string = "08:30"
	TestDateTimeStr  string = "2024-01-15 08:30:00"
	TestMinutes      int16 = 510 // 8 * 60 + 30
)

// Shared test functions for all TimeProvider implementations

func FormatDateShared(t *testing.T, tp tinytime.TimeProvider) {
	result := tp.FormatDate(TestUnixNano)
	if result != TestDateStr {
		t.Errorf("FormatDate from nano = %s; want %s", result, TestDateStr)
	}

	result = tp.FormatDate(TestDateStr)
	if result != TestDateStr {
		t.Errorf("FormatDate from string = %s; want %s", result, TestDateStr)
	}
}

func FormatTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	result := tp.FormatTime(TestUnixNano)
	if result != "08:30:00" {
		t.Errorf("FormatTime from nano = %s; want %s", result, "08:30:00")
	}

	result = tp.FormatTime(TestMinutes)
	if result != TestTimeStr {
		t.Errorf("FormatTime from minutes = %s; want %s", result, TestTimeStr)
	}

	result = tp.FormatTime(TestTimeStr)
	if result != TestTimeStr {
		t.Errorf("FormatTime from string = %s; want %s", result, TestTimeStr)
	}
}

func FormatDateTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	result := tp.FormatDateTime(TestUnixNano)
	if result != TestDateTimeStr {
		t.Errorf("FormatDateTime from nano = %s; want %s", result, TestDateTimeStr)
	}

	result = tp.FormatDateTime(TestDateTimeStr)
	if result != TestDateTimeStr {
		t.Errorf("FormatDateTime from string = %s; want %s", result, TestDateTimeStr)
	}
}

func ParseDateShared(t *testing.T, tp tinytime.TimeProvider) {
	expectedNano, _ := time.ParseInLocation("2006-01-02", TestDateStr, time.UTC)

	nano, err := tp.ParseDate(TestDateStr)
	if err != nil {
		t.Fatalf("ParseDate failed: %v", err)
	}
	if nano != expectedNano.UnixNano() {
		t.Errorf("ParseDate = %d; want %d", nano, expectedNano.UnixNano())
	}

	_, err = tp.ParseDate("invalid-date")
	if err == nil {
		t.Error("ParseDate should have failed for invalid date")
	}
}

func ParseTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	minutes, err := tp.ParseTime(TestTimeStr)
	if err != nil {
		t.Fatalf("ParseTime failed: %v", err)
	}
	if minutes != TestMinutes {
		t.Errorf("ParseTime = %d; want %d", minutes, TestMinutes)
	}

	_, err = tp.ParseTime("invalid-time")
	if err == nil {
		t.Error("ParseTime should have failed for invalid time")
	}
}

func ParseDateTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	nano, err := tp.ParseDateTime(TestDateStr, TestTimeStr)
	if err != nil {
		t.Fatalf("ParseDateTime failed: %v", err)
	}
	if nano != TestUnixNano {
		t.Errorf("ParseDateTime = %d; want %d", nano, TestUnixNano)
	}

	_, err = tp.ParseDateTime("invalid-date", "invalid-time")
	if err == nil {
		t.Error("ParseDateTime should have failed for invalid data")
	}
}

func IsTodayShared(t *testing.T, tp tinytime.TimeProvider) {
	now := tp.UnixNano()
	if !tp.IsToday(now) {
		t.Error("IsToday failed for current time")
	}

	yesterday := now - (24 * int64(time.Hour))
	if tp.IsToday(yesterday) {
		t.Error("IsToday failed for yesterday")
	}
}

func IsPastShared(t *testing.T, tp tinytime.TimeProvider) {
	past := tp.UnixNano() - 1000
	if !tp.IsPast(past) {
		t.Error("IsPast failed for past time")
	}
	if tp.IsPast(tp.UnixNano() + 1000) {
		t.Error("IsPast failed for future time")
	}
}

func IsFutureShared(t *testing.T, tp tinytime.TimeProvider) {
	future := tp.UnixNano() + 1000
	if !tp.IsFuture(future) {
		t.Error("IsFuture failed for future time")
	}
	if tp.IsFuture(tp.UnixNano() - 1000) {
		t.Error("IsFuture failed for past time")
	}
}

func DaysBetweenShared(t *testing.T, tp tinytime.TimeProvider) {
	dayInNanos := int64(24 * time.Hour)
	nano1 := TestUnixNano
	nano2 := nano1 + 7*dayInNanos

	days := tp.DaysBetween(nano1, nano2)
	if days != 7 {
		t.Errorf("DaysBetween = %d; want 7", days)
	}

	days = tp.DaysBetween(nano2, nano1)
	if days != -7 {
		t.Errorf("DaysBetween = %d; want -7", days)
	}
}

// Existing shared tests can remain below
// ...
// Shared test data for tinytime tests (used by both native and wasm tests)
var (
	// A known Unix timestamp in seconds (2021-06-22 15:32:14 UTC)
	GlobalTestUnixSeconds int64 = 1624397134

	// The corresponding time.Time value (UTC)
	GlobalTestTime = time.Unix(GlobalTestUnixSeconds, 0)

	// The expected formatted time string for UnixNanoToTime
	GlobalExpectedTimeString = GlobalTestTime.Format("15:04:05")

	// A helper to get a nano timestamp (int64)
	GlobalTestUnixNano int64 = GlobalTestUnixSeconds * 1000000000

	// Test case types for UnixNanoToTimeWithDifferentTypes
	GlobalTestCaseTypes = []string{"int64", "int", "float64", "string"}
)
func UnixNanoToTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test con timestamp conocido
	result := tp.UnixNanoToTime(GlobalTestUnixNano)
	if result != GlobalExpectedTimeString {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", GlobalTestUnixNano, result, GlobalExpectedTimeString)
	}

	// Test con string
	result = tp.UnixNanoToTime(fmt.Sprintf("%d", GlobalTestUnixNano))
	if result != GlobalExpectedTimeString {
		t.Errorf("UnixNanoToTime(string) = %s; want %s", result, GlobalExpectedTimeString)
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
func UnixNanoToTimeWithDifferentTypesShared(t *testing.T, tp tinytime.TimeProvider) {
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", GlobalTestUnixNano},
		{"int", int(GlobalTestUnixNano)},
		{"float64", float64(GlobalTestUnixNano)},
		{"string", fmt.Sprintf("%d", GlobalTestUnixNano)},
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
func UnixSecondsToDateShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test with known timestamp - check that it returns a properly formatted string
	result := tp.UnixSecondsToDate(GlobalTestUnixSeconds)
	// Just verify the format is correct: YYYY-MM-DD HH:MM
	if len(result) != 16 || result[10] != ' ' || result[4] != '-' || result[7] != '-' || result[13] != ':' {
		t.Errorf("UnixSecondsToDate(%d) returned invalid format: %s", GlobalTestUnixSeconds, result)
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

	t.Logf("UnixSecondsToDate(%d) = %s", GlobalTestUnixSeconds, result)
}
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
func TimeProviderConsistencyShared(t *testing.T, tp1, tp2 tinytime.TimeProvider) {
	// Both should return valid timestamps
	nano1 := tp1.UnixNano()
	nano2 := tp2.UnixNano()

	if nano1 <= 0 || nano2 <= 0 {
		t.Error("TimeProvider instances returned invalid timestamps")
	}

	// Test that both can format the same timestamp consistently
	testTime := tp1.UnixNanoToTime(GlobalTestUnixNano)
	expected := tp2.UnixNanoToTime(GlobalTestUnixNano)

	if testTime != expected {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testTime, expected)
	}

	// Test that both format dates consistently
	testDate := tp1.UnixSecondsToDate(GlobalTestUnixSeconds)
	expectedDate := tp2.UnixSecondsToDate(GlobalTestUnixSeconds)

	if testDate != expectedDate {
		t.Errorf("TimeProvider instances inconsistent: %s != %s", testDate, expectedDate)
	}
}
