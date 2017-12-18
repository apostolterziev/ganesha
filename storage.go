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
		"'ip' VARCHAR(15) " +
		");" +
		"CREATE UNIQUE INDEX IF NOT EXISTS 'resolver_unique_fqdn' on resolver(fqdn)",
	"Jobs": "CREATE TABLE IF NOT EXISTS 'jobs' (" +
		"'name' VARCHAR(128) PRIMARY KEY, " +
		"'project' VARCHAR(128), " +
		"'group' VARCHAR(128), " +
		"'status' INTEGER " +
		")",
	"OAuth1": "CREATE TABLE IF NOT EXISTS 'oauth1' (" +
		"'url' VARCHAR(1024) PRIMARY KEY, " +
		"'authorization_url' VARCHAR(1024) DEFAULT '', " +
		"'request_secret' VARCHAR(1024) DEFAULT '', " +
		"'request_token' VARCHAR(1024) DEFAULT '', " +
		"'token' VARCHAR(1024) DEFAULT '', " +
		"'token_secret' VARCHAR(1024) DEFAULT '' " +
		")",
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
	"AddJob":                "INSERT OR REPLACE INTO jobs(name, project, 'group', status) values(?, ?, ?, ?)",
	"RemoveJob":             "DELETE FROM jobs WHERE name=?",
	"GetAllForGroup":        "SELECT name, project, 'group', status FROM jobs WHERE 'group'=?",
	"AddAuthorizationUrl":   "INSERT OR REPLACE INTO oauth1('url', 'authorization_url', 'request_secret', request_token) values(?, ?, ?, ?)",
	"StoreToken":            "INSERT OR REPLACE INTO oauth1('url', 'token', 'token_secret') values(?, ?, ?)",
	"GetOauthConfiguration": "SELECT url, authorization_url, request_secret, request_token, token, token_secret FROM oauth1 WHERE url=?",
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
	var resolverRecord ResolverRecord
	if rows.Next() {
		resolverRecord = ResolverRecord{FQDN: fqdn}
		rows.Scan(&resolverRecord.IP)
		return resolverRecord
	}
	return resolverRecord
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

func (s *Storage) AddJob(job Job) {
	_, err := GlobalStorage.preparedStatements["AddJob"].Exec(job.Name, job.Project, job.Group, job.Status)
	checkErr(err)
}

func (s *Storage) RemoveJob(job Job) {
	_, err := GlobalStorage.preparedStatements["RemoveJob"].Exec(job.Name)
	checkErr(err)
}

func (s *Storage) GetAllGroupJobs(group string) map[string]Job {
	rows, err := s.preparedStatements["GetAllForGroup"].Query(group)
	checkErr(err)
	jobs := make(map[string]Job)
	for rows.Next() {
		job := Job{}
		rows.Scan(&job.Name, &job.Project, &job.Group, &job.Status)
		jobs[job.Name] = job
	}
	return jobs
}

func (s *Storage) readOauthConfiguration(url string) OauthConfiguration {
	statement := s.preparedStatements["GetOauthConfiguration"]
	//defer statement.Close()
	rows, err := statement.Query(url)
	checkErr(err)
	defer rows.Close()
	var oauthConfiguration OauthConfiguration
	if rows.Next() {
		oauthConfiguration = OauthConfiguration{}
		rows.Scan(&oauthConfiguration.URL, &oauthConfiguration.AuthorizationURL, &oauthConfiguration.RequestSecret,
			&oauthConfiguration.RequestToken, &oauthConfiguration.Token, &oauthConfiguration.TokenSecret)
	}
	return oauthConfiguration
}

func (s *Storage) AddAuthorizationUrl(oauthConfiguration OauthConfiguration) {
	statement := GlobalStorage.preparedStatements["AddAuthorizationUrl"]
	//defer statement.Close()

	_, err := statement.Exec(oauthConfiguration.URL,
		oauthConfiguration.AuthorizationURL,
		oauthConfiguration.RequestSecret,
		oauthConfiguration.RequestToken)
	checkErr(err)
}

func (s *Storage) StoreToken(oauthConfiguration OauthConfiguration) {
	statement := GlobalStorage.preparedStatements["StoreToken"]
	//defer statement.Close()
	_, err := statement.Exec(oauthConfiguration.URL,
		oauthConfiguration.Token,
		oauthConfiguration.TokenSecret)
	checkErr(err)
}

func (s *Storage) InitSchema(database string) {
	s.Open(database)

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

func (s *Storage) Open(database string) {
	if s.db != nil {
		return
	}
	var err error
	s.db, err = sql.Open("sqlite3", database)
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
