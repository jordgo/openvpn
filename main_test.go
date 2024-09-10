package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLongCheckConnection(t *testing.T) {
	t.Parallel()

	maxDuration := 27*time.Second
	done := make(chan string)

	go func() {
		closeCh := make(chan int, 2)
		ips := []string{"85.174.51.235"}
		// emptyIps := []string{}
		res := checkConnection(closeCh, ips)
		done <- res

	}()

	select {
	case res := <- done:
		exp := "ok"
		assert.EqualValues(t, exp, res, "Success")

	case <- time.After(maxDuration):
		t.Errorf("Test took longer than %v", maxDuration)
	}

}

func TestSome(t *testing.T) {
	assert.True(t, true, "Stub")
}