package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

var Servers []Server
var delegateMap = map[string]Scheduler{}

func main() {
	loadTestServers()
	var port int32 = 3333
	http.HandleFunc("/", handleTraffic)
	delegateMap["r"] = RandomScheduler{Servers}
	delegateMap["rr"] = RoundRobinScheduler{Servers}
	delegateMap["wrr"] = WeightedRoundRobinScheduler{Servers}

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
	server1 := Server{"", "127.0.0.1", 5555, 10}
	server2 := Server{"", "127.0.0.1", 6666, 10}
	server3 := Server{"", "127.0.0.1", 7777, 10}
	Servers = append(Servers, server1)
	Servers = append(Servers, server2)
	Servers = append(Servers, server3)
}

func handleTraffic(w http.ResponseWriter, r *http.Request) {
	routingMethod := "r"
	res := delegateMap[routingMethod].Delegate(r)
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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

type RandomScheduler struct {
	Servers []Server
}

func (S RandomScheduler) Delegate(r *http.Request) *http.Response {
	num := rand.Int() % len(S.Servers)
	r.Host = fmt.Sprintf("%s:%d", S.Servers[num].Host, S.Servers[num].Port)
	u, err := url.Parse(fmt.Sprintf("http://%s%s", r.Host, r.RequestURI))
	if err != nil {
		panic(err)
	}
	r.URL = u
	r.RequestURI = ""
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	return res
}

type RoundRobinScheduler struct {
	Servers []Server
}

func (S RoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	num := rand.Int() % len(S.Servers)
	r.Host = fmt.Sprintf("%s:%d", S.Servers[num].Host, S.Servers[num].Port)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	return res
}

type WeightedRoundRobinScheduler struct {
	Servers []Server
}

func (S WeightedRoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	num := rand.Int() % len(S.Servers)
	r.Host = fmt.Sprintf("%s:%d", S.Servers[num].Host, S.Servers[num].Port)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	return res
}
