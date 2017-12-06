package main

import (
	"github.com/bndr/gojenkins"
	"strings"
	"fmt"
)

type CI struct {
	jenkins *gojenkins.Jenkins
}

func (c *CI) connect(url string, username string, password string) {
	c.jenkins = gojenkins.CreateJenkins(nil, url, username, password)
	c.jenkins.Init()
}

func (c *CI) updateBuild(projectName string, branch string) {
	projects := GlobalStorage.GetProjects()
	jobs := GlobalStorage.GetAllGroupJobs(branch)
	var err error
	existingView, err := c.jenkins.GetView(branch)
	if err != nil {
		panic(err)
	}
	var view *gojenkins.View
	if existingView.GetName() == "" {
		view, err = c.jenkins.CreateView(branch, gojenkins.LIST_VIEW)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("View " + existingView.GetName() + " already exists");
		view, err = c.jenkins.GetView(branch)
		if err != nil {
			panic(err)
		}
	}

	for _, project := range projects {
		jobName := project.Name + "-" + branch
		var environment string
		if storedJob, ok := jobs[jobName]; ok {
			fmt.Println(storedJob.Name)
			environment = storedJob.Group
		} else {
			environment = branch
		}
		existingJob, err := c.jenkins.GetJob(jobName)
		if (err != nil && existingJob == nil) {
			if project.Name != projectName {
				environment = project.DefaultBranch
			} else if existingJob == nil {
				environment = branch
			} else if cfg,_ := existingJob.GetConfig(); strings.Contains(cfg, "<name>*/" + project.DefaultBranch + "</name>") {
				environment = branch
			}
			projectJobConfig := strings.Replace(project.JobDefinition, "{{environment}}", environment, -1)
			fmt.Println(projectJobConfig)
			projectJob, err := c.jenkins.CreateJob(projectJobConfig, jobName)
			if err != nil {
				panic(err)
			}
			view.AddJob(projectJob.GetName())
		} else {
			fmt.Println("Job " + jobName + " already exists");
			if cfg,_ := existingJob.GetConfig(); !strings.Contains(cfg, "<name>*/" + branch + "</name>") {
				c.jenkins.DeleteJob(existingJob.GetName())
				projectJobConfig := strings.Replace(project.JobDefinition, "{{environment}}", environment, -1)
				fmt.Println(projectJobConfig)
				projectJob, err := c.jenkins.CreateJob(projectJobConfig, jobName)
				if err != nil {
					panic(err)
				}
				view.AddJob(projectJob.GetName())
			}
		}
	}
}

func (c *CI) removeBuild(projectName string, branch string) {
	var err error
	existingView, err := c.jenkins.GetView(branch)
	if err != nil {
		panic(err)
	}
	if existingView.GetName() == "" {
		return
	} else {
		jobs := existingView.GetJobs()
		for _, job := range jobs {
			c.jenkins.DeleteJob(job.Name)
		}

	}
}