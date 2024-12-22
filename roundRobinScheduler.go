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

func (S *RoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	num := (S.Count) % len(*S.Servers)
	server := &(*S.Servers)[num]
	S.Count++
	for !server.Active {
		num = (S.Count) % len(*S.Servers)
		server = &(*S.Servers)[num]
		S.Count++
	}
	res := SendRequest(server, r)
	return res
}

func (S *RoundRobinScheduler) UpdateServers() {

}
