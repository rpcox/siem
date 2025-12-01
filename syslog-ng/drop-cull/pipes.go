package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

const (
	ErrNPipeExists = iota
	ErrNPipeRemove
	ErrNPipeMkFifo
	ErrNPipeOpen
)

// Determine if the named pipe already exists. If any error other than ErrNotExist, post the error and exit.  Unable to proceed.
func NamedPipeExists(pipeName string) bool {
	_, err := os.Stat(pipeName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(ErrNPipeExists)
		}
	}

	return true
}

// Open the named pipe syslog-ng will write to and drop-cull will read from. Exit is only option. We can't proceed with any of these errors.
func OpenNamedPipe(pipeName string) *os.File {
	s := "OpenNamedPipe:"
	if NamedPipeExists(pipeName) {
		err := os.Remove(pipeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Remove(%s): %v\n", s, pipeName, err)
			os.Exit(ErrNPipeRemove)
		}
	}

	if err := syscall.Mkfifo(pipeName, 0640); err != nil {
		fmt.Fprintf(os.Stderr, "%s Mkfifo(%s): %v\n", s, pipeName, err)
		os.Exit(ErrNPipeMkFifo)
	}

	fh, err := os.OpenFile(pipeName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s OpenFile(%s): %v\n", s, pipeName, err)
		os.Exit(ErrNPipeOpen)
	}

	fmt.Fprintf(os.Stderr, "%s %s\n", s, pipeName)
	return fh
}
