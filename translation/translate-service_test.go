package translation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_translateService_Translate(t *testing.T) {
	t.Run("Translate", func(t *testing.T) {
		cache := NewMockCache(t)
		cache.EXPECT().Get(mock.Anything, mock.Anything).Return("value", true, nil)

		aiClient := NewMockAIClient(t)

		ts := NewTranslateService(cache, aiClient)
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

		translatedItems, err := ts.Translate(t.Context(), data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(translatedItems))
	})
}
