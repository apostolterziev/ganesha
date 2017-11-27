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
	"Resolver": "CREATE TABLE IF NOT EXISTS 'resolver' (" +
		"'fqdn' VARCHAR(128) PRIMARY KEY, " +
		"'ip' VARCHAR(500) " +
		");" +
		"CREATE UNIQUE INDEX IF NOT EXISTS 'resolver_unique_fqdn' on resolver(fqdn)",
}

var statements = map[string]string{
	"ReadHosts":             "SELECT name, url, status from 'lxdhosts'",
	"AddHost":               "INSERT INTO lxdhosts(name, url, status) values(?, ?, ?)",
	"UpdateHost":            "UPDATE lxdhosts set url=?, status=? where name=?",
	"GetProjects":           "SELECT name, vcs_url, default_branch, job_definition from projects",
	"AddProject":            "INSERT INTO projects(name, vcs_url, default_branch, job_definition) values (?, ?, ?, ?)",
	"GetEnvironments":       "SELECT name, project_name, branch from 'environments' where name=?",
	"StoreEnvironment":      "INSERT INTO environments(name, project_name, branch) values(?, ?, ?)",
	"SetConfigValue":        "INSERT OR REPLACE INTO configuration(name, value) values(?, ?)",
	"GetConfigValue":        "SELECT value from configuration where name=?",
	"GetAllConfig":          "SELECT name, value from configuration",
	"AddResolverRecord":     "INSERT OR REPLACE INTO resolver(fqdn, ip) values(?, ?)",
	"GetFqdnIp":             "SELECT ip from resolver where fqdn=?",
	"GetAllResolverRecords": "SELECT fqdn, ip FROM resolver",
	"RemoveResolverRecord":  "DELETE FROM resolver where fqdn=?",
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

func (s *Storage) GetEnvironments(name string) EnvironmentSet {
	rows, err := s.preparedStatements["GetEnvironments"].Query(name)
	checkErr(err)
	var environments = EnvironmentSet{}
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

func (s *Storage) AddResolverRecord(resolverRecord ResolverRecord) {
	s.preparedStatements["AddResolverRecord"].Exec(resolverRecord.FQDN, resolverRecord.IP)
}

func (s *Storage) GetResolverRecord(fqdn string) ResolverRecord {
	rows, err := s.preparedStatements["GetFqdnIp"].Query(fqdn)
	checkErr(err)
	if rows.Next() {
		resolverRecord := ResolverRecord{FQDN: fqdn}
		rows.Scan(&resolverRecord.IP)
		return resolverRecord
	}
	return ResolverRecord{}
}

func (s *Storage) GetAllResolverRecords() ResolverRecords {
	rows, err := s.preparedStatements["GetAllResolverRecords"].Query()
	checkErr(err)
	resolverRecords := ResolverRecords{}
	for rows.Next() {
		resolverRecord := ResolverRecord{}
		rows.Scan(&resolverRecord.FQDN, &resolverRecord.IP)
		resolverRecords = append(resolverRecords, resolverRecord)
	}
	return resolverRecords
}

func (s *Storage) RemoveResolverRecord(fqdn string) {
	_, err := s.preparedStatements["RemoveResolverRecord"].Exec(fqdn)
	checkErr(err)
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
