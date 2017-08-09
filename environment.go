package main

type Environment struct {
	Name string `json:"name"`
}

type LxdHost struct {
	Name string `json:"name"`
	Status int8 `json:"status"`
	Url string  `json:"url"`
}

type Environments []Environment
type LxdHosts []LxdHost