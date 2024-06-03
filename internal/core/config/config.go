package config

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go-template/internal/log"
)

type Configuration struct {
	ServerUrl         string   `mapstructure:"ServerUrl" tip:"Server Url"`
	ServerPort        int      `mapstructure:"ServerPort" tip:"Server port"`
	TranslationAPIKey string   `mapstructure:"TranslationAPIKey" tip:"Translation API Key"`
	Log               Log      `mapstructure:"Log"`
	Database          Database `mapstructure:"Database"`
}
type Log struct {
	Level int `mapstructure:"Level" tip:"Minimum level to log: (-4:Debug, 0:Info, 4:Warning, 8:Error)"`
	Mode  int `mapstructure:"Mode" tip:"Log format (1: JSON, 2:Structured text)"`
}

type Database struct {
	URL string `mapstructure:"Url" tip:"The Datasource name locator"`
}

func Load() (*Configuration, error) {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Warn(ctx, "Error loading .env file, using default environment variables", "err", err)
	}

	bindEnv()

	config := &Configuration{}

	if err := viper.Unmarshal(config); err != nil {
		log.Error(ctx, "error unmarshalling configuration", "err", err)
	}

	err = checkEnvVars(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func bindEnv() {
	_ = viper.BindEnv("ServerUrl", "SERVER_URL")
	_ = viper.BindEnv("ServerPort", "SERVER_PORT")
	_ = viper.BindEnv("TranslationAPIKey", "TRANSLATION_SERVICE_API_KEY")

	_ = viper.BindEnv("Log.Level", "LOG_LEVEL")
	_ = viper.BindEnv("Log.Mode", "LOG_MODE")

	_ = viper.BindEnv("Database.URL", "TRANSLATION_DATABASE_URL")

	viper.AutomaticEnv()
}

func checkEnvVars(cfg *Configuration) error {
	if cfg.ServerUrl == "" {
		return errors.New("SERVER_URL env var is required")
	}
	if cfg.ServerPort == 0 {
		return errors.New("SERVER_PORT env var is required")
	}

	logLevels := []int{-4, 0, 4, 8}
	if isInSlice(cfg.Log.Level, logLevels) == false {
		return errors.New("LOG_LEVEL env var is required. Possible values: [-4, 0, 4, 8]")
	}
	logModes := []int{1, 2}
	if isInSlice(cfg.Log.Mode, logModes) == false {
		return errors.New("LOG_MODE env var is required. Possible values: [1, 2]")
	}
	if cfg.TranslationAPIKey == "" {
		return errors.New("TRANSLATION_SERVICE_API_KEY env var is required")
	}
	return nil
}

func isInSlice(value int, slice []int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
