package mongodbexporter

import (
    "go.opentelemetry.io/collector/config"
)

// Config 定義導出器的配置結構
type Config struct {
    config.ExporterSettings `mapstructure:",squash"`
    URI                     string `mapstructure:"uri"`
    Database                string `mapstructure:"database"`
    Collection              string `mapstructure:"collection"`
}


