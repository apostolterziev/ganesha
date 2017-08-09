package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"PostHook",
		"POST",
		"/",
		PostHook,
	},
	Route{
		"GetEnvironments",
		"GET",
		"/environments",
		GetEnvironments,
	},
	Route{
		"AddLxdHost",
		"POST",
		"/hosts",
		AddLxdHost,
	},
	Route{
		"GetHosts",
		"GET",
		"/hosts",
		GetHosts,
	},
	Route{
		"GetHosts",
		"PUT",
		"/hosts/{name}",
		UpdateHost,
	},
}