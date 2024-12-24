package main

import (
	"math/rand"
	"net/http"
)

type RandomScheduler struct {
	Servers *[]Server
}

func NewRandomScheduler(servers *[]Server) *RandomScheduler {
	scheduler := RandomScheduler{}
	scheduler.Servers = servers
	return &scheduler
}

func (S *RandomScheduler) changeServer() *Server {
	num := rand.Int() % len(*S.Servers)
	server := &(*S.Servers)[num]
	for !server.Active {
		num = rand.Int() % len(*S.Servers)
		server = &(*S.Servers)[num]
	}
	return server
}

func (S *RandomScheduler) Delegate(r *http.Request) *http.Response {
	server := S.changeServer()
	res := SendRequest(server, r)
	return res
}

func (S *RandomScheduler) UpdateServers() {

}
