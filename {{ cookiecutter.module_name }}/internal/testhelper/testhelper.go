package testhelper

import (
	"fmt"
	"testing"
)

func MustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Errorf("panic expected but did not occur")
}

func MapsEq[K comparable, V any](m1 map[K]V, m2 map[K]V) bool {
	return fmt.Sprint(m1) == fmt.Sprint(m2)
}
