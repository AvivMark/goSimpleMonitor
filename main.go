package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Host struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

func main() {
	hosts := loadHosts("./hosts.json")
	ips := getIps(hosts)
	go startWebServer()
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

func loadHosts(p string) (hosts []Host) {
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

func getIps(hl []Host) []string {
	ips := make([]string, 0)
	for _, host := range hl {
		ips = append(ips, host.IP)
		fmt.Println(ips)
	}
	return ips
}

func printHost(h Host) {
	fmt.Println(h.IP)
	fmt.Println(h.Name)
}
