package translation

import (
	"context"
	"encoding/json"
	"fmt"
)

type Processor interface {
	Translate(ctx context.Context, data []byte) ([]byte, error)
}

type processor struct {
	collector        Collector
	translateService TranslateService
}

func NewProcessor(
	collector Collector,
	translateService TranslateService,
) Processor {
	return &processor{
		collector:        collector,
		translateService: translateService,
	}
}

func (p *processor) Translate(ctx context.Context, data []byte) ([]byte, error) {
	var root any
	// Unmarshal the data into a map
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	translationItems, err := p.collector.CollectTranslationItemsJson(root)
	if err != nil {
		return nil, fmt.Errorf("failed to collect translation items: %w", err)
	}

	toBeTranslatedPayload := make([]TranslatePayloadItem, len(translationItems))
	for i, item := range translationItems {
		toBeTranslatedPayload[i] = TranslatePayloadItem{
			Path: item.Path,
			Text: item.Value,
		}
	}

	translatedPayload, err := p.translateService.Translate(ctx, toBeTranslatedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to translate payload: %w", err)
	}

	for i, item := range translationItems {
		translatedValue := translatedPayload[i].Text

		if item.Container != nil {
			switch container := item.Container.(type) {
			case map[string]any:
				container[item.Key.(string)] = translatedValue
			case []any:
				if i, ok := item.Key.(int); ok && i < len(container) {
					container[i] = translatedValue
				}
			default:
				return nil, fmt.Errorf("unsupported container type: %T", container)
			}
		}
	}

	// Marshal the modified data back to JSON
	translatedData, err := json.Marshal(root)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	return translatedData, nil
}
