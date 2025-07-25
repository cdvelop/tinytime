//go:build !wasm
// +build !wasm

package tinytime

import (
	"time"

	"github.com/cdvelop/tinystring"
)

// NewTimeProvider retorna la implementación correcta según el entorno de compilación.
func NewTimeProvider() TimeProvider {
	return timeServer{}
}

// timeServer implementa TimeProvider para Go estándar
type timeServer struct{}

func (timeServer) UnixNano() int64 {
	return time.Now().UnixNano()
}

func (timeServer) UnixSecondsToDate(unixSeconds int64) string {
	t := time.Unix(unixSeconds, 0)
	return t.Format("2006-01-02 15:04")
}

func (timeServer) UnixNanoToTime(input any) string {
	var unixNano int64
	switch v := input.(type) {
	case int64:
		unixNano = v
	case int:
		unixNano = int64(v)
	case float64:
		unixNano = int64(v)
	case string:
		val, err := tinystring.Convert(v).Int64()
		if err != nil {
			return ""
		}
		unixNano = val
	default:
		return ""
	}
	unixSeconds := unixNano / 1e9
	t := time.Unix(unixSeconds, 0)
	return t.Format("15:04:05")
}
