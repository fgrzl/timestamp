[![ci](https://github.com/fgrzl/timestamp/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/timestamp/actions/workflows/ci.yml)
[![Dependabot Updates](https://github.com/fgrzl/timestamp/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/timestamp/actions/workflows/dependabot/dependabot-updates)

# timestamp
A Go library that provides a reliable, monotonic timestamp, initialized with the current time from an NTP server or the system clock. Once initialized at application startup, the clock uses Go's built-in monotonic timekeeping to ensure consistent and increasing timestamps throughout the application's lifecycle.

## 🚀 Features
- ✅ Monotonic timestamp generation based on Go's built-in monotonic clock.
- 🌐 Fetch initial time from configurable NTP servers.
- 🔒 Fallback to system time if NTP servers are unreachable.
- 🔧 Supports custom NTP server configuration via environment variable.
- 🏎️ Lightweight with no unnecessary background tasks.

## ⚙️ Installation
```bash
go get github.com/fgrzl/timestamp
```

## 🛠️ Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/fgrzl/timestamp"
)

func main() {
	// Get the initial timestamp (in Unix milliseconds)
	ts1 := timestamp.GetTimestamp()
	fmt.Println("Initial Timestamp:", ts1)

	// Wait for 1 second
	time.Sleep(1 * time.Second)

	// Get a new timestamp
	ts2 := timestamp.GetTimestamp()
	fmt.Println("New Timestamp:", ts2)

	if ts2 > ts1 {
		fmt.Println("Monotonic clock is working correctly.")
	}
}
```

## 🌐 Configuring NTP Server
You can configure the NTP server using the FGRZL_TIME_SERVER environment variable:

Options:
- system - Use the system clock directly.
- default - Use Google's NTP server (time.google.com:123) by default.
- Custom NTP server (e.g., "pool.ntp.org:123")

```bash
export FGRZL_TIME_SERVER="time.cloudflare.com:123"
```
