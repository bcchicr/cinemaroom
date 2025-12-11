package aggregates

import "testing"

func assertOptionalEqual[T any](t *testing.T, got, want *T, name string, eq func(a, b *T) bool) {
	t.Helper()
	switch {
	case got == nil && want == nil:
		return
	case got == nil || want == nil:
		t.Fatalf("%s: got %v, want %v", name, got, want)
	default:
		if !eq(got, want) {
			t.Fatalf("%s: got %v, want %v", name, got, want)
		}
	}
}
