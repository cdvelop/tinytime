# tinytime

A minimal, portable time utility for Go and TinyGo with WebAssembly support. Automatically uses JavaScript Date APIs in WASM environments to keep binaries small.

## Quick Start

```go
import "github.com/cdvelop/tinytime"

func main() {
    tp := tinytime.NewTimeProvider()

    // Get current Unix timestamp in nanoseconds
    nano := tp.UnixNano()
    println("Current time:", nano)

    // Convert Unix seconds to formatted date
    date := tp.UnixSecondsToDate(1624397134)
    println("Date:", date) // "2021-06-22 21:25"

    // Convert Unix nanoseconds to time
    time := tp.UnixNanoToTime(nano)
    println("Time:", time) // "15:32:14"
}
```

## API Reference

### `NewTimeProvider() TimeProvider`

Creates a time provider instance. Automatically selects the appropriate implementation:
- **WASM**: Uses JavaScript Date APIs (smaller binaries)
- **Standard Go**: Uses `time` package

### `UnixNano() int64`

Returns current Unix timestamp in nanoseconds.

```go
nano := tp.UnixNano()
// Example: 1624397134562544800
```

### `UnixSecondsToDate(seconds int64) string`

Converts Unix timestamp (seconds) to formatted date string.

```go
date := tp.UnixSecondsToDate(1624397134)
// Returns: "2021-06-22 21:25"
```

### `UnixNanoToTime(input any) string`

Converts Unix timestamp (nanoseconds) to time string. Accepts multiple input types.

```go
// All of these work:
time1 := tp.UnixNanoToTime(int64(1624397134000000000))
time2 := tp.UnixNanoToTime(1624397134000000000) // int
time3 := tp.UnixNanoToTime("1624397134000000000") // string
// Returns: "17:25:34"
```

**Supported input types:** `int64`, `int`, `float64`, `string`

## WebAssembly Usage

When compiled for WebAssembly (`GOOS=js GOARCH=wasm`), tinytime automatically uses JavaScript's native Date APIs instead of bundling Go's time package. This results in significantly smaller binary sizes.

```bash
# Build for WebAssembly
GOOS=js GOARCH=wasm go build -o app.wasm .

# Run tests in browser
go install github.com/cdvelop/wasmtest@latest
wasmtest.RunTests("./tests", nil, 5*time.Minute)
```

## Testing

Run standard tests:
```bash
go test
```

Run WebAssembly tests in browser:
```bash
go test -tags=wasm
```

## Use Cases

- **Frontend applications** where binary size matters
- **TinyGo projects** requiring time utilities
- **WebAssembly modules** needing timestamp functionality
- **Cross-platform libraries** that work in both Go and browser environments

## Dependencies

- `github.com/cdvelop/tinystring` (for string parsing in WASM)