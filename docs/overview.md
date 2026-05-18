# Overview

The timestamp library solves **monotonic ordering within one process** while starting from an accurate wall-clock baseline.

## Guarantees

| Guarantee | Detail |
|-----------|--------|
| Monotonic | `GetTimestamp()` never decreases within the process |
| Accurate start | One-time NTP sync (or system clock fallback) at init |
| No drift correction | Uses Go's monotonic clock after init — immune to wall-clock jumps |
| Zero steady-state cost | No background goroutines or network calls after initialization |

## Trade-offs

- **Not wall-clock accurate** over long uptimes — no periodic re-sync
- **Not cross-process** — each process has its own timeline
- **Single initialization** — accuracy is fixed at startup

## When to use

**Good fits:** event ordering inside a service, rate limiting, metrics ordering, request IDs where monotonicity matters more than global wall time.

**Poor fits:** distributed coordination, financial timestamps, cross-service causality without logical clocks.

## Alternatives

| Need | Consider |
|------|----------|
| Cross-process sync | Redis TIME, database clock, time service |
| Wall-clock accuracy | `time.Now()` with periodic NTP |
| Distributed ordering | Lamport / vector clocks |

## Configuration

NTP servers are configurable via environment variables (see [Getting started](getting-started.md)). Initialization runs once per process, typically from `init()` or explicit `Initialize()`.
