package testhelper

import (
	"testing"
)

func MustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Errorf("panic expected but did not occur")
}
