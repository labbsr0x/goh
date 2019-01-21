package gohcmd

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestGracefulStop(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	rch := make(chan bool)

	go func() {
		osExit = func(code int) { rch <- true; return }
		GracefulStop(cancel)
	}()

	time.Sleep(10 * time.Millisecond)

	stopCh <- syscall.SIGTERM
	<-rch
}
