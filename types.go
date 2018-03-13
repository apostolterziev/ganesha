package main

type CommandLineConfiguration struct {
	LinkJira *bool
	AuthCode *string
	DatabaseFile *string
	ProcessName *string
	ConfigName *string
	ConfigValue *string
}

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

type Job struct {
	Name    string
	Project string
	Group   string
	Status  JobStatus
}

type OauthConfiguration struct {
	URL              string
	AuthorizationURL string
	RequestSecret    string
	RequestToken     string
	Token            string
	TokenSecret      string
}

type JobStatus int8

type LXDServer struct {
	Name       string
	Key        string
	Cert       string
	ServerCert string
	Url        string
}

type LXDPool struct {
	Pool map[string]*LXDServer
}

const (
	JobActive  JobStatus = 0
	JobRunning JobStatus = 1
)

type EnvironmentSet []Environment
type LxdHosts []LxdHost
type Projects []Project
type Configurations []Configuration
type ResolverRecords []ResolverRecord
type Jobs []Job
