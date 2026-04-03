package observability

import "go.uber.org/zap"

// For Phase 1 we expose a placeholder to attach tracing/metrics later.
// This keeps call sites stable when OpenTelemetry is added.
func Startup(log *zap.Logger) {
	log.Info("observability initialized", zap.String("provider", "stdout"))
}
