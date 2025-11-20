// filter example
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

	fh, err := os.OpenFile("/var/log/program.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "program exit %v\n", err)
		os.Exit(1)
	}
	// can also use (stat.Mode() & os.ModeNamedPipe) != 0
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// something is piping data to us on stdin
		reader := bufio.NewReader(os.Stdin)
		count := 1

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				} else {
					fmt.Println(err)
					continue
				}
			}

			fmt.Fprintf(fh, "%3d: %s", count, line)
			count++
		}
	} else { // not using stdin. assume a list of files on the cmd line
		for _, file := range os.Args[1:] {
			fh, err := os.Open(file)
			if err != nil {
				fmt.Println(file, err)
				continue
			}

			fmt.Println("--- ", file)
			count := 1

			reader := bufio.NewReader(fh)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					} else {
						fmt.Println(err)
						os.Exit(1)
					}
				}

				fmt.Printf("%3d: %s", count, line)
				count++
			}

			fh.Close()
		}
	}
}
