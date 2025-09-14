package tinytime_test

import (
	"testing"
	"time"

	"github.com/cdvelop/wasmtest"
)

// TestWasmIntegration runs all WebAssembly tests using the simplified wasmtest API
func TestWasmIntegration(t *testing.T) {
	// Run WebAssembly tests using the ultra-simple wasmtest API
	if err := wasmtest.RunTests("./wasm_test", func(a ...any) { t.Log(a) }, 10*time.Minute); err != nil {
		t.Errorf("WebAssembly tests failed: %v", err)
	}
}
