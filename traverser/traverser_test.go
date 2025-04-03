package traverser

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

func filterByKeyword(target string) bool {
	for _, keyword := range []string{"title", "h1", "description"} {
		if target == keyword {
			return true
		}
	}
	return false
}

func CollectTranslationItem(value any, container any, key any, path string) []TranslationItem {
	translationItems := []TranslationItem{}

	switch typedValue := value.(type) {
	// base case
	case string:
		if keyString, ok := key.(string); ok {
			if !filterByKeyword(keyString) {
				return translationItems
			}
		}
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
			if reflect.TypeOf(v).Kind() != reflect.Map {
				if !filterByKeyword(key.(string)) {
					return translationItems
				}
			}
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

func writeMapToFile(m any, filename string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
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

		writeMapToFile(root, "before.json")

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
		writeMapToFile(root, "after.json")
	})
}
