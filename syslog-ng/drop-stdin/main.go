// proof of concept - syslog-ng forwards records to this program via stdin (through a shell)
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// See if we have a pipe sending us data
	// can also use (stat.Mode() & os.ModeNamedPipe) != 0
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		count := 1

		fh, err := os.OpenFile("/var/log/drop.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "drop exit %v\n", err)
			os.Exit(1)
		}
		defer fh.Close()

		for {
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
	} else {
		fmt.Fprintln(os.Stderr, "could not obtain pipe on stdin")
	}
}
