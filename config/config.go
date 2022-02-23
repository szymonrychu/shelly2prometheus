package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ShellyUrl       string `mapstructure:"SHELLY_URL"`
	MetricsPort     int    `mapstructure:"METRICS_PORT"`
	MetricsPrefix   string `mapstructure:"METRICS_PREFIX"`
	MetricsEndpoint string `mapstructure:"METRICS_ENDPOINT"`
}

func (config *Config) Load() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("shelly2prometheus")
	v.AutomaticEnv()
	v.SetDefault("shelly_url", "http://127.0.0.1")
	v.SetDefault("metrics_prefix", "shelly25")
	v.SetDefault("metrics_port", 8080)
	v.SetDefault("metrics_endpoint", "/metrics")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.Unmarshal(config); err != nil {
		return err
	}
	log.Printf("Loaded config {shelly_url:%s, metrics_port:%d, metrics_prefix:%s, metrics_endpoint:%s}", config.ShellyUrl, config.MetricsPort, config.MetricsPrefix, config.MetricsEndpoint)
	return nil
}
