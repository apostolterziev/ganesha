package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
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
	fmt.Println(string(body))
}

func GetEnvironments(w http.ResponseWriter, r *http.Request) {
	environments := Environments{
		Environment{Name: "reporting"},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
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

func UpdateHost (w http.ResponseWriter, r *http.Request) {
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