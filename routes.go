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
		"GetEnvironment",
		"GET",
		"/environments/{name}",
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
		"UpdateHost",
		"PUT",
		"/hosts/{name}",
		UpdateHost,
	},
	Route{
		"AddProject",
		"POST",
		"/projects",
		AddProject,
	},
	Route{
		"GetConfig",
		"GET",
		"/config/{config}",
		GetConfigHandler,
	},
	Route{
		"GetAllConfig",
		"GET",
		"/config",
		GetAllConfigHandler,
	},
	Route{
		"SetConfig",
		"POST",
		"/config",
		SetConfigHandler,
	},
}