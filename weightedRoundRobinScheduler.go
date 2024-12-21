package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type WeightedRoundRobinScheduler struct {
	Servers     []Server
	Count       int
	WeightedMap map[int]*Server
}

func NewWeightedRoundRobinScheduler(servers []Server) *WeightedRoundRobinScheduler {
	scheduler := WeightedRoundRobinScheduler{}
	scheduler.Servers = servers
	scheduler.WeightedMap = createWeightMap(servers)
	scheduler.Count = 0
	return &scheduler
}

func (S *WeightedRoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	num := (S.Count) % len(S.WeightedMap)
	server := S.WeightedMap[num]
	S.Count++
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

func createWeightMap(servers []Server) map[int]*Server {
	count := 0
	result := map[int]*Server{}
	for index, element := range servers {
		for _ = range servers[index].Priority {
			result[count] = &element
			count++
		}
	}
	return result
}
