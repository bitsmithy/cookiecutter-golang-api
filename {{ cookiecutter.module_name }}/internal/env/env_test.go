package env_test

import (
	"testing"

	"{{ cookiecutter.module_path }}/internal/env"

	th "{{ cookiecutter.module_path }}/internal/testhelper"
)

func TestGetString(t *testing.T) {
	t.Setenv("FOO", "bar")
	t.Setenv("BAZ", "1")
	defaultVal := "default"

	testCases := []struct {
		key  string
		want string
	}{
		{"FOO", "bar"},
		{"BAZ", "1"},
		{"QUX", defaultVal},
	}
	for _, tc := range testCases {
		got := env.GetString(tc.key, defaultVal)

		if tc.want != got {
			t.Errorf("env.GetString(%s, %s) | want: %s, got %s", tc.key, defaultVal, tc.want, got)
		}
	}
}

func TestGetInt(t *testing.T) {
	t.Setenv("FOO", "12")
	defaultVal := 1

	testCases := []struct {
		key  string
		want int
	}{
		{"FOO", 12},
		{"BAZ", defaultVal},
	}
	for _, tc := range testCases {
		got := env.GetInt(tc.key, defaultVal)

		if tc.want != got {
			t.Errorf("env.GetInt(%s, %d) | want: %d, got %d", tc.key, defaultVal, tc.want, got)
		}
	}

	t.Run("NotParseableInt", func(t *testing.T) {
		t.Setenv("FOO", "not_an_int")
		th.MustPanic(t, func() {
			env.GetInt("FOO", 1)
		})
	})
}

func TestGetBool(t *testing.T) {
	t.Setenv("FOO", "true")
	t.Setenv("BAR", "FALSE")
	t.Setenv("BAZ", "1")
	t.Setenv("QUX", "0")
	defaultVal := true

	testCases := []struct {
		key  string
		want bool
	}{
		{"FOO", true},
		{"BAR", false},
		{"BAZ", true},
		{"QUX", false},
		{"NONEXISTENT", defaultVal},
	}
	for _, tc := range testCases {
		got := env.GetBool(tc.key, defaultVal)

		if tc.want != got {
			t.Errorf("env.GetBool(%s, %t) | want: %t, got %t", tc.key, defaultVal, tc.want, got)
		}
	}

	t.Run("NotParseableBool", func(t *testing.T) {
		t.Setenv("FOO", "not_a_bool")
		th.MustPanic(t, func() {
			env.GetBool("FOO", true)
		})
	})
}
