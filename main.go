package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var Port int
var Algorithm string
var Servers []Server
var DelegateMap = map[string]Scheduler{}
var ServerMutex = sync.RWMutex{}

func main() {
	loadConfig()
	loadSchedulers()
	http.HandleFunc("/", handleTraffic)

	go runHealthLoop()

	fmt.Println(fmt.Sprintf("Load Balancer listening on port [:%d]", Port))
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", Port), nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func runHealthLoop() {
	for {
		ServerMutex.Lock()
		for i, _ := range Servers {
			_, err := http.Get(fmt.Sprintf("http://%s:%d/health", Servers[i].Host, Servers[i].Port))
			if err != nil {
				fmt.Println(Servers[i].Name + " is Offline")
				// Set server as inactive
				Servers[i].Active = false
			} else {
				if Servers[i].Active == false {
					fmt.Println(Servers[i].Name + " is now Online again")
					Servers[i].Active = true
				}
			}
		}
		// Update schedulers
		updateSchedulers()
		ServerMutex.Unlock()
		time.Sleep(10 * time.Second)
	}
}

func updateSchedulers() {
	DelegateMap["wrr"].UpdateServers()
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
			Active:   true,
		}
		Servers = append(Servers, server)
	}
}

func loadSchedulers() {
	DelegateMap["rand"] = NewRandomScheduler(&Servers)
	DelegateMap["rr"] = NewRoundRobinScheduler(&Servers)
	DelegateMap["wrr"] = NewWeightedRoundRobinScheduler(&Servers)
}

func handleTraffic(w http.ResponseWriter, r *http.Request) {
	ServerMutex.Lock()
	defer ServerMutex.Unlock()
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

func SendRequest(server *Server, r *http.Request) *http.Response {
	r.Host = fmt.Sprintf("%s:%d", server.Host, server.Port)
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

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Priority uint8  `json:"priority"`
	Active   bool
}

type Scheduler interface {
	Delegate(r *http.Request) *http.Response
	UpdateServers()
}
