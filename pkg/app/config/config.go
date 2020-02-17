package config

type Config struct {
	InfuraBaseURL   string `split_words:"true" required:"true"`
	InfuraProjectId string `split_words:"true" required:"true"`
}
