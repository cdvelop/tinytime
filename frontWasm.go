//go:build wasm
// +build wasm

package tinytime

import (
	"syscall/js"

	. "github.com/cdvelop/tinystring"
)

// timeClient implementa TimeProvider para entornos WASM/JS
// usando la API de JavaScript.
type timeClient struct {
	// cachear el constructor Date para evitar lookups repetidos
	dateCtor js.Value
	jsDate   js.Value
	// jsTmp eliminado; usar variables locales para temporales
	// opcional:
	// dateProtoToISO js.Value

	buff *Conv
}

// NewTimeProvider retorna la implementación correcta para WASM.
func NewTimeProvider() TimeProvider {
	// cachear el constructor Date y eliminar jsTmp
	return &timeClient{
		dateCtor: js.Global().Get("Date"),
		// opcional: cachear métodos del prototype:
		// dateProtoToISO: js.Global().Get("Date").Get("prototype").Get("toISOString"),

		buff: Convert(),
	}
}

func (t *timeClient) UnixNano() int64 {
	t.jsDate = t.dateCtor.New()
	msTimestamp := t.jsDate.Call("getTime").Float()
	return int64(msTimestamp * 1e6)
}

func (t *timeClient) UnixSecondsToDate(unixSeconds int64) (date string) {
	// Crea una instancia de Date de JavaScript a partir de los segundos de Unix
	t.jsDate = t.dateCtor.New(float64(unixSeconds) * 1000)

	// Llama al método toISOString y convierte a string directamente
	date = t.jsDate.Call("toISOString").String()

	t.buff.Reset()

	// Formatea la cadena de fecha a "2006-01-02 15:04"
	t.buff.Write(date[0:10])
	t.buff.Write(" ")
	t.buff.Write(date[11:16])

	return t.buff.String()
}

func (t *timeClient) UnixNanoToTime(input any) string {

	unixNano, err := Convert(input).Int64()
	if err != nil {
		return ""
	}

	unixSeconds := unixNano / 1e9

	t.jsDate = t.dateCtor.New(unixSeconds * 1000)

	hours := t.jsDate.Call("getHours").Int()
	minutes := t.jsDate.Call("getMinutes").Int()
	seconds := t.jsDate.Call("getSeconds").Int()
	return Fmt("%02d:%02d:%02d", hours, minutes, seconds)
}
