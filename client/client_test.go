package client_test

import (
	"../client"
	"../server"
	"testing"
	"time"
)

func NewTestClient() *client.Instance {
	opts := client.NewOpts()
	c := client.New(opts)
	return c
}

func NewTestServer(init bool, listen bool) *server.Instance {
	s := server.New(server.NewOpts())
	if init {
		if err := s.Init(); err != nil {
			panic(err)
		}
	}
	if listen {
		if err := s.StartListening(); err != nil {
			panic(err)
		}
	}
	return s
}

func TestNew(t *testing.T) {
	c := NewTestClient()
	if c == nil {
		t.Error()
		return
	}

	// new series
	series := c.Series("mySeries")

	// timestamp
	now := c.Now()
	const oneMinute = 60 * 1000
	const writeValue = 10.1

	// start server
	s := NewTestServer(true, true)

	// write
	{
		result := series.Write(now, writeValue)
		if result.Error != nil {
			t.Error(result.Error)
		}
	}

	// read
	{
		result := series.QueryBuilder().From(now - oneMinute).To(now + oneMinute).Execute()
		if result.Error != nil {
			t.Error(result.Error)
		}
		if result.Results == nil {
			t.Error()
		}
		if len(result.Results) != 1 {
			t.Error(result.Results)
			return
		}
		var ts uint64
		var value float64
		for ts, value = range result.Results {
			// no need to do something
		}
		if ts != now {
			t.Error(ts, now)
		}
		if value != writeValue {
			t.Error(value)
		}
		//t.Log(ts, value)
	}

	c.Close()
	_ = s.Shutdown()
}

func TestWritePerformance(t *testing.T) {
	// start server
	s := NewTestServer(true, true)

	c := NewTestClient()
	now := c.Now()
	startTime := time.Now()
	const minTime = 1 * time.Second
	const minIters = 100
	const writeValue = 10.1
	series := c.Series("benchmarkSeriesWrite")
	var i int
	for i = 0; i < 1000*1000; i++ {
		result := series.Write(now+uint64(i), writeValue)
		if result.Error != nil {
			t.Error(result.Error)
		}
		if i > minIters && i%100 == 0 {
			if time.Since(startTime).Seconds() > minTime.Seconds() {
				break
			}
		}
	}
	tookMs := float64(time.Since(startTime).Nanoseconds()) / 1000000.0
	tookMsEach := tookMs / float64(i)
	t.Logf("write avg time %f.2ms (%d iterations)", tookMsEach, i)

	c.Close()
	_ = s.Shutdown()
}

func TestReadPerformance(t *testing.T) {
	// start server
	s := NewTestServer(true, true)

	c := NewTestClient()
	now := c.Now()
	series := c.Series("benchmarkSeriesRead")
	const writeValue = 10.1

	// write one value to prevent errors from no data found
	{
		result := series.Write(now, writeValue)
		if result.Error != nil {
			t.Error(result.Error)
		}
	}

	startTime := time.Now()
	const minTime = 1 * time.Second
	const minIters = 100
	const oneMinute = 60 * 1000

	var i int
	for i = 0; i < 1000*1000; i++ {
		result := series.QueryBuilder().From(now - oneMinute).To(now + oneMinute).Execute()
		if result.Error != nil {
			t.Error(result.Error)
		}
		if i > minIters && i%100 == 0 {
			if time.Since(startTime).Seconds() > minTime.Seconds() {
				break
			}
		}
	}
	tookMs := float64(time.Since(startTime).Nanoseconds()) / 1000000.0
	tookMsEach := tookMs / float64(i)
	t.Logf("read avg time %f.2ms (%d iterations)", tookMsEach, i)

	c.Close()
	_ = s.Shutdown()
}

func TestNoOpPerformance(t *testing.T) {
	// start server
	s := NewTestServer(true, true)

	c := NewTestClient()
	series := c.Series("benchmarkSeriesNoOp")

	startTime := time.Now()
	const minTime = 1 * time.Second
	const minIters = 100

	var i int
	for i = 0; i < 1000*1000; i++ {
		err := series.NoOp()
		if err != nil {
			t.Error(err)
		}
		if i > minIters && i%100 == 0 {
			if time.Since(startTime).Seconds() > minTime.Seconds() {
				break
			}
		}
	}
	tookMs := float64(time.Since(startTime).Nanoseconds()) / 1000000.0
	tookMsEach := tookMs / float64(i)
	t.Logf("noop avg time %f.2ms (%d iterations)", tookMsEach, i)

	c.Close()
	_ = s.Shutdown()
}
