package main

import (
	"net/http"
)

type WeightedRoundRobinScheduler struct {
	Servers     *[]Server
	Count       int
	WeightedMap map[int]*Server
}

func NewWeightedRoundRobinScheduler(servers *[]Server) *WeightedRoundRobinScheduler {
	scheduler := WeightedRoundRobinScheduler{}
	scheduler.Servers = servers
	scheduler.CreateWeightMap()
	scheduler.Count = 0
	return &scheduler
}

func (S *WeightedRoundRobinScheduler) changeServer() *Server {
	num := (S.Count) % len(S.WeightedMap)
	server := S.WeightedMap[num]
	S.Count++
	for !server.Active {
		num = (S.Count) % len(S.WeightedMap)
		server = S.WeightedMap[num]
		S.Count++
	}
	return server
}

func (S *WeightedRoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	server := S.changeServer()
	res := SendRequest(server, r)
	return res
}

func (S *WeightedRoundRobinScheduler) UpdateServers() {
	S.CreateWeightMap()
}

func (S *WeightedRoundRobinScheduler) CreateWeightMap() {
	count := 0
	result := map[int]*Server{}
	for index, element := range *S.Servers {
		for _ = range (*S.Servers)[index].Priority {
			result[count] = &element
			count++
		}
	}
	S.WeightedMap = result
}
