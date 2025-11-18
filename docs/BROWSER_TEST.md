# Browser Testing Guide for TinyTime

## Overview

TinyTime includes WebAssembly (WASM) tests that run directly in a browser environment to validate the JavaScript Date API integration. This ensures that the frontend implementation works correctly in real browser contexts.

## Quick Start

### Installation

Install the `wasmbrowsertest` tool:

```bash
go install github.com/agnivade/wasmbrowsertest@latest
```

Rename the binary to match Go's execution convention:

```bash
mv $(go env GOPATH)/bin/wasmbrowsertest $(go env GOPATH)/bin/go_js_wasm_exec
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$(go env GOPATH)/bin:$PATH
```

### Running WASM Tests

Execute tests in the browser:

```bash
cd /path/to/tinytime
GOOS=js GOARCH=wasm go test ./wasm_tests/...
```

**Note:** Tests run headlessly by default (no visible browser window).

### Running with Visible Browser

To see the browser during test execution:

```bash
WASM_HEADLESS=off GOOS=js GOARCH=wasm go test ./wasm_tests/...
```

## Test Structure

### Current WASM Tests

Location: `tinytime/wasm_tests/wasm_test.go`

**Tests Included:**
- `TestWasmUnixNano` - Validates UnixNano() returns valid timestamps
- `TestWasmUnixSecondsToDate` - Tests date formatting with JS Date API
- `TestWasmUnixNanoToTime` - Tests time formatting
- `TestWasmUnixNanoToTimeWithDifferentTypes` - Type conversion tests
- `TestWasmEnvironmentDetection` - Verifies browser environment
- `TestWasmUnixSecondsToDateEdgeCases` - Edge case handling
- `TestWasmUnixNanoToTimeEdgeCases` - Boundary value tests
- `TestWasmTimeProviderConsistency` - Cross-instance consistency

**Benchmarks:**
- `BenchmarkWasmUnixNanoToTime`
- `BenchmarkWasmUnixSecondsToDate`

### Build Tags

WASM tests use build constraints:

```go
//go:build js && wasm
// +build js,wasm
```

This ensures they only compile when targeting WebAssembly.

## How It Works

### The Magic Behind wasmbrowsertest

1. **Automatic Binary Naming:** Go looks for `go_js_wasm_exec` when `GOOS=js GOARCH=wasm`
2. **Test Compilation:** Your test is compiled to a `.wasm` file
3. **HTML Generation:** `wasmbrowsertest` creates HTML with `wasm_exec.js`
4. **Browser Launch:** Chrome/Chromium starts via ChromeDP protocol
5. **Test Execution:** Tests run in browser, results streamed back
6. **Cleanup:** Browser closes, results displayed

### Architecture

```
go test (WASM) → go_js_wasm_exec → wasmbrowsertest
                                    ↓
                              Chrome Browser
                                    ↓
                            JavaScript Date API
                                    ↓
                          tinytime frontWasm.go
```

## CI/CD Integration

### GitHub Actions

Add to `.github/workflows/ci.yml`:

```yaml
name: WASM Tests
on: [push, pull_request]

jobs:
  wasm-test:
    runs-on: ubuntu-latest
    steps:
    - name: Install Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install Chrome
      uses: browser-actions/setup-chrome@latest
    
    - name: Install wasmbrowsertest
      run: go install github.com/agnivade/wasmbrowsertest@latest
    
    - name: Setup wasmexec
      run: mv $(go env GOPATH)/bin/wasmbrowsertest $(go env GOPATH)/bin/go_js_wasm_exec
    
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Run WASM tests
      run: GOOS=js GOARCH=wasm go test ./wasm_tests/...
```

### Travis CI

Add to `.travis.yml`:

```yaml
addons:
  chrome: stable

install:
- go install github.com/agnivade/wasmbrowsertest@latest
- mv $GOPATH/bin/wasmbrowsertest $GOPATH/bin/go_js_wasm_exec
- export PATH=$GOPATH/bin:$PATH

script:
- GOOS=js GOARCH=wasm go test ./wasm_tests/...
```

## Adding New WASM Tests

### For New Methods (as per ADD_NEW_METHODS.md)

When implementing new methods, follow this pattern:

**1. Shared Logic Methods** (no WASM-specific tests needed):
- `MinutesToTime` - Test in regular test suite
- `TimeToMinutes` - Test in regular test suite  
- `IsPast` / `IsFuture` - Test in regular test suite
- `DaysBetween` - Test in regular test suite

**2. WASM-Specific Methods** (add to `wasm_test.go`):
- `UnixNanoToDate` - Requires JS Date API validation
- `DateToUnix` - Requires JS Date parsing validation
- `DateTimeToUnix` - Requires JS Date constructor validation
- `IsToday` - Requires JS timezone handling validation

### Example: Adding UnixNanoToDate Test

```go
//go:build js && wasm
// +build js,wasm

package wasm_test

func TestWasmUnixNanoToDate(t *testing.T) {
    tp := tinytime.NewTimeProvider()
    
    // Test with known timestamp
    nano := int64(1705315200000000000) // 2024-01-15 00:00:00 UTC
    result := tp.UnixNanoToDate(nano)
    expected := "2024-01-15"
    
    if result != expected {
        t.Errorf("UnixNanoToDate(%d) = %s; want %s", nano, result, expected)
    }
    
    // Test epoch
    result = tp.UnixNanoToDate(0)
    if result != "1970-01-01" {
        t.Errorf("UnixNanoToDate(0) = %s; want 1970-01-01", result)
    }
}
```

## Troubleshooting

### Chrome Not Found

**Error:** `Chrome executable not found`

**Solution:**
```bash
# Linux
sudo apt-get install chromium-browser

# macOS
brew install --cask google-chrome

# Check installation
which google-chrome || which chromium
```

### Environment Variable Limit

**Error:** `total length of command line and environment variables exceeds limit`

**Solution:** Use `cleanenv` to reduce environment variables:

```bash
go install github.com/agnivade/wasmbrowsertest/cmd/cleanenv@latest

# Remove CI variables
cleanenv -remove-prefix GITHUB_ -- GOOS=js GOARCH=wasm go test ./wasm_tests/...
```

### Port Already in Use

**Error:** `bind: address already in use`

**Solution:** `wasmbrowsertest` automatically finds free ports. If this fails, kill existing Chrome instances:

```bash
pkill chrome
pkill chromium
```

### Tests Timeout

**Error:** Tests hang or timeout

**Solution:** Increase timeout:

```bash
GOOS=js GOARCH=wasm go test -timeout 5m ./wasm_tests/...
```

## Performance Considerations

### Benchmark in Browser

Run benchmarks with:

```bash
GOOS=js GOARCH=wasm go test -bench=. ./wasm_tests/...
```

**Expected Performance (Chrome):**
- `UnixNano()` - ~50-100 ns/op
- `UnixSecondsToDate()` - ~2-5 µs/op
- `UnixNanoToTime()` - ~2-5 µs/op

### Optimization Tips

1. **Minimize JS API calls** - Cache Date constructor
2. **Reuse buffers** - Use tinystring.Convert() pooling
3. **Batch operations** - Reduce syscall/js overhead

## Best Practices

### 1. Test Coverage

**Do test in browser:**
- Methods using `syscall/js`
- Date API integration
- Timezone handling
- Browser environment detection

**Don't duplicate in browser:**
- Pure math functions (e.g., `MinutesToTime`)
- String parsing (e.g., `TimeToMinutes`)
- Simple comparisons (e.g., `IsPast`)

### 2. Shared Test Data

Keep test data separate from build tags:

```go
// data_test.go (no build tags - shared)
package tinytime_test

var TestUnixSeconds = int64(1624397134)
```

```go
// wasm_test.go (WASM-specific)
//go:build js && wasm

package wasm_test

// Duplicate constants (cannot access data_test.go)
const testUnixSeconds = int64(1624397134)
```

**Why?** Build tags prevent cross-package access.

### 3. Cross-Implementation Validation

Always test consistency between backend and WASM:

```go
func TestWasmConsistency(t *testing.T) {
    tp := tinytime.NewTimeProvider()
    
    // Same input should produce same output
    nano := int64(1705315200000000000)
    wasmResult := tp.UnixNanoToDate(nano)
    
    // Expected result from backend implementation
    expected := "2024-01-15"
    
    if wasmResult != expected {
        t.Errorf("WASM result differs from backend: %s != %s", wasmResult, expected)
    }
}
```

## Resources

- [wasmbrowsertest GitHub](https://github.com/agnivade/wasmbrowsertest)
- [Go WebAssembly Wiki](https://github.com/golang/go/wiki/WebAssembly)
- [ChromeDP Protocol](https://chromedevtools.github.io/devtools-protocol/)
- [TinyTime ADD_NEW_METHODS.md](./issues/ADD_NEW_METHODS.md) - See "Testing Requirements" section

## FAQ

**Q: Can I use Firefox for testing?**  
A: No, `wasmbrowsertest` uses ChromeDP protocol which requires Chrome/Chromium.

**Q: Do I need to install anything else?**  
A: Only Chrome/Chromium browser and the `wasmbrowsertest` binary.

**Q: Can I run `go run` in browser?**  
A: Yes! `GOOS=js GOARCH=wasm go run main.go` also works.

**Q: How do I see the browser UI?**  
A: Set `WASM_HEADLESS=off` environment variable.

**Q: Can I take CPU profiles?**  
A: Yes, use `-cpuprofile` flag. Profile is converted to pprof format automatically.

**Q: Why are some tests duplicated?**  
A: Build tags prevent sharing test code between native and WASM. Keep shared logic in production code, not tests.

---

**Related Documentation:**
- [ADD_NEW_METHODS.md](./issues/ADD_NEW_METHODS.md) - New methods implementation plan
- [TinyTime README](../README.md) - General usage guide
