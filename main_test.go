package openvpn

import (
	"log/slog"
	"testing"
	"time"
)

func TestCheckConnection(t *testing.T) {
	closeCh := make(chan int, 2)

	res := checkConnection(closeCh)
	time.Sleep(16*time.Second)
	
	if res != "ok" {
		panic("Alarm!!!")
	}

	slog.Info("OK!")
}