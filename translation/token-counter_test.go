package translation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		accCount, err := tokenCounter.GetAccCount()
		require.NoError(t, err)
		assert.Equal(t, expectedTokens, accCount)
	})
}
