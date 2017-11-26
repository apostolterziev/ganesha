package main

import (
	"log"
	"net/http"
)

var GlobalStorage = Storage{}
var GlobalConfiguration map[string]string
var GlobalResolver = Resolver{}

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
	GlobalResolver.UpdateDatabase()
	resolverPattern := GlobalConfiguration["resolver.pattern"]
	if resolverPattern != "" {
		go GlobalResolver.run(resolverPattern)
	}
	router := NewRouter()
	jenkins := CI{}
	jenkins.connect(GlobalConfiguration["jenkins.url"], GlobalConfiguration["jenkins.username"], GlobalConfiguration["jenkins.password"])
	lxd := LXDServer{
		Key:  GlobalConfiguration["lxd.key"],
		Cert: GlobalConfiguration["lxd.certificate"],
		Url:  GlobalConfiguration["lxd.url"],
	}
	lxd.Init()
	//lxd.Ping()
	//lxd.Exec("touch /tmp/hello","app1")

	log.Fatal(http.ListenAndServe(":8080", router))
}
