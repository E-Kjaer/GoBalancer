package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var Servers []Server
var DelegateMap = map[string]Scheduler{}

func main() {
	loadTestServers()
	DelegateMap["r"] = NewRandomScheduler(Servers)
	DelegateMap["rr"] = NewRoundRobinScheduler(Servers)
	DelegateMap["wrr"] = NewWeightedRoundRobinScheduler(Servers)

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
	server1 := Server{"", "127.0.0.1", 5555, 2}
	server2 := Server{"", "127.0.0.1", 6666, 4}
	server3 := Server{"", "127.0.0.1", 7777, 8}
	Servers = append(Servers, server1)
	Servers = append(Servers, server2)
	Servers = append(Servers, server3)
}

func handleTraffic(w http.ResponseWriter, r *http.Request) {
	routingMethod := "wrr"
	res := DelegateMap[routingMethod].Delegate(r)
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.WriteHeader(res.StatusCode)
	w.Write(bytes)
}

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Priority uint8  `json:"priority"`
}

type Scheduler interface {
	Delegate(r *http.Request) *http.Response
}
