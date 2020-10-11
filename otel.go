package redisext

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redisext/cmdutil"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

var tracer = global.Tracer("github.com/go-redis/redis")

type OpenTelemetryHook struct{}

var _ redis.Hook = OpenTelemetryHook{}

func (OpenTelemetryHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	ctx, span := tracer.Start(ctx, cmd.FullName())
	span.SetAttributes(
		label.String("db.system", "redis"),
		label.String("redis.cmd", cmdutil.CmdString(cmd)),
	)

	return ctx, nil
}

func (OpenTelemetryHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if err := cmd.Err(); err != nil {
		recordError(ctx, span, err)
	}
	span.End()
	return nil
}

func (OpenTelemetryHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	summary, cmdsString := cmdutil.CmdsString(cmds)

	ctx, span := tracer.Start(ctx, "pipeline "+summary)
	span.SetAttributes(
		label.String("db.system", "redis"),
		label.Int("redis.num_cmd", len(cmds)),
		label.String("redis.cmds", cmdsString),
	)

	return ctx, nil
}

func (OpenTelemetryHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if err := cmds[0].Err(); err != nil {
		recordError(ctx, span, err)
	}
	span.End()
	return nil
}

func recordError(ctx context.Context, span trace.Span, err error) {
	if err != redis.Nil {
		span.RecordError(ctx, err)
	}
}
