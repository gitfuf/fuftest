//Copyright © 2018 Fuf
//work docker data
package docker

import (
	"encoding/json"
)

type DockerItem struct {
	Image, Id, Name, Endpoint string
}

var dockerlist map[int]DockerItem

func DockerList() map[int]DockerItem {
	return dockerlist
}

func SetDockerList(data []byte) {
	// Have a JSON map, lets turn it into a Go map
	//var data map[int]DockerItem
	dockerlist = make(map[int]DockerItem)
	json.Unmarshal(data, &dockerlist)
}
