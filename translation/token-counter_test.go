package translation

import (
	"testing"
)

func Test_tokenCounter_CountTokens(t *testing.T) {
	t.Run("CountTokens", func(t *testing.T) {
		tokenCounter := NewTokenCounter("../tokenizer/tokenizer.py")
		texts := []string{
			"some text",
			"some other text",
			"some other text that is very very long i don't even know i would type something this long omg",
		}
		expectedTokens := 2 + 3 + 21 // Adjust this based on the expected token count

		tokens, err := tokenCounter.CountTokens(t.Context(), texts)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if tokens != expectedTokens {
			t.Fatalf("Expected %d tokens, got %d", expectedTokens, tokens)
		}
	})
}
