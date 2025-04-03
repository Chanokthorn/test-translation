package translation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_aiClient_Translate(t *testing.T) {
	t.Run("Translate", func(t *testing.T) {
		ac := NewAIClient()
		data := []TranslatePayloadItem{
			{
				Path: "packages[0].name",
				Text: "package name1",
			},
			{
				Path: "packages[1].name",
				Text: "package name2",
			},
		}

		result, err := ac.TranslateBatch(context.Background(), data)
		require.NoError(t, err)

		require.Equal(t, 2, len(result))
	})
}
