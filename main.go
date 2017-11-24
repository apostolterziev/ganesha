package main

import (
	"log"
	"net/http"
)

var GlobalStorage = Storage{}
var GlobalConfiguration map[string]string

func ConfigurationsToConfig(configurations Configurations) map[string]string {
	var result = make(map[string]string)
	for _, configuration := range configurations {
		result[configuration.Config] = configuration.Value
		log.Println(configuration.Config + "\t\t\t" + configuration.Value)
	}
	return result
}

func main() {
	GlobalStorage.InitSchema()
	GlobalConfiguration = ConfigurationsToConfig(GlobalStorage.GetAllConfig())
	router := NewRouter()
	jenkins := CI{}
	jenkins.connect(GlobalConfiguration["jenkins.url"], GlobalConfiguration["jenkins.username"], GlobalConfiguration["jenkins.password"])
	lxd := LXDServer{
		Key:  "/home/api/tmp/ca/client.key",
		Cert: "/home/api/tmp/ca/client.crt",
		Url:  GlobalConfiguration["lxd.url"],
	}
	lxd.Init()
	//lxd.Ping()
	//lxd.Exec("touch /tmp/hello","app1")

	log.Fatal(http.ListenAndServe(":8080", router))
}
