package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-ping/ping"
)

type Host struct {
	ID       string `json:"ID"`       // ID
	Hostname string `json:"Hostname"` // HOSTNAME
	HostIP   string `json:"HostIP"`   // HOST IP
	IsAlive  bool   `json:"isAlive"`  // say if there is connection to this host
}

var urlGetHosts string = "http://localhost:5000/hosts"

func main() {
	hosts := loadHostsAPI(urlGetHosts)
	ips := getIps(hosts)

	startMonitor(ips)
}

func testLoad(p string) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	var h Host
	err = json.Unmarshal(data, &h)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", h)

}

func loadHostsJson(p string) (hosts []Host) {
	content, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var tmp *[]Host
	err = json.Unmarshal(content, &tmp)

	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return *tmp
}

func loadHostsAPI(url string) (hosts []Host) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var tmp *[]Host
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return *tmp
}

func getIps(hl []Host) []string {
	ips := make([]string, 0)
	for _, host := range hl {
		ips = append(ips, host.HostIP)
		fmt.Println(ips)
	}
	return ips
}

func printHost(h Host) {
	fmt.Println(h.HostIP)
	fmt.Println(h.Hostname)
}

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
