package main

import (
	"log"
	"net/http"
)

var GlobalStorage = Storage{}

func main() {
	router := NewRouter()
	//jenkins := Jenkins{}
	//jenkins.connect()
	//lxd := LXDServer{
	//	Key:  "/home/api/tmp/ca/client.key",
	//	Cert: "/home/api/tmp/ca/client.crt",
	//	Url:  "https://lxd-test:8443/1.0",
	//}
	//lxd.Init()
	//lxd.Ping()
	//lxd.GetOperations()

	GlobalStorage.InitSchema()
	log.Fatal(http.ListenAndServe(":8080", router))
}
