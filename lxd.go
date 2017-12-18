package main

import (
	"github.com/lxc/lxd/client"
	"fmt"
)

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

func (s *LXDServer) Init() {
	connectionArgs := lxd.ConnectionArgs{
		TLSClientCert:      s.Cert,
		TLSClientKey:       s.Key,
		TLSServerCert:      s.ServerCert,
	}
	containerServer, err := lxd.ConnectLXD(s.Url, &connectionArgs)
	if err != nil {
		fmt.Println("Cannot connect to lxd server " + s.Url)
	}
	snapshot, etag, err := containerServer.GetContainerSnapshot("app1", "lxd-base")
	fmt.Println(snapshot.Name + "\t" + etag)
	res, err := containerServer.GetServerResources();
	fmt.Println(res.Memory.Total/1024/1024)
	op, err := containerServer.CopyContainerSnapshot(containerServer, *snapshot, &lxd.ContainerSnapshotCopyArgs{Name: "gugi"})
	if err != nil {
		panic(err)
	}
	op.Wait()
}
