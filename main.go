package main

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func main() {
	ips := []string{
		"10.10.10.10",
		"1.1.1.2",
		"ansible-server.mark.home",
		"jenkins-server.mark.home",
		"rancher-server.mark.home",
	}
	completed := make(chan string)
	for _, ip := range ips {
		go monitor(ip, completed)

	}

	for ip := range completed {
		go func(ipAddress string) {
			time.Sleep(2 * time.Second)
			monitor(ipAddress, completed)
		}(ip)
	}
}

func monitor(ip string, c chan string) {
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
