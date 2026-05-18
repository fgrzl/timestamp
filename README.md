[![ci](https://github.com/fgrzl/timestamp/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/timestamp/actions/workflows/ci.yml)
[![Dependabot Updates](https://github.com/fgrzl/timestamp/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/timestamp/actions/workflows/dependabot/dependabot-updates)

# timestamp
A Go library that provides **process-local monotonic timestamps** with accurate initialization. The library fetches the initial time from NTP servers (or system clock) at startup, then uses Go's built-in monotonic clock to ensure consistent, always-increasing timestamps throughout the application's lifecycle.

## 🎯 Design Philosophy
This library solves a specific problem: **ensuring monotonic timestamp ordering within a process** while starting with accurate wall-clock time.

### Key Guarantees:
- ✅ **Monotonic**: Timestamps always increase, never go backwards
- ⏰ **Accurate Start**: Initial time synchronized with NTP servers
- 🔒 **No Drift**: Uses monotonic clock, immune to system clock adjustments
- 🚀 **High Performance**: No network calls after initialization
- 🎯 **Process-Scoped**: Consistent ordering within application lifecycle

### Trade-offs:
- ❌ **Not wall-clock accurate** over long periods (no re-synchronization)
- ❌ **Process-local only** (not synchronized across different processes)
- ❌ **Single initialization** (accuracy locked at startup)

## 🚀 Features
- ✅ **Process-local monotonic timestamps** based on Go's built-in monotonic clock
- 🌐 **NTP initialization** from configurable servers for accurate startup time
- 🔒 **Automatic fallback** to system time if NTP servers are unreachable
- 🔧 **Environment-based configuration** for NTP server selection
- 🏎️ **Zero overhead** after initialization - no background tasks or network calls
- ⚡ **Thread-safe** timestamp generation

## 🎯 When to Use This Library

### ✅ Perfect For:
- **Event ordering** within a microservice or application
- **Rate limiting** with consistent time progression
- **Metrics collection** where monotonic progression matters
- **Request tracing** with guaranteed timestamp ordering
- **Database operations** requiring consistent time-based sorting
- **Avoiding issues** with system clock adjustments (NTP corrections, leap seconds)

### ❌ Not Suitable For:
- **Cross-process synchronization** (each process has its own timeline)
- **Long-running services** requiring wall-clock accuracy over days/weeks
- **Time-sensitive calculations** that need real-time precision
- **Distributed systems** where processes need synchronized timestamps
- **Financial systems** requiring exact wall-clock time

### 🤔 Alternative Solutions:
- **Need cross-process sync?** → Use external time service (Redis, database)
- **Need wall-clock accuracy?** → Use `time.Now()` with periodic NTP sync
- **Need distributed sync?** → Consider logical clocks (Lamport, Vector clocks)

## ⚙️ Installation
```bash
go get github.com/fgrzl/timestamp
```

## Documentation

Guides: **[docs/](docs/README.md)** — [overview](docs/overview.md), [getting started](docs/getting-started.md)

## 🛠️ Usage

### Basic Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/fgrzl/timestamp"
)

func main() {
	// Get timestamps - guaranteed to be monotonic
	ts1 := timestamp.GetTimestamp()
	fmt.Println("Initial Timestamp:", ts1)

	time.Sleep(100 * time.Millisecond)

	ts2 := timestamp.GetTimestamp()
	fmt.Println("Later Timestamp:", ts2)

	// This will ALWAYS be true (monotonic guarantee)
	if ts2 > ts1 {
		fmt.Println("✅ Monotonic behavior confirmed")
		fmt.Printf("Elapsed: %d milliseconds\n", ts2-ts1)
	}
}
```

### Practical Examples

#### Event Logging with Guaranteed Ordering
```go
type Event struct {
    ID        string `json:"id"`
    Timestamp int64  `json:"timestamp"`
    Message   string `json:"message"`
}

func logEvent(message string) Event {
    return Event{
        ID:        generateID(),
        Timestamp: timestamp.GetTimestamp(), // Always increasing
        Message:   message,
    }
}
```

#### Rate Limiting
```go
type RateLimiter struct {
    lastRequest int64
    minInterval int64 // milliseconds
}

func (rl *RateLimiter) Allow() bool {
    now := timestamp.GetTimestamp()
    if now-rl.lastRequest >= rl.minInterval {
        rl.lastRequest = now
        return true
    }
    return false
}
```

#### Performance Measurement
```go
func measureOperation() {
    start := timestamp.GetTimestamp()
    
    // Your operation here
    performWork()
    
    end := timestamp.GetTimestamp()
    duration := end - start
    fmt.Printf("Operation took %d milliseconds\n", duration)
}
```

## 🌐 Configuration

### NTP Server Configuration
Control the initial time source using the `FGRZL_TIME_SERVER` environment variable:

| Value | Behavior | Use Case |
|-------|----------|----------|
| `system` | Use system clock directly | Testing, offline environments |
| `default` or empty | Try multiple NTP servers | Production (recommended) |
| Custom server | Use specific NTP server | Corporate networks, specific requirements |

### Examples
```bash
# Use system time (fastest startup, may be inaccurate)
export FGRZL_TIME_SERVER="system"

# Use default NTP servers (recommended for production)
export FGRZL_TIME_SERVER="default"
# OR simply don't set the variable

# Use specific NTP server
export FGRZL_TIME_SERVER="time.cloudflare.com:123"
export FGRZL_TIME_SERVER="pool.ntp.org:123"
export FGRZL_TIME_SERVER="time.nist.gov:123"
```

### Default NTP Servers
When using `default` mode, the library tries these servers in order:
- `time.google.com:123`
- `time.aws.com:123`
- `time.cloudflare.com:123`
- `time.windows.com:123`

If all NTP servers fail, the library automatically falls back to system time.

## 🔧 Implementation Details

### How It Works
1. **Initialization** (once per process):
   - Fetches accurate time from NTP server (or uses system time)
   - Records the initial timestamp and Go's monotonic clock reference
   - This happens automatically when the package is first imported

2. **Timestamp Generation** (every call to `GetTimestamp()`):
   - Calculates elapsed time using Go's monotonic clock
   - Adds elapsed time to initial timestamp
   - Returns result in Unix milliseconds

### Performance Characteristics
- **Initialization**: ~100ms (network-dependent)
- **Per-call overhead**: ~10ns (pure calculation)
- **Memory usage**: <100 bytes (single clock instance)
- **Thread safety**: Concurrent-safe reads

### Precision & Accuracy
- **Resolution**: 1 millisecond
- **Initial accuracy**: ±100ms (NTP-dependent)
- **Long-term drift**: None (monotonic clock based)
- **Thread consistency**: Perfect (same monotonic reference)
