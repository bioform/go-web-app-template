package config

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// EventConfig configuration for Kafka consumer
type EventConfig struct {
	Topic string `yaml:"topic" json:"topic" mapstructure:"topic"`
}

// fixKafkaConfig retrieves a Kafka.ConfigMap compatible struct from
// our configuration. Viper supports nested configuration. However, we
// need a flatten struct for Kafka
func (c Settings) fixKafkaConfig() *kafka.ConfigMap {
	if c.Kafka == nil {
		return nil
	}

	cm := &kafka.ConfigMap{}

	rawConfig := make(map[string]any)
	for k, v := range *c.Kafka {
		rawConfig[k] = v
	}

	flattenKafkaConfigMap("", rawConfig, cm)

	return cm
}

// flattenKafkaConfigMap converts a nested struct into a flatten config map
// TODO this is specifcally for Kafka.ConfigMap. Maybe open this up to other
// structs
func flattenKafkaConfigMap(prefix string, src map[string]any, cm *kafka.ConfigMap) {
	if prefix != "" {
		prefix += "."
	}

	for k, v := range src {
		switch child := v.(type) {
		case map[string]any:
			flattenKafkaConfigMap(prefix+k, child, cm)
		default:
			key := prefix + k
			err := cm.SetKey(key, child)
			if err != nil {
				log.Fatalf("error setting key %s: %v", key, err)
			}
		}
	}
}
