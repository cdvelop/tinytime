package tinytime

// TimeProvider define la interfaz para utilidades de tiempo
// implementada tanto para Go estándar como para WASM/JS.
type TimeProvider interface {
	// UnixNano retrieves the current Unix timestamp in nanoseconds.
	// It creates a new JavaScript Date object, gets the timestamp in milliseconds,
	// converts it to nanoseconds, and returns the result as an int64.
	// eg: 1624397134562544800
	UnixNano() int64
	//	ts := int64(1609459200) // January 1, 2021 00:00:00 UTC
	//	formattedDate := UnixSecondsToDate(ts)
	//	println(formattedDate) // Output: "2021-01-01 00:00:00"
	UnixSecondsToDate(int64) string
	// UnixNanoToTime converts a Unix timestamp in nanoseconds to a formatted time string.
	// Format: "15:04:05" (hour:minute:second)
	// It accepts a parameter of type any and attempts to convert it to an int64 Unix timestamp in nanoseconds.
	// eg: 1624397134562544800 -> "15:32:14"
	// supported types: int64, int, float64, string
	UnixNanoToTime(any) string
}
