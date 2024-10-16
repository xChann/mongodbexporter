package mongodbexporter

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/exporter/exporterhelper"
    "go.opentelemetry.io/collector/config"
)

func NewFactory() component.ExporterFactory {
    return component.NewExporterFactory(
        "mongodb",
        createDefaultConfig,
        component.WithLogsExporter(createLogsExporter),
    )
}

func createDefaultConfig() component.ExporterConfig {
    return &Config{
        ExporterSettings: config.NewExporterSettings(component.NewID("mongodb")),
        URI:              "mongodb://localhost:27017",
        Database:         "otel_logs",
        Collection:       "logs",
    }
}

func createLogsExporter(
    ctx context.Context,
    set component.ExporterCreateSettings,
    cfg component.ExporterConfig,
) (component.LogsExporter, error) {
    return newLogsExporter(cfg.(*Config))
}

