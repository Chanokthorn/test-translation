package translation

import (
	"context"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"log"
)

type TranslateService interface {
	Translate(ctx context.Context, payload []TranslatePayloadItem) ([]TranslatePayloadItem, error)
}

type translateService struct {
	cache    Cache
	aiClient AIClient
}

func NewTranslateService(cache Cache, aiClient AIClient) TranslateService {
	return &translateService{
		cache:    cache,
		aiClient: aiClient,
	}
}

func hash(s string) string {
	h := fnv.New64()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (t *translateService) Translate(ctx context.Context, payloads []TranslatePayloadItem) ([]TranslatePayloadItem, error) {
	translatedPayloads := make([]TranslatePayloadItem, len(payloads))
	for i, payload := range payloads {
		hashedValue := hash(payload.Text)

		cachedValue, found, err := t.cache.Get(ctx, hashedValue)
		if err != nil {
			return nil, fmt.Errorf("failed to get value from cache: %w", err)
		}

		if found {
			translatedPayloads[i] = TranslatePayloadItem{
				Path: payload.Path,
				Text: cachedValue.(string),
			}
			log.Println("path:", translatedPayloads[i].Path)
			log.Println("text:", translatedPayloads[i].Text)
			continue
		}

		translatedPayload, err := t.aiClient.Translate(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("failed to translate payload: %w", err)
		}

		translatedPayloads[i] = translatedPayload
		if err := t.cache.Set(ctx, hashedValue, translatedPayload.Text); err != nil {
			return nil, fmt.Errorf("failed to set value in cache: %w", err)
		}
	}

	return translatedPayloads, nil
}
