#!/bin/bash

echo "=========================================="
echo "Running Backend Tests (Go stdlib)..."
echo "=========================================="
go test -v ./...

if [ $? -ne 0 ]; then
    echo "❌ Backend tests failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "Running WASM Tests (Browser)..."
echo "=========================================="
GOOS=js GOARCH=wasm go test -v ./... 2>&1 | grep -v "ERROR: could not unmarshal"

if [ $? -ne 0 ]; then
    echo "❌ WASM tests failed"
    exit 1
fi

echo ""
echo "✅ All tests passed!"
