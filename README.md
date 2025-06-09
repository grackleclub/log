# log
simple structured log wrapper for go, with tint by [lmittmann/tint](https://github.com/lmittmann/tint) 🌈

![badge](https://github.com/ddbgio/log/actions/workflows/test.yml/badge.svg)

## examples
```go
import logger "github.com/grackleclub/log"

var log *slog.Logger // a global logger can be set

// If empty handler options are provided, the DEBUG
// env var will determine log level and source info.
_ := os.Setenv("DEBUG", "1")

// No color env var is always respected,
// but no effort is made to autodetect.
_ := os.Unsetenv("NO_COLOR")

log, _ = logger.New(slog.HandlerOptions{})

// Because this is regular slog, you can use log.With
// or attrs or any other slog features.
log = log.With("service", "cool_service")


log.Info("logging initialized", "msg", "hello world")
```
