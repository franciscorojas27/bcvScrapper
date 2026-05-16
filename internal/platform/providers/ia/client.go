package ia

import (
	"os"

	"github.com/tmc/langchaingo/llms/openai"
)

func ClientIA() (*openai.LLM, error) {
	if os.Getenv("GROQ_API_KEY") == "" {
		return nil, nil
	}
	if os.Getenv("GROQ_MODEL") == "" {
		return nil, nil
	}
	llm, err := openai.New(
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(os.Getenv("GROQ_API_KEY")),
		openai.WithModel(os.Getenv("GROQ_MODEL")),
	)
	if err != nil {
		return nil, err
	}
	return llm, nil
}
