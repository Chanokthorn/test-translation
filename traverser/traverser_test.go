package traverser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type TranslationItem struct {
	Container any
	Key       any
	Value     string
	Path      string
}

func CollectTranslationItemsJson(root any) ([]TranslationItem, error) {
	return CollectTranslationItem(root, nil, nil, ""), nil
}

func CollectTranslationItem(value any, container any, key any, path string) []TranslationItem {
	translationItems := []TranslationItem{}

	switch typedValue := value.(type) {
	// base case
	case string:
		translationItems = append(translationItems, TranslationItem{
			Container: container,
			Key:       key,
			Value:     typedValue,
			Path:      path,
		})

	// common case
	case map[string]any:
		for k, v := range typedValue {
			newPath := fmt.Sprintf("%s.%s", path, k)
			translationItems = append(
				translationItems,
				CollectTranslationItem(v, typedValue, k, newPath)...,
			)
		}
	case []any:
		for i, v := range typedValue {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			translationItems = append(
				translationItems,
				CollectTranslationItem(v, typedValue, i, newPath)...,
			)
		}
	default:
		break
	}

	return translationItems
}

func Test_object(t *testing.T) {
	t.Run("Test object", func(t *testing.T) {
		data, err := os.ReadFile("small_data.json")
		if err != nil {
			t.Fatalf("failed to read data.json: %v", err)
		}

		var root any
		if err := json.Unmarshal(data, &root); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		TranslationItems, err := CollectTranslationItemsJson(root)
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

		// print translated data
		spew.Dump(root)
	})
}
