package tinytime

// TimeProvider define la interfaz para utilidades de tiempo
// implementada tanto para Go estándar como para WASM/JS.
type TimeProvider interface {
	UnixNano() int64
	UnixSecondsToDate(int64) string
	UnixNanoToTime(any) string
}
