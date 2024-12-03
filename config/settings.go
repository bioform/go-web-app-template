package config

import (
	"embed"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/bioform/go-web-app-template/pkg/env"
	"github.com/bioform/go-web-app-template/pkg/logging"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

// Settings main configuration struct for the service
type Settings struct {
	Event    EventConfig      `yaml:"event" json:"event" mapstructure:"event"`
	Kafka    *kafka.ConfigMap `yaml:"kafka" json:"kafka" mapstructure:"kafka"`
	Database Database         `yaml:"database" json:"database" mapstructure:"database"`
	Email    Email            `yaml:"email" json:"email" mapstructure:"email"`
}

var (
	//go:embed settings
	settingsDir embed.FS

	App *Settings
)

func init() {
	logging.InitLogger()

	var err error
	App, err = readConfig()
	if err != nil {
		slog.Error("Error reading config file", slog.Any("error", err))
		panic(err)
	}

	App.Kafka = App.fixKafkaConfig()
}

func readConfig() (*Settings, error) {
	viper.SetConfigType("yml")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Viper looks for values in environment variables first like this DATABASE.URL
	// instead of URL but since shell does not support dot format for environment variables,
	// we transform it to dash format using "strings.NewReplacer"
	// See https://99devops.com/overriding-configuration-variables-with-env-vars-in-go/
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := readFromSettingsFile("default.yml"); err != nil {
		return nil, err
	}

	envSpecificFile := strings.ToLower(env.App())
	_ = readFromSettingsFile(envSpecificFile + ".yml")

	_ = readFromSettingsFile("local.yml")

	conf := &Settings{}

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func readFromSettingsFile(path string) error {
	file, err := settingsDir.Open(filepath.Join("settings", path))
	if err != nil {
		return err
	}
	defer file.Close()

	if err := viper.ReadConfig(file); err != nil {
		return err
	}

	return nil
}
