//Copyright © 2018 Fuf fufproxy
//for work with tcp tasks
package tcp

import (
	"testing"

	docker "github.com/gitfuf/fuftest/fufproxy/docker"
)

//p.AddRoute(":80", tcpproxy.To("10.0.0.1:8081"))

func TestFillinGpRoutes(t *testing.T) {

	entry1 := docker.DockerItem{
		Image:    "postgres",
		Id:       "03bf6d1c4232",
		Name:     "pg_fuf1",
		Endpoint: "127.0.0.1:3435->5432/tcp",
	}
	entry2 := docker.DockerItem{
		Image:    "postgres",
		Id:       "82a7cc1184b4",
		Name:     "pg_fuf",
		Endpoint: "127.0.0.1:3437->5432/tcp",
	}

	list := make(map[int]docker.DockerItem)
	list[0] = entry1
	list[1] = entry2
	routes, err := FillinGpRoutes(list)
	if err == nil {
		for key, value := range routes {
			if key != "7900" && key != "7901" {
				t.Errorf("port not valid. expect 7900 or 7901, got %s", key)
			}
			if value != entry1 && value != entry2 {
				t.Errorf("docker entry not valid, got %v", value)
			}
		}
	}

}
