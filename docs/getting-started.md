# Getting started

## Install

```bash
go get github.com/fgrzl/timestamp
```

## Basic usage

```go
import "github.com/fgrzl/timestamp"

ts1 := timestamp.GetTimestamp()
// ... work ...
ts2 := timestamp.GetTimestamp()
// ts2 >= ts1 always within this process
```

Timestamps are millisecond-resolution `int64` values.

## Initialization

The clock initializes automatically when the package is **imported** (`init` + `sync.Once`). There is no public `Initialize()` function.

Check `timestamp.GetInitializationError()` if you need to detect NTP fallback failures.

## Configuration (`FGRZL_TIME_SERVER`)

| Value | Behavior |
|-------|----------|
| unset or `default` | Try built-in NTP server list, then system clock |
| `system` | Use system clock only (no NTP) |
| any other string | Single NTP host (e.g. `pool.ntp.org:123`) |

```bash
export FGRZL_TIME_SERVER=default
# or
export FGRZL_TIME_SERVER=time.google.com:123
```

There is no comma-separated multi-server env var; multiple hosts are only used in the built-in default list.

## Logging

```go
timestamp.SetLogger(myLogger)
// or
timestamp.DisableLogging()
```

## Thread safety

`GetTimestamp()` is safe for concurrent use after package init completes.

## Tests

```bash
go test ./...
```

See the root [README](../README.md) for design trade-offs and when not to use this library.
