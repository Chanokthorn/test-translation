package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func Test_collector_CollectTranslationItem(t *testing.T) {

	t.Run("Test object", func(t *testing.T) {
		whiteList := []string{"title", "h1", "description"}
		c := NewCollector(whiteList)

		data, err := os.ReadFile("small_data.json")
		if err != nil {
			t.Fatalf("failed to read data.json: %v", err)
		}

		var root any
		if err := json.Unmarshal(data, &root); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		writeMapToFile(root, "collector_test_before.json")

		TranslationItems, err := c.CollectTranslationItemsJson(root)
		if err != nil {
			t.Fatalf("failed to collect translation items: %v", err)
		}

		for _, item := range TranslationItems {
			fmt.Printf("Path: %s\n",
				item.Path)
		}

		// simulate translation
		for _, item := range TranslationItems {
			translatedValue := "translated_" + item.Value

			if item.Container != nil {
				switch container := item.Container.(type) {
				case map[string]any:
					container[item.Key.(string)] = translatedValue
				case []any:
					if i, ok := item.Key.(int); ok && i < len(container) {
						container[i] = translatedValue
					}
				default:
					t.Fatalf("unsupported container type: %T", container)
				}
			}
		}

		writeMapToFile(root, "collector_test_after.json")
	})
}
