package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	prefix_docker = "hap-"
)

func getHaList() ([]string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)

		for _, n := range container.Names {
			if strings.HasPrefix(n, prefix_docker) {
				names = append(names, n)
			}
		}
	}

	return names, nil
}

func runHaproxy(name string, haproxycfg string, hostport string, rm bool) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()
	imageName := "haproxy:1.8"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		//	Cmd:   []string{"ls", "-al", "/usr/local/etc/haproxy"},
		ExposedPorts: nat.PortSet{"80/tcp": {}},
	}, &container.HostConfig{
		Binds:        []string{haproxycfg + ":/usr/local/etc/haproxy/haproxy.cfg"},
		PortBindings: nat.PortMap{"80/tcp": {{HostPort: hostport}}},
		AutoRemove:   true,
	}, nil, prefix_docker+name)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	fmt.Println(resp.ID)

	// statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case <-statusCh:
	// }

	// out, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	// if err != nil {
	// 	panic(err)
	// }

	// io.Copy(os.Stdout, out)
	return nil
}
