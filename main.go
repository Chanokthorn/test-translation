package main

import (
	"encoding/json"
	"io"
	"net/http"
	"translate_service_poc/translation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := echo.New()

	rd := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cache := translation.NewCache(rd)
	aiClient := translation.NewAIClient()
	translateService := translation.NewTranslateService(cache, aiClient)
	collector := translation.NewCollector([]string{"title", "h1", "category_id", "listing_seo"})
	processor := translation.NewProcessor(collector, translateService)

	handler := &Handler{
		processor: processor,
	}

	e.POST("/", handler.HandleTranslate)

	e.Logger.Fatal(e.Start(":3000"))
}

type Handler struct {
	processor translation.Processor
}

func (h *Handler) HandleTranslate(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	res, err := h.processor.Translate(c.Request().Context(), body)
	if err != nil {
		log.Errorf("failed to process translation: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to process translation"})
	}

	var jsonRes any
	if err := json.Unmarshal(res, &jsonRes); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to unmarshal response"})
	}

	return c.JSON(http.StatusOK, jsonRes)
}
