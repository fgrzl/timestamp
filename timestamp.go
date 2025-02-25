package timestamp

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// TimeServer environment variable for NTP server.
const TimeServer = "FGRZL_TIME_SERVER"

var (
	globalClock *clock
	once        sync.Once
)

// Clock holds the start time for the monotonic clock.
type clock struct {
	startTime int64 // Unix timestamp in milliseconds
	start     time.Time
}

// Initialize the global clock once during application startup.
func init() {
	once.Do(func() {
		t := getCurrentTime()
		globalClock = &clock{
			startTime: t.UnixMilli(),
			start:     t,
		}
	})
}

// GetTimestamp returns a timestamp using monotonic elapsed time.
func GetTimestamp() int64 {
	elapsed := time.Since(globalClock.start)
	return globalClock.startTime + elapsed.Milliseconds()
}

// GetTimeServer fetches the configured NTP server from the environment.
func GetTimeServer() string {
	return os.Getenv(TimeServer)
}

// Default list of NTP servers.
var ntpServers = []string{
	"time.google.com:123",
	"time.aws.com:123",
	"time.cloudflare.com:123",
	"time.windows.com:123",
}

// getCurrentTime attempts to fetch time from NTP or falls back to system time.
func getCurrentTime() time.Time {
	server := GetTimeServer()

	if server == "system" {
		return time.Now()
	}

	if server == "default" || server == "" {
		for _, s := range ntpServers {
			if t, err := ntpTime(s); err == nil {
				return t
			} else {
				log.Printf("NTP server %s failed: %v", s, err)
			}
		}
		log.Println("All NTP servers failed. Falling back to system time.")
		return time.Now()
	}

	ntpTime, err := ntpTime(server)
	if err != nil {
		log.Printf("Failed to fetch time from NTP server (%s): %v. Falling back to system time.", server, err)
		return time.Now()
	}

	return ntpTime
}

// ntpTime fetches the current time from the specified NTP server.
func ntpTime(server string) (time.Time, error) {
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to resolve address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Set a timeout for the connection
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	req := make([]byte, 48)
	req[0] = 0x1B // NTP version 3, client mode

	if _, err = conn.Write(req); err != nil {
		return time.Time{}, fmt.Errorf("failed to send request: %w", err)
	}

	resp := make([]byte, 48)
	if _, err = conn.Read(resp); err != nil {
		return time.Time{}, fmt.Errorf("failed to read response: %w", err)
	}

	seconds := binary.BigEndian.Uint32(resp[40:44])
	fraction := binary.BigEndian.Uint32(resp[44:48])

	ntpSeconds := float64(seconds) + float64(fraction)/0x100000000
	unixSeconds := ntpSeconds - 2208988800 // NTP epoch offset
	return time.Unix(int64(unixSeconds), 0), nil
}
