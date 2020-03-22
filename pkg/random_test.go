package pkg

import "testing"

func TestGenerateRandomString(t *testing.T) {
	t.Run("Non-positive length", func(t *testing.T) {
		if got := GenerateRandomString(0); len(got) != 0 {
			t.Errorf("GenerateRandomString() = %v, want %v", got, "")
		}

		if got := GenerateRandomString(-1); len(got) != 0 {
			t.Errorf("GenerateRandomString() = %v, want %v", got, "")
		}
	})

	t.Run("Randomness", func(t *testing.T) {
		for i := 1; i < 1000; i++ {
			got1 := GenerateRandomString(i)
			got2 := GenerateRandomString(i)

			if got1 == got2 {
				t.Errorf("GenerateRandomString() = %v and %v, should differ", got1, got2)
			}
		}
	})
}
