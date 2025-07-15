# tinytime

TinyGo-compatible time package for Go, using syscall/js on WASM to keep binaries small and leverage native JS time APIs.


Minimal and portable time utility package for Go and TinyGo with WebAssembly support.
When compiled for wasm (GOOS=js GOARCH=wasm), it uses syscall/js to access native JavaScript time APIs (Date.now, performance.now, etc.), drastically reducing binary size by avoiding the Go standard library.
For non-WASM targets, it falls back to the standard time package.
Ideal for frontend projects where binary size and compatibility with JavaScript environments matter.