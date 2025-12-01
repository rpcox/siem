package main

import (
	"fmt"
	"os"
	"syscall"
)

func sigHandler(sigChan chan os.Signal, done chan bool) {
	for sig := range sigChan {
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			fmt.Fprintln(os.Stderr, "received: terminating signal")
			close(done)
			return
		} else if sig == syscall.SIGHUP {
			fmt.Fprintln(os.Stderr, "received: SIGHUP")
		} else if sig == syscall.SIGUSR1 {
			fmt.Fprintln(os.Stderr, "received: SIGUSR1")
		} else if sig == syscall.SIGUSR2 {
			fmt.Fprintln(os.Stderr, "received: SIGUSR2")
		}
	}
}

func InitDropLog(enable bool, logName string) *os.File {
	var fh *os.File
	var err error

	if enable {
		fh, err = os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "drop log not enabled: %v\n", err)
			return nil
		}
	}

	return fh
}
