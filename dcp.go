package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"os/exec"
)

func fail(e error,msg string) {
	if e!=nil {
		log.Fatal(msg,"\n\t",e)
	}
}

func warn(e error,msg string) {
	if e!=nil {
		log.Println(msg,"\n\t-->",e)
	}
}

func main() {
	if len(os.Args)!=3 {
		log.Fatal("Usage: dcp {currentDockerImage} {NewDockerImage}")
	}

	FROM:=os.Args[1]
	TO:=os.Args[2]
	var err error
	ctx:=context.Background()
	fail(err,"Unable to create sync context")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	fail(err,"Failed to set up Docker client")
	r, err := cli.ImagePull(ctx, FROM, types.ImagePullOptions{})
	io.Copy(os.Stdout, r)
	warn(err,"Failed to pull image")
	err=cli.ImageTag(ctx,FROM,TO)
	fail(err,"Tagging failed")
	//r,err=cli.ImagePush(ctx,TO,types.ImagePushOptions{authBytes})
	//fail(err,"Failed to push image")
	//io.Copy(os.Stdout, r)
	RunCmd("docker",[]string{"push",TO})
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	fail(err,"Failed to list containers")
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	is,e:=cli.ImageList(ctx, types.ImageListOptions{})
	fail(e,"Failed to list images")
	for _,img:=range is {
		fmt.Println(img)
	}
}


func RunCmd(in string, args []string) error {
	cmd := exec.Command(in, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// SecurityTokenScript -u ndb338  -p xxxxxx
	//out, err := cmd.CombinedOutput()
	err := cmd.Run()
	return err
}

