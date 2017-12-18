package main

import (
	"io/ioutil"
	"regexp"
)

var issueRe, _ = regexp.Compile(`^\[(.+)\].*`)

func readFile(filePath string) string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func extractIssueFromMessage(message string) string {
	issues := issueRe.FindStringSubmatch(message)
	if len(issues) == 2 {
		return issues[1]
	}
	return ""
}
