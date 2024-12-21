package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type RoundRobinScheduler struct {
	Servers []Server
	Count   int
}

func NewRoundRobinScheduler(servers []Server) *RoundRobinScheduler {
	scheduler := RoundRobinScheduler{}
	scheduler.Servers = servers
	scheduler.Count = 0
	return &scheduler
}

func (S *RoundRobinScheduler) Delegate(r *http.Request) *http.Response {
	num := (S.Count) % len(S.Servers)
	S.Count++
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
