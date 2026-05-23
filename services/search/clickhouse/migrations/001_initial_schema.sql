CREATE TABLE IF NOT EXISTS security_metrics (
    timestamp DateTime,
    event_type LowCardinality(String),
    source_ip IPv4,
    dest_ip IPv4,
    source_port UInt16,
    dest_port UInt16,
    user_name String,
    severity UInt8,
    action LowCardinality(String),
    outcome LowCardinality(String),
    tenant_id String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (tenant_id, timestamp, event_type);

CREATE TABLE IF NOT EXISTS kafka_ingest_stream (
    payload String
) ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
         kafka_topic_list = 'normalized-events',
         kafka_group_name = 'clickhouse-ingest',
         kafka_format = 'JSONAsString';

-- Materialized view to transform JSON from Kafka into structured table
-- Note: Simplified mapping for demonstration
CREATE MATERIALIZED VIEW IF NOT EXISTS security_metrics_mv TO security_metrics AS
SELECT
    parseDateTimeBestEffort(JSONExtractString(payload, '@timestamp')) AS timestamp,
    JSONExtractString(payload, 'event', 'kind') AS event_type,
    IPv4StringToNumOrDefault(JSONExtractString(payload, 'source', 'ip'), '0.0.0.0') AS source_ip,
    IPv4StringToNumOrDefault(JSONExtractString(payload, 'destination', 'ip'), '0.0.0.0') AS dest_ip,
    JSONExtractUInt(payload, 'source', 'port') AS source_port,
    JSONExtractUInt(payload, 'destination', 'port') AS dest_port,
    JSONExtractString(payload, 'user', 'name') AS user_name,
    JSONExtractUInt(payload, 'event', 'severity') AS severity,
    JSONExtractString(payload, 'event', 'action') AS action,
    JSONExtractString(payload, 'event', 'outcome') AS outcome,
    JSONExtractString(payload, 'metadata', 'tenant_id') AS tenant_id
FROM kafka_ingest_stream;
