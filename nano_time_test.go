package tinytime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cdvelop/tinytime"
)

func TestUnixNanoToTime(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	// Test con timestamp conocido en la zona horaria local
	testUnixSeconds := int64(1624397134) // 2021-06-22 15:32:14 UTC
	expectedTime := time.Unix(testUnixSeconds, 0)
	expected := expectedTime.Format("15:04:05")

	nanoTimestamp := testUnixSeconds * 1e9 // convertir a nanosegundos

	result := tp.UnixNanoToTime(nanoTimestamp)
	if result != expected {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", nanoTimestamp, result, expected)
	}

	// Test con string
	result = tp.UnixNanoToTime((fmt.Sprintf("%d", nanoTimestamp)))
	if result != expected {
		t.Errorf("UnixNanoToTime(string) = %s; want %s", result, expected)
	}

	// Test con timestamps secuenciales para verificar orden
	now := time.Now()
	baseNano := now.UnixNano()

	var results []string
	for i := 0; i < 3; i++ {
		nano := baseNano + int64(i)*int64(time.Second) // Incrementar 1 segundo
		timeStr := tp.UnixNanoToTime(nano)
		results = append(results, timeStr)
		t.Logf("Nano: %d -> Time: %s", nano, timeStr)
	}

	// Verificar que los tiempos están en orden
	for i := 1; i < len(results); i++ {
		if results[i] <= results[i-1] {
			t.Errorf("Los timestamps no están en orden: %s <= %s", results[i], results[i-1])
		}
	}
}

func TestUnixNanoToTimeWithDifferentTypes(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	now := time.Now()
	nanoTimestamp := now.UnixNano()

	// Test con diferentes tipos de entrada
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", nanoTimestamp},
		{"int", int(nanoTimestamp)},
		{"float64", float64(nanoTimestamp)},
		{"string", (fmt.Sprintf("%d", nanoTimestamp))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tp.UnixNanoToTime(tc.input)

			if result == "" {
				t.Errorf("UnixNanoToTime devolvió string vacío para tipo %s", tc.name)
			}

			t.Logf("Tipo %s: %v -> %s", tc.name, tc.input, result)
		})
	}

	// Test con tipo no soportado
	invalidResult := tp.UnixNanoToTime(make(chan int))
	if invalidResult != "" {
		t.Error("UnixNanoToTime debería devolver string vacío para tipos no soportados")
	}
}
