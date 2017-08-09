package main

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

var schema = map[string]string{
	"LxdHosts": "CREATE TABLE IF NOT EXISTS 'lxdhosts' (" +
		"'name' VARCHAR(128) PRIMARY KEY, " +
		"'url' VARCHAR(128), " +
		"'status' NUMBER" +
		")",
}

var statements = map[string]string{
	"ReadHosts": "SELECT name, url, status from 'lxdhosts'",
	"AddHost": "INSERT INTO lxdhosts(name, url, status) values(?,?,?)",
	"UpdateHost": "UPDATE lxdhosts set url=?, status=? where name=?",
}

type Storage struct {
	db                 *sql.DB
	preparedStatements map[string]*sql.Stmt
}

func (s* Storage) ReadHosts() LxdHosts {
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

func (s* Storage) AddHost(host* LxdHost)  {
	s.preparedStatements["AddHost"].Exec(host.Name, host.Url, host.Status)
}

func (s* Storage) UpdateHost(host* LxdHost)  {
	s.preparedStatements["UpdateHost"].Exec(host.Url, host.Status, host.Name)
}

func (s*Storage) InitSchema() {
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

func (s*Storage) Open() {
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
