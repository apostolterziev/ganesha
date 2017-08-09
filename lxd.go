package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"time"
)

type LXDServer struct {
	Name      string
	Key, Cert string
	Url       string
	Transport *http.Transport
}

type LXDPool struct {
	Pool map[string]*LXDServer
}

func (s *LXDServer) Init() {
	cer, err := tls.LoadX509KeyPair(s.Cert, s.Key)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	s.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cer},
		},
		MaxIdleConns:2,
		IdleConnTimeout: 3* time.Second,
	}
}

func (s *LXDServer) Ping() {

	client := &http.Client{Transport: s.Transport}
	resp, err := client.Get(s.Url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func (s *LXDServer) GetOperations() {
	client := &http.Client{Transport: s.Transport}
	resp, err := client.Get(s.Url + "/operations")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}