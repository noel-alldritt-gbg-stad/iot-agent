package mqtt

import (
	"os"

	"github.com/rs/zerolog"
)

type Client interface {
	Start() error
	Stop()
}

type Config struct {
	host     string
	user     string
	password string
}

func NewClient(logger zerolog.Logger, cfg Config) (Client, error) {
	return &mqttClient{
		cfg:    cfg,
		log:    logger.With().Str("mqtt-host", cfg.host).Logger(),
		topics: []string{"application/53/device/#"},
	}, nil
}

func NewConfigFromEnvironment() (Config, error) {
	cfg := Config{
		host:     os.Getenv("MQTT_HOST"),
		user:     os.Getenv("MQTT_USER"),
		password: os.Getenv("MQTT_PASSWORD"),
	}

	return cfg, nil
}
