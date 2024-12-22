package main

import (
	"math/rand"
	"net/http"
)

type RandomScheduler struct {
	Servers []Server
}

func NewRandomScheduler(servers []Server) *RandomScheduler {
	scheduler := RandomScheduler{}
	scheduler.Servers = servers
	return &scheduler
}

func (S *RandomScheduler) Delegate(r *http.Request) *http.Response {
	num := rand.Int() % len(S.Servers)
	server := &S.Servers[num]
	res := SendRequest(server, r)
	return res
}
