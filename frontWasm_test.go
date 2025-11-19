//go:build wasm
// +build wasm

package tinytime_test

import (
	"testing"

	"github.com/cdvelop/tinytime"
)

// Tests for WASM environment using shared validation functions

func TestAllSharedWasm(t *testing.T) {
	tp := tinytime.NewTimeProvider()

	t.Run("FormatDate", func(t *testing.T) { FormatDateShared(t, tp) })
	t.Run("FormatTime", func(t *testing.T) { FormatTimeShared(t, tp) })
	t.Run("FormatDateTime", func(t *testing.T) { FormatDateTimeShared(t, tp) })
	t.Run("ParseDate", func(t *testing.T) { ParseDateShared(t, tp) })
	t.Run("ParseTime", func(t *testing.T) { ParseTimeShared(t, tp) })
	t.Run("ParseDateTime", func(t *testing.T) { ParseDateTimeShared(t, tp) })
	t.Run("IsToday", func(t *testing.T) { IsTodayShared(t, tp) })
	t.Run("IsPast", func(t *testing.T) { IsPastShared(t, tp) })
	t.Run("IsFuture", func(t *testing.T) { IsFutureShared(t, tp) })
	t.Run("DaysBetween", func(t *testing.T) { DaysBetweenShared(t, tp) })

	// Existing tests
	t.Run("UnixNanoToTime", func(t *testing.T) { UnixNanoToTimeShared(t, tp) })
	t.Run("UnixNanoToTimeWithDifferentTypes", func(t *testing.T) { UnixNanoToTimeWithDifferentTypesShared(t, tp) })
	t.Run("UnixNano", func(t *testing.T) { UnixNanoShared(t, tp) })
	t.Run("UnixSecondsToDate", func(t *testing.T) { UnixSecondsToDateShared(t, tp) })
	t.Run("UnixSecondsToDateEdgeCases", func(t *testing.T) { UnixSecondsToDateEdgeCasesShared(t, tp) })
	t.Run("UnixNanoToTimeEdgeCases", func(t *testing.T) { UnixNanoToTimeEdgeCasesShared(t, tp) })

	tp2 := tinytime.NewTimeProvider()
	t.Run("TimeProviderConsistency", func(t *testing.T) {TimeProviderConsistencyShared(t, tp, tp2)})
}
