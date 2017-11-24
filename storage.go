package main

import "database/sql"
import (
	_ "github.com/mattn/go-sqlite3"
)

var schema = map[string]string{
	"LxdHosts": "CREATE TABLE IF NOT EXISTS 'lxdhosts' (" +
		"'name' VARCHAR(128) PRIMARY KEY, " +
		"'url' VARCHAR(128), " +
		"'status' NUMBER" +
		")",
	"Environments": "CREATE TABLE IF NOT EXISTS 'environments' (" +
		"'name' VARCHAR(128) PRIMARY KEY, " +
		"'project_name' VARCHAR(128), " +
		"'branch' VARCHAR (128)" +
		")",
	"Projects": "CREATE TABLE IF NOT EXISTS 'projects' (" +
		"'name' VARCHAR(128) PRIMARY KEY, " +
		"'vcs_url' VARCHAR(500), " +
		"'default_branch' VARCHAR(128)," +
		"'job_definition' VARCHAR" +
		")",
	"Configuration": "CREATE TABLE IF NOT EXISTS 'configuration' (" +
		"'name' VARCHAR(128) PRIMARY KEY, " +
		"'value' VARCHAR(500) " +
		")",
}

var statements = map[string]string{
	"ReadHosts":        "SELECT name, url, status from 'lxdhosts'",
	"AddHost":          "INSERT INTO lxdhosts(name, url, status) values(?, ?, ?)",
	"UpdateHost":       "UPDATE lxdhosts set url=?, status=? where name=?",
	"GetProjects":      "SELECT name, vcs_url, default_branch, job_definition from projects",
	"AddProject":       "INSERT INTO projects(name, vcs_url, default_branch, job_definition) values (?, ?, ?, ?)",
	"GetEnvironments":  "SELECT name, project_name, branch from 'environments' where name=?",
	"StoreEnvironment": "INSERT INTO environments(name, project_name, branch) values(?, ?, ?)",
	"SetConfigValue":   "INSERT OR REPLACE INTO configuration(name, value) values(?, ?)",
	"GetConfigValue":   "SELECT value from configuration where name=?",
	"GetAllConfig":     "SELECT name, value from configuration",
}

type Storage struct {
	db                 *sql.DB
	preparedStatements map[string]*sql.Stmt
}

func (s *Storage) ReadHosts() LxdHosts {
	rows, err := s.preparedStatements["ReadHosts"].Query()
	checkErr(err)
	var hosts = LxdHosts{}
	for rows.Next() {
		var host = LxdHost{}
		rows.Scan(&host.Name, &host.Url, &host.Status)
		hosts = append(hosts, host)
	}
	rows.Close()
	return hosts
}

func (s *Storage) AddHost(host *LxdHost) {
	s.preparedStatements["AddHost"].Exec(host.Name, host.Url, host.Status)
}

func (s *Storage) UpdateHost(host *LxdHost) {
	s.preparedStatements["UpdateHost"].Exec(host.Url, host.Status, host.Name)
}

func (s *Storage) GetProjects() Projects {
	rows, err := s.preparedStatements["GetProjects"].Query()
	checkErr(err)
	projects := Projects{}
	for rows.Next() {
		var project = Project{}
		rows.Scan(&project.Name, &project.VCSUrl, &project.DefaultBranch, &project.JobDefinition)
		projects = append(projects, project)
	}
	return projects
}

func (s *Storage) GetConfig(config string) Configuration {
	rows, err := s.preparedStatements["GetConfigValue"].Query(config)
	checkErr(err)
	if rows.Next() {
		configuration := Configuration{Config: config}
		rows.Scan(&configuration.Value)
		return configuration
	}
	return Configuration{}
}

func (s *Storage) GetAllConfig() Configurations {
	rows, err := s.preparedStatements["GetAllConfig"].Query()
	checkErr(err)
	configs := Configurations{}
	for rows.Next() {
		configuration := Configuration{}
		rows.Scan(&configuration.Config, &configuration.Value)
		configs = append(configs, configuration)
	}
	return configs
}


func (s *Storage) SetConfig(configuration Configuration) {
	s.preparedStatements["SetConfigValue"].Exec(configuration.Config, configuration.Value)
}

func (s *Storage) AddProject(project *Project) {
	s.preparedStatements["AddProject"].Exec(project.Name, project.VCSUrl, project.DefaultBranch, project.JobDefinition)
}

func (s *Storage) GetEnvironments(name string) Environments {
	rows, err := s.preparedStatements["GetEnvironments"].Query(name)
	checkErr(err)
	var environments = Environments{}
	for rows.Next() {
		var environment = Environment{}
		rows.Scan(&environment.Name, &environment.ProjectName, &environment.Branch)
		environments = append(environments, environment)
	}
	return environments
}

func (s *Storage) StoreEnvironment(environment *Environment) {
	s.preparedStatements["StoreEnvironment"].Exec(environment.Name, environment.ProjectName)
}

func (s *Storage) InitSchema() {
	s.Open()

	for _, sql := range schema {
		_, err := s.db.Exec(sql)
		checkErr(err)
	}

	s.preparedStatements = make(map[string]*sql.Stmt)
	for statementKey, statement := range statements {
		stmt, err := s.db.Prepare(statement)
		checkErr(err)
		s.preparedStatements[statementKey] = stmt
	}
}

func (s *Storage) Open() {
	if s.db != nil {
		return
	}
	var err error
	s.db, err = sql.Open("sqlite3", "./ganesha.db")
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
