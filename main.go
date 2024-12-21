package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var Port int
var Algorithm string
var Servers []Server
var DelegateMap = map[string]Scheduler{}

func main() {
	loadConfig()
	loadSchedulers()
	//var port int32 = 3333
	http.HandleFunc("/", handleTraffic)

	fmt.Println(fmt.Sprintf("Load Balancer listening on port [:%d]", Port))
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", Port), nil)
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

	Port = int(j["lb_port"].(float64))

	Algorithm = j["lb_algo"].(string)

	servers := j["be_servers"].([]interface{})
	for _, element := range servers {
		obj := element.(map[string]interface{})
		server := Server{
			Name:     obj["name"].(string),
			Host:     obj["host"].(string),
			Port:     uint16(obj["port"].(float64)),
			Priority: uint8(obj["priority"].(float64)),
		}
		Servers = append(Servers, server)
	}
}

func loadSchedulers() {
	DelegateMap["rand"] = NewRandomScheduler(Servers)
	DelegateMap["rr"] = NewRoundRobinScheduler(Servers)
	DelegateMap["wrr"] = NewWeightedRoundRobinScheduler(Servers)
}

func handleTraffic(w http.ResponseWriter, r *http.Request) {
	res := DelegateMap[Algorithm].Delegate(r)
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
