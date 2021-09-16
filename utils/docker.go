package utils

import (
	"context"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/juju/errors"
)

func ListContainers(cli *client.Client, all bool) ([]types.Container, error) {
	containers, err := cli.ContainerList(context.TODO(), types.ContainerListOptions{All: all})
	if err != nil {
		err = errors.Trace(err)
	}
	return containers, err
}

func ExecInContainer(cli *client.Client, containerID string, command []string) (err error) {
	config := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}
	cmd, err := cli.ContainerExecCreate(context.Background(), containerID, config)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	hijacked, err := cli.ContainerExecAttach(context.Background(), cmd.ID, types.ExecStartCheck{})
	if err != nil {
		err = errors.Trace(err)
		return
	}
	defer hijacked.Close()

	b, err := ioutil.ReadAll(hijacked.Reader)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	resp, err := cli.ContainerExecInspect(context.Background(), cmd.ID)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	if resp.ExitCode != 0 {
		err = errors.Errorf("Exec command failed: %+v, exit code: %d, output: %s", command, resp.ExitCode, string(b))
	}
	return
}
