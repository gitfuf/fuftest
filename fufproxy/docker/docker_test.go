//Docker test package for proxy

package docker

import (
	"testing"
)

//test fillin docker list
func TestFillList(t *testing.T) {

	/*data := []string{
		"postgres 03bf6d1c4232 pg_fuf1 127.0.0.1:3437->5432/tcp",
		"postgres 82a7cc1184b4 pg_fuf 127.0.0.1:3435->5432/tcp"
	}
	*/

	testData := "postgres 03bf6d1c4232 pg_fuf1 127.0.0.1:3437->5432/tcp\npostgres 82a7cc1184b4 pg_fuf 127.0.0.1:3435->5432/tcp"
	testDockerlist := fillList(testData)

	if testDockerlist[0].Image != "postgres" {
		t.Errorf("expect image = postgres, got %s", testDockerlist[0].Image)
	}

	if testDockerlist[0].Id != "03bf6d1c4232" {
		t.Errorf("expect id = 03bf6d1c4232, got %s", testDockerlist[0].Id)
	}

	if testDockerlist[0].Name != "pg_fuf1" {
		t.Errorf("expect image = pg_fuf1, got %s", testDockerlist[0].Name)
	}

	if testDockerlist[0].Endpoint != "127.0.0.1:3437->5432/tcp" {
		t.Errorf("expect endpoint = 127.0.0.1:3437->5432/tcp, got %s", testDockerlist[0].Name)
	}

	if testDockerlist[1].Image != "postgres" {
		t.Errorf("expect image = postgres, got %s", testDockerlist[0].Image)
	}

	if testDockerlist[1].Id != "82a7cc1184b4" {
		t.Errorf("expect id = 03bf6d1c4232, got %s", testDockerlist[0].Id)
	}

	if testDockerlist[1].Name != "pg_fuf" {
		t.Errorf("expect image = pg_fuf, got %s", testDockerlist[0].Name)
	}

	if testDockerlist[1].Endpoint != "127.0.0.1:3435->5432/tcp" {
		t.Errorf("expect endpoint = 127.0.0.1:3435->5432/tcp, got %s", testDockerlist[0].Name)
	}

}

//Test correct short endpoint
func TestShortEndpoint(t *testing.T) {

	entry := DockerItem{
		Image:    "postgres",
		Id:       "03bf6d1c4232",
		Name:     "pg_fuf1",
		Endpoint: "127.0.0.1:3437->5432/tcp"}

	short := entry.ShortEndpoint()
	if short != "127.0.0.1:3437" {
		t.Errorf("expect short endpoint = 127.0.0.1:3437, got %s", short)
	}
}
