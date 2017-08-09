package main

import (
	"github.com/bndr/gojenkins"
	"fmt"
)

type Jenkins struct {

}

func (j *Jenkins) connect() {
	jenkins := gojenkins.CreateJenkins("http://ci.controlcho.int:8080/", "apostol.terziev", "chunjurlejka10")
	jenkins.Init()
	build, err := jenkins.GetJob("guardian-master")
	if err != nil {
		panic("Job Does Not Exist")
	}
	job, error := build.Copy("new-guardian-branch")
	if error != nil {
		panic("Couldn't copy a job")
	}
	fmt.Println(job.GetConfig())
}