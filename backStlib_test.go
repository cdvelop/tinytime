package tinytime_test

import (
	"testing"

	"github.com/cdvelop/tinytime"
)

// Tests for standard Go environment using shared validation functions

func TestUnixNanoToTime(t *testing.T) {
	tp := tinytime.NewTimeProvider()
	UnixNanoToTimeShared(t, tp)
}

func TestUnixNanoToTimeWithDifferentTypes(t *testing.T) {
	tp := tinytime.NewTimeProvider()
	UnixNanoToTimeWithDifferentTypesShared(t, tp)
}

func TestUnixNanoFunction(t *testing.T) {
	tp := tinytime.NewTimeProvider()
	UnixNanoShared(t, tp)
}

func TestUnixSecondsToDate(t *testing.T) {
	tp := tinytime.NewTimeProvider()
	UnixSecondsToDateShared(t, tp)
}

func TestUnixSecondsToDateEdgeCases(t *testing.T) {
	tp := tinytime.NewTimeProvider()
	UnixSecondsToDateEdgeCasesShared(t, tp)
}

func TestUnixNanoToTimeEdgeCases(t *testing.T) {
	tp := tinytime.NewTimeProvider()
	UnixNanoToTimeEdgeCasesShared(t, tp)
}

func TestTimeProviderConsistency(t *testing.T) {
	tp1 := tinytime.NewTimeProvider()
	tp2 := tinytime.NewTimeProvider()
	TimeProviderConsistencyShared(t, tp1, tp2)
}
