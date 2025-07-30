package config

type SecretConfig interface {
	TokenSecret() string
}

type secretConfigImpl struct {
	tokenSecret string
}

func (s *secretConfigImpl) TokenSecret() string {
	return s.tokenSecret
}

func NewSecretConfig(tokenSecret string) SecretConfig {
	return &secretConfigImpl{tokenSecret: tokenSecret}
}
