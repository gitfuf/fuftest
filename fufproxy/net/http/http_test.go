//Copyright © 2018 Fuf
//work with http tasks for proxy

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	docker "github.com/gitfuf/fuftest/fufproxy/docker"
	fuftcp "github.com/gitfuf/fuftest/fufproxy/net/tcp"
)

type TestCase struct {
	Container string
	Port      string
	//Response   string
	StatusCode int
}

func createTestGpRoutes() (map[string]docker.DockerItem, error) {
	entry1 := docker.DockerItem{
		Image:    "postgres",
		Id:       "82a7cc1184b4",
		Name:     "pg_fuf",
		Endpoint: "127.0.0.1:3437->5432/tcp",
	}

	list := make(map[int]docker.DockerItem)
	list[0] = entry1
	routes, err := fuftcp.FillinGpRoutes(list)
	return routes, err

}

func createTestDockerList() map[int]docker.DockerItem {
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
	return list
}

func TestCheckPortRequest(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Container: "pg_fuf",
			Port:      "7900",
			//Response:   `{"status":200, "resp": {"port":"7900"}}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Container: "",
			Port:      "",
			//Response:   `{"status":400, "err":"no_port"}`,
			StatusCode: http.StatusBadRequest,
		},
	}
	testGplist, _ := createTestGpRoutes()

	for caseNum, item := range cases {
		testUrl := "http://local.com/api/askproxyport?container=" + item.Container
		req := httptest.NewRequest("GET", testUrl, nil)

		w := httptest.NewRecorder()
		container, _ := checkPortRequest(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: expected %d, got %d", caseNum, item.StatusCode, w.Code)
		}
		port, _ := searchLnPort(container, testGplist)
		if port != item.Port {
			t.Errorf("[%d] wrong Port: expected %s, got %s", item.Port, port)
		}
	}

}

func TestDockerListToJson(t *testing.T) {
	expectedJS := `{"0":{"Image":"postgres","Id":"03bf6d1c4232","Name":"pg_fuf1","Endpoint":"127.0.0.1:3435-\u003e5432/tcp"},"1":{"Image":"postgres","Id":"82a7cc1184b4","Name":"pg_fuf","Endpoint":"127.0.0.1:3437-\u003e5432/tcp"}}`

	jsonS, _ := dockerlistTojson(createTestDockerList())
	if jsonS != expectedJS {
		t.Errorf("Wrong Json: expected %s, got %s ", expectedJS, jsonS)
	}
}
