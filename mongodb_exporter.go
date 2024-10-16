package mongodbexporter

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/model/pdata"
)

type mongodbExporter struct {
    client     *mongo.Client
    collection *mongo.Collection
}

func newLogsExporter(cfg *Config) (component.LogsExporter, error) {
    clientOptions := options.Client().ApplyURI(cfg.URI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }

    collection := client.Database(cfg.Database).Collection(cfg.Collection)

    exporter := &mongodbExporter{
        client:     client,
        collection: collection,
    }

    return exporterhelper.NewLogsExporter(
        cfg,
        exporter.pushLogs,
        exporterhelper.WithShutdown(exporter.Shutdown),
    )
}

func (e *mongodbExporter) pushLogs(ctx context.Context, ld pdata.Logs) error {
    // 遍歷日誌記錄並插入到 MongoDB
    for i := 0; i < ld.ResourceLogs().Len(); i++ {
        resourceLogs := ld.ResourceLogs().At(i)
        for j := 0; j < resourceLogs.ScopeLogs().Len(); j++ {
            scopeLogs := resourceLogs.ScopeLogs().At(j)
            for k := 0; k < scopeLogs.LogRecords().Len(); k++ {
                logRecord := scopeLogs.LogRecords().At(k)
                // 將 logRecord 轉換為適合 MongoDB 的格式
                doc := map[string]interface{}{
                    "timestamp":    logRecord.Timestamp().AsTime(),
                    "severity":     logRecord.SeverityText(),
                    "body":         logRecord.Body().AsString(),
                    "attributes":   logRecord.Attributes().AsRaw(),
                    "trace_id":     logRecord.TraceID().HexString(),
                    "span_id":      logRecord.SpanID().HexString(),
                    "flags":        logRecord.Flags(),
                    "observedTime": logRecord.ObservedTimestamp().AsTime(),
                }
                // 插入到 MongoDB
                _, err := e.collection.InsertOne(ctx, doc)
                if err != nil {
                    return fmt.Errorf("failed to insert log into MongoDB: %w", err)
                }
            }
        }
    }
    return nil
}

func (e *mongodbExporter) Shutdown(ctx context.Context) error {
    return e.client.Disconnect(ctx)
}

