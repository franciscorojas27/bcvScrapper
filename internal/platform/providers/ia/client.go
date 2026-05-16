package ia

import (
	"os"

	"github.com/tmc/langchaingo/llms/openai"
)

func ClientIA() (*openai.LLM, error) {
	llm, err := openai.New(
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(os.Getenv("GROQ_API_KEY")),
		openai.WithModel("allam-2-7b"),
	)
	if err != nil {
		return nil, err
	}
	return llm, nil
}
