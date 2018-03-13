package main

import (
	"log"
	"net/http"
	"fmt"
	"flag"
	"os"
	"strings"
	"path/filepath"
)

var GlobalStorage = Storage{}
var GlobalConfiguration map[string]string
var GlobalResolver = Resolver{}
var GlobalCI = CI{}
var GlobalFlags = CommandLineConfiguration{}

func configurationsToConfig(configurations Configurations) map[string]string {
	var result = make(map[string]string)
	for _, configuration := range configurations {
		result[configuration.Config] = configuration.Value
	}
	return result
}

func parseCommandLine() {
	GlobalFlags.ProcessName = &os.Args[0]
	GlobalFlags.LinkJira = flag.Bool("link-jira", false, "Link up with jira")
	GlobalFlags.AuthCode = flag.String("auth-code", "", "Authentication code")
	GlobalFlags.DatabaseFile = flag.String("database", "", "Database File")
	GlobalFlags.ConfigName = flag.String("config-name", "", "Configuration Name")
	GlobalFlags.ConfigValue = flag.String("config-value", "", "Configuration Value")

	flag.Parse()
}

func main() {
	parseCommandLine()
	if *GlobalFlags.DatabaseFile == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		*GlobalFlags.DatabaseFile = dir + "/ganesha.db"
	}

	GlobalStorage.InitSchema(*GlobalFlags.DatabaseFile)
	GlobalConfiguration = configurationsToConfig(GlobalStorage.GetAllConfig())

	// Set a confing from the command line and exit
	if *GlobalFlags.ConfigName != "" && *GlobalFlags.ConfigValue != "" {
		GlobalStorage.SetConfig(Configuration{
			Config: *GlobalFlags.ConfigName,
			Value:  *GlobalFlags.ConfigValue,
		})
		return
	}

	taskProvider := NewTaskProvider(GlobalConfiguration["jira.key"],
		GlobalConfiguration["jira.url"],
		GlobalConfiguration["jira.consumer_key"])

	if *GlobalFlags.LinkJira {
		authorizationUrl := taskProvider.storeAuthorizationUrl()
		fmt.Println(authorizationUrl)
		return
	}

	if *GlobalFlags.AuthCode != "" {
		taskProvider.storeAccessToken(*GlobalFlags.AuthCode)
		return
	}

	//Work as a git hook
	if strings.HasSuffix(*GlobalFlags.ProcessName, "commit-msg") {
		commitMessage := readFile(os.Args[1])
		issueId := extractIssueFromMessage(commitMessage)
		if issueId == "" {
			fmt.Println("No issue id found in commit message!")
			os.Exit(1)
		}
		jiraClient := taskProvider.getJIRAClient()
		if jiraClient == nil {
			fmt.Println("Cannot connect to jira! Probably a valid jira token is missing!")
			os.Exit(3)
		}
		issue, _, err := jiraClient.Issue.Get(issueId, nil)
		if (err != nil) {
			fmt.Println("Issue " + issueId + " cannot be found in issue tracker!")
			os.Exit(2)
		}
		fmt.Println("Commiting for issue: " + issueId)
		fmt.Println(issue.Fields.Summary)
		os.Exit(0)
	}

	// Start a DNS server in a separate go-routine
	resolverPattern := GlobalConfiguration["resolver.pattern"]
	GlobalResolver.UpdateDatabase()
	if resolverPattern != "" {
		go GlobalResolver.run(resolverPattern)
	}

	if GlobalConfiguration["jenkins.url"] != "" {
		GlobalCI.connect(GlobalConfiguration["jenkins.url"], GlobalConfiguration["jenkins.username"], GlobalConfiguration["jenkins.password"])
	}
	/*lxd := LXDServer{
		Key:        GlobalConfiguration["lxd.key"],
		Cert:       GlobalConfiguration["lxd.certificate"],
		ServerCert: GlobalConfiguration["lxd.server_certificate"],
		Url:        GlobalConfiguration["lxd.url"],
	}*/
	//lxd.Init()
	//lxd.Ping()
	//lxd.Exec("touch /tmp/hello","app1")
	router := NewRouter()
	err := http.ListenAndServe(":8080", router)
	if (err != nil) {
		log.Fatal(err)
	}
}
