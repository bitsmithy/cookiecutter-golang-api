package validator_test

import (
	"testing"

	"{{ cookiecutter.module_path }}/internal/validator"
)

func TestAddError(t *testing.T) {
	v := validator.New()

	if v.HasErrors() {
		t.Error("v.HasErrors() | want: false, got: true")
	}

	v.AddError("this is an error")

	if !v.HasErrors() {
		t.Error("v.HasErrors() | want: true, got: false")
	}
}

func TestAddFieldError(t *testing.T) {
	v := validator.New()

	if v.HasErrors() {
		t.Error("v.HasErrors() | want: false, got: true")
	}

	v.AddFieldError("foo", "an error on foo")

	if !v.HasErrors() {
		t.Error("v.HasErrors() | want: true, got: false")
	}
}

func TestCheck(t *testing.T) {
	testCases := []struct {
		check     bool
		hasErrors bool
	}{
		{2+2 == 4, false},
		{2+3 == 4, true},
		{validator.NotBlank(""), true},
	}

	for _, tc := range testCases {
		v := validator.New()

		if v.HasErrors() {
			t.Error("v.HasErrors() | want: false, got: true")
		}

		v.Check(tc.check, "error message")

		if v.HasErrors() != tc.hasErrors {
			t.Errorf("v.HasErrors() | want: %t, got: %t", tc.hasErrors, v.HasErrors())
		}
	}
}

func TestCheckField(t *testing.T) {
	testCases := []struct {
		check     bool
		hasErrors bool
	}{
		{2+2 == 4, false},
		{2+3 == 4, true},
		{validator.NotBlank(""), true},
	}

	for _, tc := range testCases {
		v := validator.New()

		if v.HasErrors() {
			t.Error("v.HasErrors() | want: false, got: true")
		}

		v.CheckField(tc.check, "fieldName", "error message")

		if v.HasErrors() != tc.hasErrors {
			t.Errorf("v.HasErrors() | want: %t, got: %t", tc.hasErrors, v.HasErrors())
		}
	}
}
