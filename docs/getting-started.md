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
// ts2 >= ts1 always
```

Timestamps are millisecond-resolution `int64` values suitable for sorting and storage keys.

## Initialization

The library initializes automatically on first use. For explicit control or custom NTP servers, call the package initialization API before generating timestamps (see package godoc for `Initialize` and environment variables such as `TIMESTAMP_NTP_SERVERS`).

If NTP is unreachable, initialization falls back to the system clock and logs the outcome.

## Environment

Configure NTP endpoints when defaults are unsuitable in your network:

```bash
export TIMESTAMP_NTP_SERVERS="pool.ntp.org,time.google.com"
```

## Thread safety

`GetTimestamp()` is safe for concurrent use from any goroutine after initialization completes.

## Tests

```bash
go test ./...
```

Use test helpers or dependency injection patterns in your app if you need deterministic timestamps in unit tests.
