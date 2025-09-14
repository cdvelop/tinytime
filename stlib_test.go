package tinytime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cdvelop/tinytime"
)

func TestUnixNanoToTime(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test con timestamp conocido en la zona horaria local
	expected := ExpectedTimeString

	nanoTimestamp := TestUnixNano // convertir a nanosegundos

	result := tp.UnixNanoToTime(nanoTimestamp)
	if result != expected {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", nanoTimestamp, result, expected)
	}

	// Test con string
	result = tp.UnixNanoToTime((fmt.Sprintf("%d", nanoTimestamp)))
	if result != expected {
		t.Errorf("UnixNanoToTime(string) = %s; want %s", result, expected)
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

func TestUnixNanoToTimeWithDifferentTypes(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test con diferentes tipos de entrada
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", TestUnixNano},
		{"int", int(TestUnixNano)},
		{"float64", float64(TestUnixNano)},
		{"string", (fmt.Sprintf("%d", TestUnixNano))},
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

func TestUnixNanoFunction(t *testing.T) {
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

	// Test that timestamp is recent (within last hour)
	now := time.Now().UnixNano()
	if nano > now || nano < now-3600000000000 { // 1 hour in nanoseconds
		t.Errorf("UnixNano() returned timestamp not within reasonable range: %d", nano)
	}

	t.Logf("UnixNano: %d", nano)
}

func TestUnixSecondsToDate(t *testing.T) {
	tp := tinytime.NewTimeProvider()

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

func TestUnixSecondsToDateEdgeCases(t *testing.T) {
	tp := tinytime.NewTimeProvider()

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

func TestUnixNanoToTimeEdgeCases(t *testing.T) {
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

func TestTimeProviderConsistency(t *testing.T) {
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
