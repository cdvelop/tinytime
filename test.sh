#!/bin/bash

echo "=========================================="
echo "Running Backend Tests (Go stdlib)..."
echo "=========================================="
go test -v -cover ./...

if [ $? -ne 0 ]; then
    echo "❌ Backend tests failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "Running WASM Tests (Browser)..."
echo "=========================================="

# Check if wasmbrowsertest is installed
if ! command -v wasmbrowsertest &> /dev/null; then
    echo "⚠️  wasmbrowsertest not found. Install it with:"
    echo "   go install github.com/agnivade/wasmbrowsertest@latest"
    echo "   export PATH=\$PATH:\$(go env GOPATH)/bin"
    exit 1
fi

# Run WASM tests
GOOS=js GOARCH=wasm go test -v -cover ./... 2>&1 | grep -v "ERROR: could not unmarshal"
WASM_EXIT_CODE=$?

if [ $WASM_EXIT_CODE -ne 0 ]; then
    echo ""
    echo "❌ WASM tests failed"
    exit 1
fi

echo ""
echo "✅ All tests passed!"
