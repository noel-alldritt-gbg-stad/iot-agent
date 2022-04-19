package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type mqttClient struct {
	cfg     Config
	log     zerolog.Logger
	options *mqtt.ClientOptions
}

func (c *mqttClient) Start() error {

	if !c.cfg.enabled {
		c.log.Warn().Msg("mqtt has been explicitly disabled with MQTT_DISABLED=true and will therefore not start")
		return nil
	}

	go c.run()

	return nil
}

var keepRunning bool = false // Temporary solution to be replaced with proper channels

func (c *mqttClient) run() {
	keepRunning = true

	client := mqtt.NewClient(c.options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		c.log.Fatal().Err(token.Error()).Msg("connection error")
	}

	for keepRunning {
		time.Sleep(1 * time.Second)
	}
}

func (c *mqttClient) Stop() {
	keepRunning = false
}
