package main

import (
	"bytes"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var _sync sync.Once

func launchMainOnce(t *testing.T) {
	f := func() {
		go func() {
			main()
			t.Fatal("main failed")
		}()
	}
	_sync.Do(f)
}

func assertInLog(t *testing.T, logExpSubStr string, logReal string) {
	if !strings.Contains(logReal, logExpSubStr) {
		t.Errorf("'%s' not found in log '%s'", logExpSubStr, logReal)
	}
}

func assertNotInLog(t *testing.T, logExpSubStr string, logReal string) {
	if strings.Contains(logReal, logExpSubStr) {
		t.Errorf("'%s' found in log '%s'", logExpSubStr, logReal)
	}
}

func Test_Increment(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.Incr("example_metric.increment", []string{"environment:dev"}, 1)
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "example_metric.increment:1|c|#environment:dev", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_Decrement(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.Decr("example_metric.decrement", []string{"environment:dev"}, 1)
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "example_metric.decrement:-1|c|#environment:dev", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_Count(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.Count("example_metric.count", 2, []string{"environment:dev"}, 1)
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "example_metric.count:2|c|#environment:dev", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_Gauge(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.Gauge("example_metric.gauge", 2, []string{"environment:dev"}, 1)
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "example_metric.gauge:2|g|#environment:dev", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_Set(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.Set("example_metric.set", "10", []string{"environment:dev"}, 1)

	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "example_metric.set:10|s|#environment:dev", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_Histogram(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.Histogram("example_metric.histogram", float64(42), []string{"environment:dev"}, 1)
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "example_metric.histogram:42|h|#environment:dev", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_Event(t *testing.T) {
	t.Skip("Not implemented")
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.SimpleEvent("Something occurred", "Some message")
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "_e{18,12}:Something occurred|Some message", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_ServiceCheck(t *testing.T) {
	t.Skip("Not implemented")
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithoutTelemetry())
	if err != nil {
		t.Fatal(err)
	}
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	statsd.SimpleServiceCheck("application.service_check", 0)
	statsd.Flush()
	time.Sleep(1 * time.Second)
	realLog := logOut.String()
	assertInLog(t, "__sc|application.service_check|0", realLog)
	assertNotInLog(t, "Invalid event", realLog)
}

func Test_ErrorCase(t *testing.T) {
	var logOut bytes.Buffer

	log.SetOutput(&logOut)
	launchMainOnce(t)
	time.Sleep(1 * time.Second)

	conn, _ := net.Dial("udp", "127.0.0.1:8125")
	conn.Write([]byte("Some error"))
	time.Sleep(1 * time.Second)

	realLog := logOut.String()
	assertInLog(t, "Invalid event", realLog)
}
