//Copyright © 2018 Fuf
//Package for work with docker container
package docker

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

//neccesary docker container fields
type DockerItem struct {
	Image, Id, Name, Endpoint string
	//State           bool
}

var dockerlist map[int]DockerItem

var chosenContainer string

//Update or create new list of docker running containers
func updateList() error {
	out, err := exec.Command("docker", "ps", "--filter", "status=running", "--format", "{{.Image}} {{.ID}} {{.Names}} {{.Ports}}").Output()
	if err != nil {
		fmt.Errorf("error running docker ps: %v", err)
		log.Fatal(err)
	}

	//fill in dockerlist
	data := string(out)
	if len(data) != 0 {
		//fmt.Println("docker list :", data)
		dockerlist = fillList(data)
	} else {
		return errors.New("No available routes")
	}

	return nil
}

func fillList(str string) map[int]DockerItem {
	dockerlist = make(map[int]DockerItem)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		arr := strings.SplitN(line, " ", 4)
		if len(arr) < 4 {
			continue
		}
		dockerlist[i] = DockerItem{
			arr[0], arr[1], arr[2], arr[3],
		}
		log.Println(dockerlist[i])
	}
	return dockerlist
}

func DockerList() (map[int]DockerItem, error) {
	err := updateList()
	return dockerlist, err
}

/*
func SetChosen(name string) {
	chosenContainer = name
	log.Println("chosen container = ", chosenContainer)
}
*/
//return endpoint for chosen database container
/*func endPoint() string {
	var endpoint string
	for _, docker := range dockerlist {
		fmt.Println("docker ", docker.Image, docker.Id, docker.Name)
		if docker.Image == chosenContainer || docker.Id == chosenContainer || docker.Name == chosenContainer {

			str := docker.Endpoint
			endpoint = str[0:strings.Index(str, "->")]
			return endpoint
		}
	}
	return ""
}
*/
func (d DockerItem) ShortEndpoint() string {
	endpoint := d.Endpoint[0:strings.Index(d.Endpoint, "->")]
	return endpoint
}

/*
"context"
"github.com/docker/docker/api/types"
"github.com/docker/docker/client"

func AltFormDockerList() {

	cli, err := client.NewClientWithOpts(client.WithVersion("1.35"))
	if err != nil {
		panic(err)

	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:12], container.Image)
	}
}

*/
