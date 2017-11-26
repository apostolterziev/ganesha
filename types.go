package main

type Environment struct {
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	Branch      string `json:"branch"`
}

type LxdHost struct {
	Name   string `json:"name"`
	Status int8   `json:"status"`
	Url    string `json:"url"`
}

type Project struct {
	Name          string `json:"name"`
	VCSUrl        string `json:"vcs_url"`
	DefaultBranch string `json:"default_branch"`
	JobDefinition string `json:"job_definition"`
}

type Configuration struct {
	Config string `json:"config"`
	Value  string `json:"value"`
}

type ResolverRecord struct {
	FQDN string `json:"fqdn"`
	IP   string `json:"ip"`
}

type Environments []Environment
type LxdHosts []LxdHost
type Projects []Project
type Configurations []Configuration
type ResolverRecords []ResolverRecord
