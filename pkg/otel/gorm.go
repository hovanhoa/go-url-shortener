package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GormWithTracing wraps GORM operations with OpenTelemetry tracing.
func GormWithTracing(ctx context.Context, db *gorm.DB, operation string, fn func(*gorm.DB) error) error {
	tr := Tracer()
	ctx, span := tr.Start(ctx, "gorm."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", operation),
		),
	)
	defer span.End()

	err := fn(db.WithContext(ctx))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}
