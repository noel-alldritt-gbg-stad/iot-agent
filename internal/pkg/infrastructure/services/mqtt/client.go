package mqtt

import (
	"crypto/tls"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type mqttClient struct {
	cfg    Config
	log    zerolog.Logger
	topics []string
}

func (c *mqttClient) Start() error {
	options := mqtt.NewClientOptions()

	connectionString := fmt.Sprintf("tls://%s:8883", c.cfg.host)
	options.AddBroker(connectionString)

	options.Username = c.cfg.user
	options.Password = c.cfg.password

	options.SetClientID("diwise/iot-agent/" + uuid.NewString())
	options.SetDefaultPublishHandler(MessageHandler)

	options.OnConnect = func(mc mqtt.Client) {
		c.log.Info().Msg("connected")
		for _, topic := range c.topics {
			c.log.Info().Msgf("subscribing to %s", topic)
			mc.Subscribe(topic, 0, nil)
		}
	}

	options.OnConnectionLost = func(mc mqtt.Client, err error) {
		c.log.Fatal().Err(err).Msg("connection lost")
	}

	options.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	go c.run(options)

	return nil
}

var keepRunning bool = false // Temporary solution to be replaced with proper channels

func (c *mqttClient) run(options *mqtt.ClientOptions) {
	keepRunning = true

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		c.log.Fatal().Err(token.Error()).Msg("connection error")
	}

	for keepRunning == true {
		time.Sleep(1 * time.Second)
	}
}

func (c *mqttClient) Stop() {
	keepRunning = false
}
