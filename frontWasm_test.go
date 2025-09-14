package tinytime_test

import (
	"testing"

	"github.com/cdvelop/wasmtest"
)

// TestWasmIntegration runs all WebAssembly tests using the simplified wasmtest API
func TestWasmIntegration(t *testing.T) {
	// Run WebAssembly tests using the ultra-simple wasmtest API
	if err := wasmtest.RunTests(); err != nil {
		t.Fatal(err)
	}
}
