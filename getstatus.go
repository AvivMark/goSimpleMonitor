package main

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func startMonitor(ips []string) {

	completed := make(chan string)
	for _, ip := range ips {
		go monitorIP(ip, completed)

	}

	for ip := range completed {
		go func(ipAddress string) {
			time.Sleep(2 * time.Second)
			monitorIP(ipAddress, completed)
		}(ip)
	}

}

func monitorIP(ip string, c chan string) {
	pinger, err := ping.NewPinger(ip)
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Duration(time.Millisecond * 300)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		c <- err.Error()
	}
	if pinger.PacketsRecv > 0 {
		fmt.Println(ip + " is working")
		c <- ip
	} else {
		fmt.Println("failed to connect to " + ip)
		c <- ip
	}
}
