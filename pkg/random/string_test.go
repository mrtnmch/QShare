package random

import "testing"

func TestGenerateRandomString(t *testing.T) {
	t.Run("Non-positive length", func(t *testing.T) {
		generator := NewStringGenerator()

		if got := generator(0); len(got) != 0 {
			t.Errorf("generator() = %v, want %v", got, "")
		}

		if got := generator(-1); len(got) != 0 {
			t.Errorf("generator() = %v, want %v", got, "")
		}
	})

	t.Run("Randomness", func(t *testing.T) {
		generator := NewStringGenerator()

		for i := 1; i < 1_000; i++ {
			got1 := generator(i)
			got2 := generator(i)

			if got1 == got2 {
				t.Errorf("generator() = %v and %v, should differ", got1, got2)
			}
		}
	})
}
