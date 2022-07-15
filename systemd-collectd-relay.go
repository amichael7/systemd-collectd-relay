package main

import (
	"context"
	"fmt"
	"os"
	"log"
	"strconv"
	"time"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
)

/*
 * systemd-collectd-relay
 *
 * Description:	Connects to systemd and converts status to
 *				collectd-friendly output on STDOUT
 * 
 * Example Usage:
 *		./systemd-collectd-relay ssh postgresql redis
 *
 * See also dbus API docs:
 * 		https://www.freedesktop.org/wiki/Software/systemd/dbus/
 */


func getUnitStatus(
	baseCtx context.Context, 
	conn *dbus.Conn, 
	units []string, 
	hostname string,
) {
	// Note: dbus calls take less than a millisecond to execute in testing
	ctx, cancel := context.WithTimeout(baseCtx, 10 * time.Millisecond)
	defer cancel()
	
	unitNamesLong := make([]string, len(units))
    for ix, unit := range units {
    	unitNamesLong[ix] =  unit + ".service"
    }

    unitStates := make(map[string]float64)
	dbusUnitStateArr, err := conn.ListUnitsByNamesContext(ctx, unitNamesLong)
	if err != nil {
		log.Fatal(err)
	}
	for _, status := range dbusUnitStateArr {
		unitName := status.Name
		if strings.HasSuffix(unitName, ".service") {
			unitName = unitName[:len(unitName) - 8]
		}

		if status.ActiveState == "active" {
			unitStates[unitName] = 1.0
		} else {
			unitStates[unitName] = 0.0
		}
	}

	for unit, state := range unitStates {
		fmt.Printf("PUTVAL %s/systemd/%s N:%f\n", hostname, unit, state)	
	}
}

func main() {
	collectd_interval := os.Getenv("COLLECTD_INTERVAL")
	collectd_hostname := os.Getenv("COLLECTD_HOSTNAME")

	fmt.Println("# COLLECTD_INTERVAL:", collectd_interval)
	fmt.Println("# COLLECTD_HOSTNAME:", collectd_hostname)

	interval := 2
	if "COLLECTD_INTERVAL" != "" {
		parsed, err := strconv.Atoi(collectd_interval)
		if err == nil {
			interval = parsed
		}
	}

	ctx := context.Background()
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for true {
		go getUnitStatus(ctx, conn, os.Args, collectd_hostname)
		fmt.Printf("# Sleeping for %ds...\n", interval)
		time.Sleep(time.Duration(interval) * time.Second)	
	}
}

