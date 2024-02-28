package validator_test

import (
	"regexp"
	"testing"

	"{{ cookiecutter.module_path }}/internal/validator"
)

func TestNotBlank(t *testing.T) {
	testCases := []struct {
		value string
		want  bool
	}{
		{"foobar", true},
		{"", false},
	}
	for _, tc := range testCases {
		got := validator.NotBlank(tc.value)

		if tc.want != got {
			t.Errorf("validator.NotBlank(%s) | want: %t, got %t", tc.value, tc.want, got)
		}
	}
}

func TestMinRunes(t *testing.T) {
	testCases := []struct {
		value string
		min   int
		want  bool
	}{
		{"foobar", 6, true},
		{"foobar", 5, true},
		{"foobar", 7, false},
	}
	for _, tc := range testCases {
		got := validator.MinRunes(tc.value, tc.min)

		if tc.want != got {
			t.Errorf("validator.MinRunes(%s, %d) | want: %t, got %t", tc.value, tc.min, tc.want, got)
		}
	}
}

func TestMaxRunes(t *testing.T) {
	testCases := []struct {
		value string
		min   int
		want  bool
	}{
		{"foobar", 6, true},
		{"foobar", 5, false},
		{"foobar", 7, true},
	}
	for _, tc := range testCases {
		got := validator.MaxRunes(tc.value, tc.min)

		if tc.want != got {
			t.Errorf("validator.MaxRunes(%s, %d) | want: %t, got %t", tc.value, tc.min, tc.want, got)
		}
	}
}

func TestBetween(t *testing.T) {
	t.Run("IntCases", func(t *testing.T) {
		intTestCases := []struct {
			value int
			min   int
			max   int
			want  bool
		}{
			{2, 1, 3, true},
			{1, 1, 3, true},
			{3, 1, 3, true},
			{0, 1, 3, false},
			{4, 1, 3, false},
		}
		for _, tc := range intTestCases {
			got := validator.Between(tc.value, tc.min, tc.max)

			if tc.want != got {
				t.Errorf("validator.Between(%d, %d, %d) | want: %t, got %t", tc.value, tc.min, tc.max, tc.want, got)
			}
		}
	})

	t.Run("FloatCases", func(t *testing.T) {
		floatTestCases := []struct {
			value float32
			min   float32
			max   float32
			want  bool
		}{
			{1.5, 1.0, 2.0, true},
			{1.000, 1.0, 2.0, true},
			{1.001, 1.0, 2.0, true},
			{2.000, 1.0, 2.0, true},
			{2.001, 1.0, 2.0, false},
			{0.999, 1.0, 2.0, false},
		}
		for _, tc := range floatTestCases {
			got := validator.Between(tc.value, tc.min, tc.max)

			if tc.want != got {
				t.Errorf("validator.Between(%f, %f, %f) | want: %t, got %t", tc.value, tc.min, tc.max, tc.want, got)
			}
		}
	})
}

func TestMatches(t *testing.T) {
	testCases := []struct {
		value string
		want  bool
	}{
		{"color", true},
		{"colour", true},
		{"coulor", false},
		{"", false},
	}
	for _, tc := range testCases {
		p := "colou?r"
		rx := regexp.MustCompile(p)
		got := validator.Matches(tc.value, rx)

		if tc.want != got {
			t.Errorf("validator.Matches(%s, %s) | want: %t, got %t", tc.value, p, tc.want, got)
		}
	}
}

func TestIn(t *testing.T) {
	t.Run("IntCases", func(t *testing.T) {
		intTestCases := []struct {
			value    int
			safelist []int
			want     bool
		}{
			{2, []int{1, 2, 3}, true},
			{0, []int{1, 2, 3}, false},
		}
		for _, tc := range intTestCases {
			got := validator.In(tc.value, tc.safelist...)

			if tc.want != got {
				t.Errorf("validator.In(%d, %v) | want: %t, got %t", tc.value, tc.safelist, tc.want, got)
			}
		}
	})

	t.Run("FloatCases", func(t *testing.T) {
		floatTestCases := []struct {
			value    float32
			safelist []float32
			want     bool
		}{
			{2.0, []float32{1.0, 2.0, 3.0}, true},
			{2.00000, []float32{1.0, 2.0, 3.0}, true},
			{2.01, []float32{1.0, 2.0, 3.0}, false},
			{3.01, []float32{1.0, 2.0, 3.0}, false},
		}
		for _, tc := range floatTestCases {
			got := validator.In(tc.value, tc.safelist...)

			if tc.want != got {
				t.Errorf("validator.In(%f, %v) | want: %t, got %t", tc.value, tc.safelist, tc.want, got)
			}
		}
	})
}

func TestNotIn(t *testing.T) {
	t.Run("IntCases", func(t *testing.T) {
		intTestCases := []struct {
			value    int
			safelist []int
			want     bool
		}{
			{2, []int{1, 2, 3}, false},
			{0, []int{1, 2, 3}, true},
		}
		for _, tc := range intTestCases {
			got := validator.NotIn(tc.value, tc.safelist...)

			if tc.want != got {
				t.Errorf("validator.NotIn(%d, %v) | want: %t, got %t", tc.value, tc.safelist, tc.want, got)
			}
		}
	})

	t.Run("FloatCases", func(t *testing.T) {
		floatTestCases := []struct {
			value    float32
			safelist []float32
			want     bool
		}{
			{2.0, []float32{1.0, 2.0, 3.0}, false},
			{2.00000, []float32{1.0, 2.0, 3.0}, false},
			{2.01, []float32{1.0, 2.0, 3.0}, true},
			{3.01, []float32{1.0, 2.0, 3.0}, true},
		}
		for _, tc := range floatTestCases {
			got := validator.NotIn(tc.value, tc.safelist...)

			if tc.want != got {
				t.Errorf("validator.NotIn(%f, %v) | want: %t, got %t", tc.value, tc.safelist, tc.want, got)
			}
		}
	})
}

func TestAllIn(t *testing.T) {
	t.Run("IntCases", func(t *testing.T) {
		intTestCases := []struct {
			values   []int
			safelist []int
			want     bool
		}{
			{[]int{1}, []int{1, 2, 3}, true},
			{[]int{1, 3}, []int{1, 2, 3}, true},
			{[]int{1, 4}, []int{1, 2, 3}, false},
			{[]int{0, 4}, []int{1, 2, 3}, false},
		}
		for _, tc := range intTestCases {
			got := validator.AllIn(tc.values, tc.safelist...)

			if tc.want != got {
				t.Errorf("validator.AllIn(%v, %v) | want: %t, got %t", tc.values, tc.safelist, tc.want, got)
			}
		}
	})

	t.Run("FloatCases", func(t *testing.T) {
		floatTestCases := []struct {
			values   []float32
			safelist []float32
			want     bool
		}{
			{[]float32{2.0}, []float32{1.0, 2.0, 3.0}, true},
			{[]float32{1.0, 2.0}, []float32{1.0, 2.0, 3.0}, true},
			{[]float32{1.0, 2.1}, []float32{1.0, 2.0, 3.0}, false},
			{[]float32{0.1, 2.1}, []float32{1.0, 2.0, 3.0}, false},
		}
		for _, tc := range floatTestCases {
			got := validator.AllIn(tc.values, tc.safelist...)

			if tc.want != got {
				t.Errorf("validator.AllIn(%v, %v) | want: %t, got %t", tc.values, tc.safelist, tc.want, got)
			}
		}
	})
}

func TestNoDuplicates(t *testing.T) {
	t.Run("IntCases", func(t *testing.T) {
		intTestCases := []struct {
			values []int
			want   bool
		}{
			{[]int{1, 2, 3}, true},
			{[]int{1, 2, 2}, false},
		}
		for _, tc := range intTestCases {
			got := validator.NoDuplicates(tc.values)

			if tc.want != got {
				t.Errorf("validator.NoDuplicates(%v) | want: %t, got %t", tc.values, tc.want, got)
			}
		}
	})

	t.Run("FloatCases", func(t *testing.T) {
		floatTestCases := []struct {
			values []float32
			want   bool
		}{
			{[]float32{1.0, 2.0, 3.0}, true},
			{[]float32{1.0, 2.0, 2.0}, false},
		}
		for _, tc := range floatTestCases {
			got := validator.NoDuplicates(tc.values)

			if tc.want != got {
				t.Errorf("validator.NoDuplicates(%v) | want: %t, got %t", tc.values, tc.want, got)
			}
		}
	})

	t.Run("StringCases", func(t *testing.T) {
		strTestCases := []struct {
			values []string
			want   bool
		}{
			{[]string{"foo", "bar", "baz"}, true},
			{[]string{"foo", "bar", "bar"}, false},
		}
		for _, tc := range strTestCases {
			got := validator.NoDuplicates(tc.values)

			if tc.want != got {
				t.Errorf("validator.NoDuplicates(%v) | want: %t, got %t", tc.values, tc.want, got)
			}
		}
	})
}

func TestIsEmail(t *testing.T) {
	testCases := []struct {
		email string
		want  bool
	}{
		{"foo@bar.com", true},
		{"foo@bar", true},
		{"foo@", false},
		{"foo", false},
		{"", false},
	}
	for _, tc := range testCases {
		got := validator.IsEmail(tc.email)

		if tc.want != got {
			t.Errorf("validator.IsEmail(%s) | want: %t, got %t", tc.email, tc.want, got)
		}
	}
}

func TestIsURL(t *testing.T) {
	testCases := []struct {
		url  string
		want bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"http://example.com/foo/bar", true},
		{"https://example.com/foo?bar=baz", true},
		{"https://john@example.com", true},
		{"ssh://example.com", true},
		{"ftp://example.com", true},
		{"example.com", false},
		{"example.com/foo/bar", false},
		{"example.com/foo?bar=baz", false},
		{"http://example", true},
		{"https://example", true},
		{"example", false},
		{"john@example.com", false},
		{"", false},
	}
	for _, tc := range testCases {
		got := validator.IsURL(tc.url)

		if tc.want != got {
			t.Errorf("validator.IsURL(%s) | want: %t, got %t", tc.url, tc.want, got)
		}
	}
}
