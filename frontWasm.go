//go:build wasm
// +build wasm

package tinytime

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	. "github.com/cdvelop/tinystring"
)

// timeClient implements TimeProvider for WASM/JS environments using the JavaScript Date API.
type timeClient struct {
	dateCtor js.Value
}

// NewTimeProvider returns the correct implementation for WASM.
func NewTimeProvider() TimeProvider {
	return &timeClient{
		dateCtor: js.Global().Get("Date"),
	}
}

func (tc *timeClient) UnixNano() int64 {
	jsDate := tc.dateCtor.New()
	msTimestamp := jsDate.Call("getTime").Float()
	// Convert milliseconds to nanoseconds
	return int64(msTimestamp) * 1000000
}

func (tc *timeClient) FormatDate(value any) string {
	switch v := value.(type) {
	case int64:
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		return jsDate.Call("toISOString").String()[0:10]
	case string:
		if _, err := time.Parse("2006-01-02", v); err == nil {
			return v
		}
	}
	return ""
}

func (tc *timeClient) FormatTime(value any) string {
	switch v := value.(type) {
	case int64: // UnixNano
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		hours := jsDate.Call("getUTCHours").Int()
		minutes := jsDate.Call("getUTCMinutes").Int()
		seconds := jsDate.Call("getUTCSeconds").Int()
		return Fmt("%02d:%02d:%02d", hours, minutes, seconds)
	case int16: // Minutes since midnight
		hours := v / 60
		minutes := v % 60
		return Fmt("%02d:%02d", hours, minutes)
	case string:
		if strings.Count(v, ":") >= 1 {
			return v
		}
	}
	return ""
}

func (tc *timeClient) FormatDateTime(value any) string {
	switch v := value.(type) {
	case int64:
		jsDate := tc.dateCtor.New(float64(v) / 1e6)
		iso := jsDate.Call("toISOString").String()
		return iso[0:10] + " " + iso[11:19]
	case string:
		if _, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return v
		}
	}
	return ""
}

func (tc *timeClient) ParseDate(dateStr string) (int64, error) {
	// Parse using time.Parse to validate format strictly
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %s", dateStr)
	}

	jsDate := tc.dateCtor.New(dateStr + "T00:00:00Z")
	if jsDate.Call("toString").String() == "Invalid Date" {
		return 0, fmt.Errorf("invalid date format: %s", dateStr)
	}

	// Verify date components match (JS Date auto-corrects invalid dates like Feb 30)
	year := jsDate.Call("getUTCFullYear").Int()
	month := jsDate.Call("getUTCMonth").Int() + 1
	day := jsDate.Call("getUTCDate").Int()
	expected := Fmt("%04d-%02d-%02d", year, month, day)
	if expected != dateStr {
		return 0, fmt.Errorf("invalid date: %s (auto-corrected to %s)", dateStr, expected)
	}

	ms := jsDate.Call("getTime").Float()
	return int64(ms) * 1000000, nil
}

func (tc *timeClient) ParseTime(timeStr string) (int16, error) {
	return parseTime(timeStr)
}

func (tc *timeClient) ParseDateTime(dateStr, timeStr string) (int64, error) {
	if len(timeStr) == 5 {
		timeStr += ":00"
	}
	isoStr := dateStr + "T" + timeStr + "Z"
	jsDate := tc.dateCtor.New(isoStr)
	if jsDate.Call("toString").String() == "Invalid Date" {
		return 0, fmt.Errorf("invalid date/time format: %s %s", dateStr, timeStr)
	}
	ms := jsDate.Call("getTime").Float()
	return int64(ms) * 1000000, nil
}

func (tc *timeClient) IsToday(nano int64) bool {
	jsDate := tc.dateCtor.New(float64(nano) / 1e6)
	now := tc.dateCtor.New()
	return jsDate.Call("toDateString").String() == now.Call("toDateString").String()
}

func (tc *timeClient) IsPast(nano int64) bool {
	return nano < tc.UnixNano()
}

func (tc *timeClient) IsFuture(nano int64) bool {
	return nano > tc.UnixNano()
}

func (tc *timeClient) DaysBetween(nano1, nano2 int64) int {
	return daysBetween(nano1, nano2)
}
