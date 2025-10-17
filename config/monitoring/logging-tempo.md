# Logging & Tracing Conventions

## Loki
- Label `app="virtualization-gateway"` for REST layer logs.
- Required fields: `trace_id`, `vm_name`, `operation`, `result`, `reason`.
- Structured log example (JSON):
```
{"ts":"2024-01-01T00:00:00Z","level":"info","trace_id":"1234","vm_name":"demo","operation":"powerOn","result":"success"}
```

## Tempo
- Propagate `traceparent` header from frontend to gateway and controller.
- Use span names `gateway.<resource>.<verb>` and `controller.<resource>.<phase>`.

## Export
- Configure Promtail to scrape `/var/log/containers/*virtualization*.log` with `multiline` disabled.
- Tempo receive endpoint `tempo-distributor.monitoring.svc:4317` (OTLP gRPC).
