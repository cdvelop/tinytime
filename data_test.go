package tinytime_test

import (
	"testing"

	"github.com/cdvelop/tinytime"
)

// Global test constants
const (
	GlobalTestUnixSeconds    int64  = 1609459200 // 2021-01-01 00:00:00 UTC
	GlobalExpectedTimeString string = "00:00:00"
)

var (
	GlobalTestUnixNano int64 = GlobalTestUnixSeconds * 1000000000
)

// Test FormatDate
func FormatDateShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test with UnixNano (int64)
	result := tp.FormatDate(GlobalTestUnixNano)
	if len(result) != 10 || result[4] != '-' || result[7] != '-' {
		t.Errorf("FormatDate(%d) returned invalid format: %s", GlobalTestUnixNano, result)
	}

	// Test with string passthrough
	result = tp.FormatDate("2024-01-15")
	if result != "2024-01-15" {
		t.Errorf("FormatDate(string) = %s; want 2024-01-15", result)
	}

	// Test with zero value
	result = tp.FormatDate(int64(0))
	if result != "1970-01-01" {
		t.Errorf("FormatDate(0) = %s; want 1970-01-01", result)
	}

	t.Logf("FormatDate tests passed")
}

// Test FormatTime
func FormatTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test with int64 (UnixNano)
	result := tp.FormatTime(GlobalTestUnixNano)
	if len(result) != 8 || result[2] != ':' || result[5] != ':' {
		t.Errorf("FormatTime(int64) returned invalid format: %s", result)
	}

	// Test with int16 (minutes since midnight)
	result = tp.FormatTime(int16(510)) // 08:30
	if result != "08:30" {
		t.Errorf("FormatTime(int16(510)) = %s; want 08:30", result)
	}

	// Test with string passthrough
	result = tp.FormatTime("14:45")
	if result != "14:45" {
		t.Errorf("FormatTime(string) = %s; want 14:45", result)
	}

	t.Logf("FormatTime tests passed")
}

// Test FormatDateTime
func FormatDateTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test with UnixNano (int64)
	result := tp.FormatDateTime(GlobalTestUnixNano)
	// Format: YYYY-MM-DD HH:MM:SS (19 chars)
	if len(result) != 19 || result[10] != ' ' || result[13] != ':' || result[16] != ':' {
		t.Errorf("FormatDateTime(%d) returned invalid format: %s", GlobalTestUnixNano, result)
	}

	// Test with zero value
	result = tp.FormatDateTime(int64(0))
	if result != "1970-01-01 00:00:00" {
		t.Errorf("FormatDateTime(0) = %s; want 1970-01-01 00:00:00", result)
	}

	// Test with current time
	currentNano := tp.UnixNano()
	result = tp.FormatDateTime(currentNano)
	if result == "" {
		t.Error("FormatDateTime returned empty string for current time")
	}

	t.Logf("FormatDateTime tests passed")
}

// Test FormatDateTimeShort
func FormatDateTimeShortShared(t *testing.T, tp tinytime.TimeProvider) {
	// Test with int64 (UnixNano timestamp)
	nano := int64(1705307400000000000) // 2024-01-15 08:30:00 UTC
	result := tp.FormatDateTimeShort(nano)
	if result != "2024-01-15 08:30" {
		t.Errorf("FormatDateTimeShort(nano) = %q, want %q", result, "2024-01-15 08:30")
	}

	// Test with valid string passthrough
	result = tp.FormatDateTimeShort("2024-01-15 08:30")
	if result != "2024-01-15 08:30" {
		t.Errorf("FormatDateTimeShort(string) = %q, want %q", result, "2024-01-15 08:30")
	}

	// Test with zero timestamp (epoch)
	result = tp.FormatDateTimeShort(int64(0))
	if result != "1970-01-01 00:00" {
		t.Errorf("FormatDateTimeShort(0) = %q, want %q", result, "1970-01-01 00:00")
	}

	// Test with current timestamp (should be 16 chars)
	currentNano := tp.UnixNano()
	result = tp.FormatDateTimeShort(currentNano)
	if len(result) != 16 {
		t.Errorf("FormatDateTimeShort(current) length = %d, want 16", len(result))
	}

	t.Logf("FormatDateTimeShort tests passed")
}

// Test UnixNano
func UnixNanoShared(t *testing.T, tp tinytime.TimeProvider) {
	nano := tp.UnixNano()

	// Check it's a reasonable timestamp (not zero, not negative, not too far in future)
	if nano <= 0 {
		t.Errorf("UnixNano() returned non-positive value: %d", nano)
	}

	// Test that timestamp is recent (within last 10 seconds to allow for clock drift)
	now := tp.UnixNano()
	diff := nano - now
	if diff < 0 {
		diff = -diff
	}
	if diff > 10000000000 {
		t.Errorf("UnixNano() returned timestamp too far from current time: %d (diff: %d ns)", nano, diff)
	}

	t.Logf("UnixNano: %d", nano)
}

// Test ParseDate
func ParseDateShared(t *testing.T, tp tinytime.TimeProvider) {
	// Valid date
	nano, err := tp.ParseDate("2024-01-15")
	if err != nil {
		t.Errorf("ParseDate(2024-01-15) failed: %v", err)
	}
	if nano <= 0 {
		t.Errorf("ParseDate returned invalid nano: %d", nano)
	}

	// Invalid date
	_, err = tp.ParseDate("invalid")
	if err == nil {
		t.Error("ParseDate(invalid) should return error")
	}

	// Invalid date (Feb 30)
	_, err = tp.ParseDate("2024-02-30")
	if err == nil {
		t.Error("ParseDate(2024-02-30) should return error")
	}

	t.Logf("ParseDate tests passed")
}

// Test ParseTime
func ParseTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	// Valid time
	minutes, err := tp.ParseTime("08:30")
	if err != nil {
		t.Errorf("ParseTime(08:30) failed: %v", err)
	}
	if minutes != 510 {
		t.Errorf("ParseTime(08:30) = %d; want 510", minutes)
	}

	// With seconds (should ignore)
	minutes, err = tp.ParseTime("08:30:45")
	if err != nil {
		t.Errorf("ParseTime(08:30:45) failed: %v", err)
	}
	if minutes != 510 {
		t.Errorf("ParseTime(08:30:45) = %d; want 510", minutes)
	}

	// Invalid time
	_, err = tp.ParseTime("invalid")
	if err == nil {
		t.Error("ParseTime(invalid) should return error")
	}

	// Invalid hours
	_, err = tp.ParseTime("25:00")
	if err == nil {
		t.Error("ParseTime(25:00) should return error")
	}

	t.Logf("ParseTime tests passed")
}

// Test ParseDateTime
func ParseDateTimeShared(t *testing.T, tp tinytime.TimeProvider) {
	// Valid date + time
	nano, err := tp.ParseDateTime("2024-01-15", "08:30")
	if err != nil {
		t.Errorf("ParseDateTime failed: %v", err)
	}
	if nano <= 0 {
		t.Errorf("ParseDateTime returned invalid nano: %d", nano)
	}

	// Invalid date
	_, err = tp.ParseDateTime("invalid", "08:30")
	if err == nil {
		t.Error("ParseDateTime(invalid date) should return error")
	}

	// Invalid time
	_, err = tp.ParseDateTime("2024-01-15", "invalid")
	if err == nil {
		t.Error("ParseDateTime(invalid time) should return error")
	}

	t.Logf("ParseDateTime tests passed")
}

// Test IsToday
func IsTodayShared(t *testing.T, tp tinytime.TimeProvider) {
	// Current time should be today
	now := tp.UnixNano()
	if !tp.IsToday(now) {
		t.Error("IsToday(now) should return true")
	}

	// Yesterday should not be today
	yesterday := now - (24 * 60 * 60 * 1000000000)
	if tp.IsToday(yesterday) {
		t.Error("IsToday(yesterday) should return false")
	}

	// Tomorrow should not be today
	tomorrow := now + (24 * 60 * 60 * 1000000000)
	if tp.IsToday(tomorrow) {
		t.Error("IsToday(tomorrow) should return false")
	}

	t.Logf("IsToday tests passed")
}

// Test IsPast
func IsPastShared(t *testing.T, tp tinytime.TimeProvider) {
	now := tp.UnixNano()

	// Past timestamp
	past := now - 1000000000
	if !tp.IsPast(past) {
		t.Error("IsPast(past) should return true")
	}

	// Future timestamp
	future := now + 1000000000
	if tp.IsPast(future) {
		t.Error("IsPast(future) should return false")
	}

	t.Logf("IsPast tests passed")
}

// Test IsFuture
func IsFutureShared(t *testing.T, tp tinytime.TimeProvider) {
	now := tp.UnixNano()

	// Future timestamp
	future := now + 1000000000
	if !tp.IsFuture(future) {
		t.Error("IsFuture(future) should return true")
	}

	// Past timestamp
	past := now - 1000000000
	if tp.IsFuture(past) {
		t.Error("IsFuture(past) should return false")
	}

	t.Logf("IsFuture tests passed")
}

// Test DaysBetween
func DaysBetweenShared(t *testing.T, tp tinytime.TimeProvider) {
	// 7 days apart
	nano1 := int64(1705276800000000000) // 2024-01-15
	nano2 := int64(1705881600000000000) // 2024-01-22

	days := tp.DaysBetween(nano1, nano2)
	if days != 7 {
		t.Errorf("DaysBetween = %d; want 7", days)
	}

	// Reversed (negative)
	days = tp.DaysBetween(nano2, nano1)
	if days != -7 {
		t.Errorf("DaysBetween(reversed) = %d; want -7", days)
	}

	// Same day
	days = tp.DaysBetween(nano1, nano1)
	if days != 0 {
		t.Errorf("DaysBetween(same) = %d; want 0", days)
	}

	t.Logf("DaysBetween tests passed")
}
