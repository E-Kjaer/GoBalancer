package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var Servers []Server

func main() {
	loadTestServers()
	var port int32 = 3333
	http.HandleFunc("/", handleTraffic)

	fmt.Println(fmt.Sprintf("Load Balancer listening on port [:%d]", port))
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func loadConfig() {
	dat, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	var j map[string]interface{}
	err = json.Unmarshal(dat, &j)
	servers := j["servers"].([]interface{})[0]
	fmt.Println(servers)
}

func loadTestServers() {
	server1 := Server{"Server 1", "127.0.0.1", 5555, 10}
	server2 := Server{"Server 2", "127.0.0.1", 6666, 10}
	server3 := Server{"Server 3", "127.0.0.1", 7777, 10}
	Servers = append(Servers, server1)
	Servers = append(Servers, server2)
	Servers = append(Servers, server3)
}

func handleTraffic(w http.ResponseWriter, r *http.Request) {
	num := rand.Int() % len(Servers)
	hostDetails := strings.Split(r.Host, ":")
	fmt.Println(fmt.Sprintf("Host: %s, Port: %d, Server: %d", hostDetails[0], hostDetails[1], num))
}

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Priority uint8  `json:"priority"`
}
