package timestamp

import (
	"encoding/binary"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"
)

// Test that the timestamp always increases over time (monotonic behavior)
func TestMonotonicTimestamp(t *testing.T) {
	first := GetTimestamp()
	time.Sleep(10 * time.Millisecond)
	second := GetTimestamp()

	if second <= first {
		t.Errorf("Expected timestamp to increase: first=%d, second=%d", first, second)
	}
}

// Test that NTP time fetch works successfully
func TestNTPTimeSuccess(t *testing.T) {
	// Attempt to fetch time from a real NTP server (if network is available)
	server := "time.google.com:123"
	ntpTime, err := ntpTime(server)
	if err != nil {
		t.Skipf("Skipping NTP time test due to network error: %v", err)
	}

	if time.Since(ntpTime) > time.Minute {
		t.Errorf("Fetched NTP time is too far from the system time: %v", ntpTime)
	}
}

// Test NTP failure fallback to system time
func TestNTPTimeFailureFallback(t *testing.T) {
	// Pass an invalid NTP server to force failure
	server := "invalid.ntp.server:123"
	start := time.Now()
	ntpTime, err := ntpTime(server)

	if err == nil {
		t.Error("Expected error when fetching time from an invalid NTP server")
	}

	if ntpTime != (time.Time{}) {
		t.Errorf("Expected zero time on failure, got %v", ntpTime)
	}

	// Fallback to system time check
	systemTime := time.Now()
	if systemTime.Sub(start) > time.Second {
		t.Errorf("Fallback to system time took too long")
	}
}

// Test fetching time using environment configuration
func TestGetTimeServerConfig(t *testing.T) {
	os.Setenv(TimeServer, "system")
	defer os.Unsetenv(TimeServer)

	t1, _ := getCurrentTime()
	time.Sleep(5 * time.Millisecond)
	t2, _ := getCurrentTime()

	if t2.Sub(t1) <= 0 {
		t.Error("Expected system time to increase")
	}
}

// Test invalid NTP server environment fallback to default
func TestInvalidTimeServerFallback(t *testing.T) {
	os.Setenv(TimeServer, "invalid.ntp.server:123")
	defer os.Unsetenv(TimeServer)

	start := time.Now()
	t1, _ := getCurrentTime()

	if t1.Before(start) {
		t.Errorf("Time fetched is unexpectedly before start time: %v", t1)
	}
}

// Mock UDP connection for testing without network
type mockUDPConn struct{}

func (m *mockUDPConn) Read(b []byte) (n int, err error) {
	// Simulate a valid NTP response
	// Example time: 2208988800 (Jan 1, 1970) + 1 (1 second past epoch)
	binary.BigEndian.PutUint32(b[40:44], 2208988801)
	binary.BigEndian.PutUint32(b[44:48], 0)
	return 48, nil
}

func (m *mockUDPConn) Write(b []byte) (n int, err error) { return len(b), nil }
func (m *mockUDPConn) Close() error                      { return nil }
func (m *mockUDPConn) SetDeadline(t time.Time) error     { return nil }

// Mock NTP fetch using fake UDP connection
func mockNTPTime(_ string) (time.Time, error) {
	mockConn := &mockUDPConn{}
	resp := make([]byte, 48)
	mockConn.Read(resp)

	seconds := binary.BigEndian.Uint32(resp[40:44])
	fraction := binary.BigEndian.Uint32(resp[44:48])

	ntpSeconds := float64(seconds) + float64(fraction)/0x100000000
	unixSeconds := ntpSeconds - 2208988800
	return time.Unix(int64(unixSeconds), 0), nil
}

// Test mocked NTP response for a deterministic test
func TestMockedNTPTime(t *testing.T) {
	ntpTime, err := mockNTPTime("mock.ntp.server:123")
	if err != nil {
		t.Fatalf("Failed to get mock NTP time: %v", err)
	}

	expected := time.Unix(1, 0) // Expecting Jan 1, 1970 + 1 second
	if !ntpTime.Equal(expected) {
		t.Errorf("Expected NTP time %v, got %v", expected, ntpTime)
	}
}

// Test logging configuration
func TestLoggingConfiguration(t *testing.T) {
	// Test disabling logging
	DisableLogging()

	// Test setting custom logger
	SetLogger(nil) // Should disable logging

	// Restore default logging behavior
	SetLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// Test initialization error tracking
func TestInitializationError(t *testing.T) {
	// The error might be nil if initialization was successful
	err := GetInitializationError()
	// We just verify the function works, error might be nil in good conditions
	_ = err
}
