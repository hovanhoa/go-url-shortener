package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// RedisWithTracing wraps Redis operations with OpenTelemetry tracing.
func RedisWithTracing(ctx context.Context, operation, key string, fn func() error) error {
	tr := Tracer()
	ctx, span := tr.Start(ctx, "redis."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", operation),
			attribute.String("db.redis.key", key),
		),
	)
	defer span.End()

	err := fn()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}
