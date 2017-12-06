package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func PostHook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	branch := strings.Split(m["ref"].(string), "/")[2]
	project := m["project"].(map[string]interface{})["name"].(string)
	fmt.Println(string(body))
	fmt.Println("Project->" + project)
	fmt.Println("Branch->" + branch)
	jenkins := CI{}
	job := Job{
		Name:    project + "-" + branch,
		Project: project,
		Group:   branch,
		Status:  JobActive,
	}
	jenkins.connect(GlobalConfiguration["jenkins.url"], GlobalConfiguration["jenkins.username"], GlobalConfiguration["jenkins.password"])
	if m["checkout_sha"] == nil { // Branch deleted
		GlobalStorage.RemoveJob(job)
		jenkins.removeBuild(project, branch)
	} else {
		GlobalStorage.AddJob(job)
		jenkins.updateBuild(project, branch)
	}
}

func GetEnvironments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	environments := GlobalStorage.GetEnvironments(name)
	if err := json.NewEncoder(w).Encode(environments); err != nil {
		panic(err)
	}
}

func AddLxdHost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var result LxdHost
	json.Unmarshal(body, &result)
	GlobalStorage.AddHost(&result)
}

func GetHosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	hosts := GlobalStorage.ReadHosts()
	if err := json.NewEncoder(w).Encode(hosts); err != nil {
		panic(err)
	}
}

func UpdateHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var result LxdHost
	json.Unmarshal(body, &result)
	result.Name = name
	GlobalStorage.UpdateHost(&result)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var result Project
	json.Unmarshal(body, &result)
	GlobalStorage.AddProject(&result)
}

func SetConfigHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var configuration Configuration
	json.Unmarshal(body, &configuration)
	GlobalStorage.SetConfig(configuration)
}

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	config := vars["config"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(GlobalStorage.GetConfig(config)); err != nil {
		panic(err)
	}
}

func GetAllConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(GlobalStorage.GetAllConfig()); err != nil {
		panic(err)
	}
}

func AddResolverRecordHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var resolverRecord ResolverRecord
	json.Unmarshal(body, &resolverRecord)
	GlobalStorage.AddResolverRecord(resolverRecord)
	GlobalResolver.UpdateDatabase()
}

func GetResolverRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fqdn := vars["fqdn"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(GlobalStorage.GetResolverRecord(fqdn)); err != nil {
		panic(err)
	}
}

func GetAllResolverRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(GlobalStorage.GetAllResolverRecords()); err != nil {
		panic(err)
	}
}

func DeleteResolverRecordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fqdn := vars["fqdn"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	GlobalStorage.RemoveResolverRecord(fqdn)
	GlobalResolver.UpdateDatabase()
}
