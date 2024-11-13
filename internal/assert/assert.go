package assert

import (
	"strings"
	"testing"
)

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	// Note: The t.Helper() function that we’re using in the code
	// above indicates to the Go test runner that our Equal() function
	// is a test helper. This means that when t.Errorf() is called
	// from our Equal() function, the Go test runner will report the
	// filename and line number of the code which called our Equal()
	// function in the output.

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
