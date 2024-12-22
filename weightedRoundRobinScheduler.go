package main

import (
	"net/http"
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
	res := SendRequest(server, r)
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
