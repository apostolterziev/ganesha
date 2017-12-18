package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"github.com/andygrunwald/go-jira"
	"github.com/dghubble/oauth1"
	"golang.org/x/net/context"
	"log"
)

type TaskProveder struct {
	JiraUrl *url.URL
	config  *oauth1.Config
}

func NewTaskProvider(pk string, jiraUrl string, consumerKey string) *TaskProveder {
	jiraUrlParsed, _ := url.Parse(jiraUrl)
	keyDERBlock, _ := pem.Decode([]byte(pk))
	if keyDERBlock == nil {
		log.Fatal("unable to decode key PEM block")
	}
	if !(keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY")) {
		log.Fatalf("unexpected key DER block type: %s", keyDERBlock.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	if err != nil {
		log.Fatalf("unable to parse PKCS1 private key. %v", err)
	}
	config := &oauth1.Config{
		ConsumerKey: consumerKey,
		CallbackURL: "oob", /* for command line usage */
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: jiraUrl + "plugins/servlet/oauth/request-token",
			AuthorizeURL:    jiraUrl + "plugins/servlet/oauth/authorize",
			AccessTokenURL:  jiraUrl + "plugins/servlet/oauth/access-token",
		},
		Signer: &oauth1.RSASigner{
			PrivateKey: privateKey,
		},
	}
	return &TaskProveder{
		JiraUrl: jiraUrlParsed,
		config:  config,
	}
}

func (j *TaskProveder) storeAuthorizationUrl() string {
	requestToken, requestSecret, err := j.config.RequestToken()
	if err != nil {
		log.Fatalf("Unable to get request token. %v", err)
	}
	authorizationURL, err := j.config.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatalf("Unable to get authorization url. %v", err)
	}

	GlobalStorage.AddAuthorizationUrl(OauthConfiguration{URL: j.JiraUrl.String(),
		AuthorizationURL: authorizationURL.String(),
		RequestSecret: requestSecret,
		RequestToken: requestToken})
	return authorizationURL.String()
}

func (j *TaskProveder) storeAccessToken(code string) {
	oauthConfiguration := GlobalStorage.readOauthConfiguration(j.JiraUrl.String())
	accessToken, accessSecret, err := j.config.AccessToken(oauthConfiguration.RequestToken,
		oauthConfiguration.RequestSecret,
		code)
	if err != nil {
		log.Fatalf("Unable to get access token. %v", err)
	}
	token := oauth1.NewToken(accessToken, accessSecret)
	oauthConfiguration.Token = token.Token
	oauthConfiguration.TokenSecret = token.TokenSecret
	GlobalStorage.StoreToken(oauthConfiguration)
}

func (j *TaskProveder) jiraTokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(j.JiraUrl.Host+".json")), err
}

func jiraTokenFromFile(file string) (*oauth1.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth1.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveJIRAToken(file string, token *oauth1.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (j *TaskProveder) getJIRAClient() *jira.Client {
	ctx := context.Background()
	oauthConfiguration := GlobalStorage.readOauthConfiguration(j.JiraUrl.String())
	if oauthConfiguration.Token == "" {
		return nil;
	}
	jiraClient, err := jira.NewClient(j.config.Client(ctx,
		oauth1.NewToken(oauthConfiguration.Token, oauthConfiguration.TokenSecret)),
		j.JiraUrl.String())
	if err != nil {
		log.Fatalf("unable to create new JIRA client. %v", err)
	}
	return jiraClient
}
