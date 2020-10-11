# go-redis extensions

## Tracing using OpenTelemetry

For more details see [documentation](https://redis.uptrace.dev/tracing/):

```go
import "github.com/go-redis/redisext"

rdb := redis.NewClient(&redis.Options{...})
rdb.AddHook(redisext.OpenTelemetryHook{})
```

## Tracing using OpenCensus

Installation:

```bash
go get github.com/go-redis/redisext/rediscensus
```

Setup outline:

```
import "github.com/go-redis/redisext/rediscensus"

rdb := redis.NewClient(&redis.Options{...})
rdb.AddHook(rediscensus.OpenCensusHook{})
```
