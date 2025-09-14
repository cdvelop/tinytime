package tinytime_test

import "time"

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
