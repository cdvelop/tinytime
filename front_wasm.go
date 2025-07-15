//go:build wasm
// +build wasm

package tinytime

import (
	"syscall/js"

	"github.com/cdvelop/tinystring"
)

// timeClient implementa TimeProvider para entornos WASM/JS
// usando la API de JavaScript.
type timeClient struct{}

func (timeClient) UnixNano() int64 {
	jsDate := js.Global().Get("Date").New()
	msTimestamp := jsDate.Call("getTime").Float()
	return int64(msTimestamp * 1e6)
}

func (timeClient) UnixSecondsToDate(unixSeconds int64) string {
	jsDate := js.Global().Get("Date").New(float64(unixSeconds) * 1000)
	dateJSValue := jsDate.Call("toISOString")
	date := dateJSValue.String()
	return date[0:10] + " " + date[11:16]
}

func (timeClient) UnixNanoToTime(input any) string {
	var unixNano int64
	switch v := input.(type) {
	case int64:
		unixNano = v
	case int:
		unixNano = int64(v)
	case float64:
		unixNano = int64(v)
	case string:
		parsed := int64(0)
		multiplier := int64(1)
		for i := len(v) - 1; i >= 0; i-- {
			if v[i] >= '0' && v[i] <= '9' {
				parsed += int64(v[i]-'0') * multiplier
				multiplier *= 10
			} else {
				return ""
			}
		}
		unixNano = parsed
	default:
		return ""
	}
	unixSeconds := unixNano / 1e9
	jsDate := js.Global().Get("Date").New(unixSeconds * 1000)
	hours := jsDate.Call("getHours").Int()
	minutes := jsDate.Call("getMinutes").Int()
	seconds := jsDate.Call("getSeconds").Int()
	return tinystring.Fmt("%02d:%02d:%02d", hours, minutes, seconds)
}

// NewTimeProvider retorna la implementación correcta para WASM.
func NewTimeProvider() TimeProvider {
	return timeClient{}
}
