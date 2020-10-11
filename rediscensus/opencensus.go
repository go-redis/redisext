package rediscensus

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redisext/cmdutil"
	"go.opencensus.io/trace"
)

type OpenCensusHook struct{}

var _ redis.Hook = OpenCensusHook{}

func (OpenCensusHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx, span := trace.StartSpan(ctx, cmd.FullName())
	span.AddAttributes(trace.StringAttribute("db.system", "redis"),
		trace.StringAttribute("redis.cmd", cmdutil.CmdString(cmd)))

	return ctx, nil
}

func (OpenCensusHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.FromContext(ctx)
	if err := cmd.Err(); err != nil {
		recordErrorOnOCSpan(ctx, span, err)
	}
	span.End()
	return nil
}

func (OpenCensusHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (OpenCensusHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func recordErrorOnOCSpan(ctx context.Context, span *trace.Span, err error) {
	if err != redis.Nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate([]trace.Attribute{trace.StringAttribute("Error", "redis error")}, err.Error())
	}
}
