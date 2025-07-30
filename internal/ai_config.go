package internal

type AiClientConfig struct {
	ApiKey string
	Model  string
	Url    string
}

func NewAiClientConfig(apiKey string, model string, url string) *AiClientConfig {
	return &AiClientConfig{
		ApiKey: apiKey,
		Model:  model,
		Url:    url,
	}
}
