package main

import (
	"net/http"
)

type RoundRobinScheduler struct {
	Servers *[]Server
	Count   int
}

func NewRoundRobinScheduler(servers *[]Server) *RoundRobinScheduler {
	scheduler := RoundRobinScheduler{}
	scheduler.Servers = servers
	scheduler.Count = 0
	return &scheduler
}

func (S *RoundRobinScheduler) changeServer() *Server {
	num := (S.Count) % len(*S.Servers)
	server := &(*S.Servers)[num]
	S.Count++
	for !server.Active {
		num = (S.Count) % len(*S.Servers)
		server = &(*S.Servers)[num]
		S.Count++
	}
	return server
}

func (S *RoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	server := S.changeServer()
	res := SendRequest(server, r)
	return res
}

func (S *RoundRobinScheduler) UpdateServers() {

}
