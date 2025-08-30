# tinytime

TinyGo-compatible time package for Go, using syscall/js on WASM to keep binaries small and leverage native JS time APIs.


Minimal and portable time utility package for Go and TinyGo with WebAssembly support.
When compiled for wasm (GOOS=js GOARCH=wasm), it uses syscall/js to access native JavaScript time APIs (Date.now, performance.now, etc.), drastically reducing binary size by avoiding the Go standard library.
For non-WASM targets, it falls back to the standard time package.
Ideal for frontend projects where binary size and compatibility with JavaScript environments matter.

## API Usage

The `tinytime` package provides a `TimeProvider` interface that abstracts time operations for both standard Go and WebAssembly (WASM) environments.

### `NewTimeProvider() TimeProvider`

This is the entry point of the library. It returns an implementation of the `TimeProvider` interface that is appropriate for the current build target (standard Go or WASM).

### `TimeProvider` Interface

The `TimeProvider` interface has the following methods:

- **`UnixNano() int64`**: Returns the current Unix timestamp in nanoseconds.

- **`UnixSecondsToDate(unixSeconds int64) string`**: Converts a Unix timestamp in seconds to a formatted date string (`YYYY-MM-DD HH:MM`).

- **`UnixNanoToTime(input any) string`**: Converts a Unix timestamp in nanoseconds to a formatted time string (`HH:MM:SS`). It accepts `int64`, `int`, `float64`, or a numeric `string` as input.

### Example

Here is a basic example of how to use the `tinytime` library:

```go
package main

import (
	"fmt"
	"github.com/cdvelop/tinytime"
)

func main() {
	// Get the time provider
	t := tinytime.NewTimeProvider()

	// Get the current time in Unix nanoseconds
	nano := t.UnixNano()
	fmt.Printf("Current Unix Nano: %d\n", nano)

	// Convert Unix seconds to a date string
	// Example timestamp for January 1, 2025 00:00:00 UTC
	seconds := int64(1735689600)
	dateStr := t.UnixSecondsToDate(seconds)
	fmt.Printf("Date from seconds: %s\n", dateStr)

	// Convert Unix nanoseconds to a time string
	timeStr := t.UnixNanoToTime(nano)
	fmt.Printf("Time from nanoseconds: %s\n", timeStr)
}
```