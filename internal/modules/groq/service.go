package groq

import (
	"bcv/internal/domain"
	"bcv/internal/modules/news"
	"bcv/internal/modules/trade"
	"bcv/internal/platform/database"
	"bcv/internal/platform/server"
	"bcv/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type forecastOutput struct {
	Action       string   `json:"action"`
	Rationale    string   `json:"rationale"`
	KeyFactors   []string `json:"key_factors"`
	WinPoints    float64  `json:"win_points"`
	AccuracyRate float64  `json:"accuracy_rate"`
}

func GetTradeSignal(app *server.App) error {
	data, err := news.FetchNewsTitles()
	if err != nil {
		slog.Error("Error fetching news titles", "error", err)
	}

	rates, err := database.GetListOfLatestRates(app.DB)
	if err != nil {
		slog.Error("Error fetching latest rates", "error", err)
	}

	binance, err := trade.FetchBinanceRates()
	if err != nil {
		slog.Error("Error fetching Binance rates", "error", err)
		return fmt.Errorf("cannot proceed without binance rates: %w", err)
	}
	res, err := GenerateTradeSignal(app.IA, data, rates, *binance)

	if err != nil {
		slog.Error("Error generating trade signal", "error", err)
		return fmt.Errorf("failed to generate trade signal: %w", err)
	}

	err = database.SaveTradeSignal(app.DB, &res)
	if err != nil {
		slog.Error("Error saving trade signal", "error", err)
		return fmt.Errorf("failed to save trade signal: %w", err)
	}

	slog.Info("Generated and saved trade signal successfully", "action", res.Action, "win_points", res.WinPoints, "accuracy_rate", res.AccuracyRate)
	return nil
}

func GenerateTradeSignal(llm *openai.LLM, newsTitles []string, rates []models.CurrencyRate, binanceRate domain.BinanceRate) (models.TradeSignal, error) {
	type r struct {
		Symbol string `json:"s"`
		Price  string `json:"p"`
	}

	var compactRates []r
	for _, rr := range rates {
		compactRates = append(compactRates, r{Symbol: rr.Symbol, Price: fmt.Sprintf(`%.2f`, rr.Price)})
	}

	if len(newsTitles) == 0 {
		return models.TradeSignal{}, fmt.Errorf("no news titles available for analysis")
	}
	if binanceRate.BuyPrice <= 0 || binanceRate.SellPrice <= 0 {
		return models.TradeSignal{}, fmt.Errorf("invalid Binance rate data")
	}

	numRegex := regexp.MustCompile("(?i)(b\\s*s|usd|\\$|[\\d.,]+%?|millones|mil)")

	var topHeadlines []string
	for i, h := range newsTitles {
		if i >= 8 {
			break
		}

		cleaned := numRegex.ReplaceAllString(h, "")
		cleaned = strings.TrimSpace(cleaned)
		cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

		if len(cleaned) > 140 {
			cleaned = cleaned[:137] + "..."
		}
		topHeadlines = append(topHeadlines, cleaned)
	}

	payload := map[string]any{
		"bcv_rates_official": compactRates,
		"binance_p2p_real": map[string]string{
			"buy":  fmt.Sprintf("%.2f", binanceRate.BuyPrice),
			"sell": fmt.Sprintf("%.2f", binanceRate.SellPrice),
		},
		"cleaned_headlines_context": topHeadlines,
	}

	b, _ := json.Marshal(payload)

	userContent := fmt.Sprintf(
		"INFORME FINANCIERO CAMBIARIO (ENTRADA DE DATOS):\n%s",
		string(b),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	systemPrompt := "Eres un algoritmo automatizado de análisis cuantitativo para fondos de cobertura. Tu única tarea es evaluar el INFORME FINANCIERO proporcionado por el usuario y generar una señal de trading en formato JSON estricto.\n\n" +
		"REGLAS MATEMÁTICAS ABSOLUTAS:\n" +
		"1. Está estrictamente prohibido inventar, alucinar o usar números que no estén explícitamente declarados en 'bcv_rates_official' o 'binance_p2p_real'. No uses tasas del pasado ni cifras externas.\n" +
		"2. Calcula la brecha real restando el USD del BCV del precio 'buy' de Binance P2P.\n" +
		"3. Redacta el 'rationale' en tercera persona neutral (ej: 'El sistema detecta...', 'Se recomienda...'). Debe ser un análisis técnico, frío y fluido de un solo párrafo sin ideas repetidas. Debe justificar la acción elegida basándose en la brecha calculada y el sentimiento macroeconómico de los titulares.\n" +
		"4. Coherencia de variables: Si el análisis determina que el spread es desfavorable no te dejes llevar por lo que te he dicho y recuerda que los valores que debes de tomar en cuenta de los rates son usd y eur por ves, la variable 'action' debe ser \"HOLD\" y el texto debe reflejar esa prudencia.\n\n" +
		"ESTRUCTURA DEL OBJETO JSON DE SALIDA:\n" +
		"{\n" +
		"  \"action\": \"BUY\",\n" +
		"  \"rationale\": \"Análisis cuantitativo de la brecha real y evaluación neutral del entorno de prensa y riesgo temporal.\",\n" +
		"  \"key_factors\": [\"Factor de mercado 1\", \"Factor de mercado 2\"],\n" +
		"  \"win_points\": 75.5,\n" +
		"  \"accuracy_rate\": 85\n" +
		"}\n\n" +
		"RESTRICCIONES DE FORMATO:\n" +
		"- 'action': Solo se permite \"BUY\", \"SELL\" o \"HOLD\".\n" +
		"- 'win_points': Asigna obligatoriamente un valor decimal flotante mayor a 0 que mida la fuerza de la señal.\n" +
		"- 'accuracy_rate': Debe ser un número entero calculado entre 1 and 100. Bajo ninguna circunstancia puede ser 0.\n" +
		"- No respondas con texto plano, ni introducciones, ni bloques de código markdown ```json. Entrega solo el objeto JSON puro."

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userContent),
	}

	resp, err := llm.GenerateContent(ctx, messages, llms.WithJSONMode())
	if err != nil {
		return models.TradeSignal{}, fmt.Errorf("langchaingo call error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return models.TradeSignal{}, fmt.Errorf("empty response from langchaingo")
	}

	content := strings.TrimSpace(resp.Choices[0].Content)
	if strings.HasPrefix(content, "```") {
		parts := strings.SplitN(content, "\n", 2)
		if len(parts) > 1 {
			content = strings.Trim(parts[1], "`\n ")
		}
	}

	var out forecastOutput
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return models.TradeSignal{}, fmt.Errorf("failed to unmarshal forecast JSON: %w; raw=%s", err, content)
	}

	ts := models.TradeSignal{
		Action:       out.Action,
		Rationale:    out.Rationale,
		KeyFactors:   models.JSONB(out.KeyFactors),
		WinPoints:    out.WinPoints,
		AccuracyRate: out.AccuracyRate,
		CreatedAt:    time.Now(),
	}

	return ts, nil
}
