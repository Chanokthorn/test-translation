package translation

import (
	"context"
	"log"
	"time"
)

type TranslatePayloadItem struct {
	Path string `json:"path"`
	Text string `json:"text"`
}

type AIClient interface {
	Translate(ctx context.Context, data TranslatePayloadItem) (TranslatePayloadItem, error)
	TranslateBatch(ctx context.Context, data []TranslatePayloadItem) ([]TranslatePayloadItem, error)
}

type aiClient struct {
}

func NewAIClient() AIClient {
	return &aiClient{}
}

func (a *aiClient) Translate(ctx context.Context, data TranslatePayloadItem) (TranslatePayloadItem, error) {
	log.Printf("AIClient Translate called")

	// Simulate translation
	data.Text = "Translated: " + data.Text
	time.Sleep(100 * time.Millisecond) // Simulate network delay

	return data, nil
}

func (a *aiClient) TranslateBatch(ctx context.Context, data []TranslatePayloadItem) ([]TranslatePayloadItem, error) {
	log.Printf("AIClient Translate called")

	for i := range len(data) {
		data[i].Text = "Translated: " + data[i].Text
	}

	return data, nil
}
