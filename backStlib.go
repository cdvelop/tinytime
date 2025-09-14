//go:build !wasm
// +build !wasm

package tinytime

import (
	"time"

	"github.com/cdvelop/tinystring"
)

// NewTimeProvider retorna la implementación correcta según el entorno de compilación.
func NewTimeProvider() TimeProvider {
	return &timeServer{}
}

// timeServer implementa TimeProvider para Go estándar
type timeServer struct{}

func (timeServer) UnixNano() int64 {
	return time.Now().UnixNano()
}

func (timeServer) UnixSecondsToDate(unixSeconds int64) string {
	return time.Unix(unixSeconds, 0).UTC().Format("2006-01-02 15:04")
}

func (timeServer) UnixNanoToTime(input any) string {

	unixNano, err := tinystring.Convert(input).Int64()
	if err != nil {
		return ""
	}

	unixSeconds := unixNano / 1e9

	return time.Unix(unixSeconds, 0).Format("15:04:05")
}
