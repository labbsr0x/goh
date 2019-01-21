package gohcmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var stopCh = make(chan os.Signal)

// GracefulStop cancels gracefully the running goRoutines
func GracefulStop(cancel context.CancelFunc) {
	signal.Notify(stopCh, syscall.SIGTERM)
	signal.Notify(stopCh, syscall.SIGINT)

	<-stopCh // waits for a stop signal
	stop(0, cancel)
}

var osExit = os.Exit

// stop stops this program
func stop(returnCode int, cancel context.CancelFunc) {
	logrus.Infof("Stopping execution...")
	cancel()
	time.Sleep(2 * time.Second)
	osExit(returnCode)
}
