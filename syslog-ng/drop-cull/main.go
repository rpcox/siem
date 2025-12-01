// proof of concept - syslog-ng forwards records to this program via a named pipe
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var (
	m Metrics
)

func main() {
	_readPipe := flag.String("read-pipe", "/tmp/read-pipe", "Specify the named pipe to read from")
	_writePipe := flag.String("write-pipe", "/tmp/write-pipe", "Specify the named pipe to write to")
	_dropLog := flag.String("log", "/var/log/drop.log", "Specify location of drop log")
	_enable_log := flag.Bool("enable-log", false, "Enable drop log")
	_dropModulus := flag.Uint64("drop-fraction", 10, "Specify the fraction of records to drop")
	flag.Parse()

	fh := InitDropLog(*_enable_log, *_dropLog)
	if fh != nil {
		defer fh.Close()
	}

	writePipe := OpenNamedPipe(*_writePipe)
	defer writePipe.Close()
	readPipe := OpenNamedPipe(*_readPipe)
	defer readPipe.Close()

	reader := bufio.NewReader(readPipe)
	done := make(chan bool)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go sigHandler(sigChan, done)

	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				m.ReadErrorCount++
				if err == io.EOF { // end of pipe
					return
				} else {
					// throw errors and continue
					fmt.Println(err)
					continue
				}
			}

			m.RecordCount++
			if fh != nil && (m.RecordCount%*_dropModulus == 1) {
				fmt.Fprintf(fh, "%3d: %s", m.DropWriteCount, line)
				m.DropWriteCount++
			}
		}

	}()

	<-done
	fmt.Fprintf(os.Stderr, "exiting\n")
	if err := os.Remove(*_readPipe); err != nil {
		fmt.Fprintf(os.Stderr, "remove read pipe on exit: %v\n", err)
	}
	if err := os.Remove(*_writePipe); err != nil {
		fmt.Fprintf(os.Stderr, "remove write pipe on exit: %v\n", err)
	}
}
