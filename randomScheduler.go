package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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
