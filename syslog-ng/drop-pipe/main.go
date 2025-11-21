// proof of concept - syslog-ng forwards records to this program via a named pipe
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func PipeExists(pipeName string) bool {
	_, err := os.Stat(pipeName)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	fmt.Println("named pipe Stat():", err)
	return false
}

func sigHandler(sigChan chan os.Signal, done chan interface{}) {
	for sig := range sigChan {
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			close(done)
		}
	}
}

func main() {
	fh, err := os.OpenFile("/var/log/drop.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "drop exit %v\n", err)
		os.Exit(1)
	}
	defer fh.Close()

	pipeName := "/tmp/mypipe"
	if PipeExists(pipeName) {
		err := os.Remove(pipeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "remove pipe: %v\n", err)
			os.Exit(1)
		}
	}

	err = syscall.Mkfifo(pipeName, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mkfifo(): %v\n", err)
		os.Exit(1)
	}

	inPipe, err := os.OpenFile(pipeName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open pipe: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(inPipe)
	count := 1
	done := make(chan interface{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go sigHandler(sigChan, done)

	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "exit main: %v\n", err)
			inPipe.Close()
			if err = os.Remove(pipeName); err != nil {
				fmt.Fprintf(os.Stderr, "remove pipe on exit: %v\n", err)
			}
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF { // end of pipe
					return
				} else {
					// throw errors and continue
					fmt.Println(err)
					continue
				}
			}

			// write the records to drop to /var/log/drop.log
			fmt.Fprintf(fh, "%3d: %s", count, line)
			count++
		}
	}

}
