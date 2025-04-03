package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type TranslationItem struct {
	Container any
	Key       any
	Value     string
	Path      string
}

type Collector interface {
	CollectTranslationItemsJson(root any) ([]TranslationItem, error)
}

type collector struct {
	whiteList []string
}

func NewCollector(whiteList []string) Collector {
	return &collector{
		whiteList: whiteList,
	}
}

func filterByKeyword(target string, list []string) bool {
	for _, keyword := range list {
		if target == keyword {
			return true
		}
	}
	return false
}

func (c *collector) CollectTranslationItemsJson(root any) ([]TranslationItem, error) {
	return c.CollectTranslationItem(root, nil, nil, ""), nil
}

func (c *collector) CollectTranslationItem(value any, container any, key any, path string) []TranslationItem {
	translationItems := []TranslationItem{}

	switch typedValue := value.(type) {
	// base case
	case string:
		if keyString, ok := key.(string); ok {
			if !filterByKeyword(keyString, c.whiteList) {
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
				c.CollectTranslationItem(v, typedValue, k, newPath)...,
			)
		}
	case []any:
		for i, v := range typedValue {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			if reflect.TypeOf(v).Kind() != reflect.Map {
				if !filterByKeyword(key.(string), c.whiteList) {
					return translationItems
				}
			}
			translationItems = append(
				translationItems,
				c.CollectTranslationItem(v, typedValue, i, newPath)...,
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
