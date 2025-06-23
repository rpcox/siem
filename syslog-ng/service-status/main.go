package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/rpcox/pkg/exit"
)

const (
	ExitMatch = iota
	ExitDbusNotFound
	ExitDbusConnErr
	ExitDbusListErr
	ExitRegexpErr
)

func main() {
	_timeout := flag.Int("to", 3, "Specify the timeout for the dbus connection in seconds")
	_unitName := flag.String("name", "dbus", "Specify the systemd unit name")
	_unitType := flag.String("type", "service", "Specify the systemd unit type")
	//	_version := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*_timeout)*time.Second)
	defer cancel()

	dbconn, err := dbus.NewSystemConnectionContext(ctx)
	exit.IfErr(err != nil, err, ExitDbusConnErr)
	defer dbconn.Close()

	units, err := dbconn.ListUnitsContext(ctx)
	exit.IfErr(err != nil, err, ExitDbusListErr)

	re, err := regexp.Compile(`^` + *_unitName + `.` + *_unitType)
	exit.IfErr(err != nil, err, ExitRegexpErr)

	for _, unit := range units {
		if match := re.MatchString(unit.Name); match {
			fmt.Fprintf(os.Stdout, " Name: %s\n", unit.Name)
			fmt.Fprintf(os.Stdout, "State: %s %s %s\n", unit.LoadState, unit.ActiveState, unit.SubState)
			os.Exit(ExitMatch)
		}
	}

	os.Exit(ExitDbusNotFound)
}
