//Copyright © 2018 Fuf
//work with http tasks
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	docker "github.com/gitfuf/fuftest/fufproxy/docker"
	fuftcp "github.com/gitfuf/fuftest/fufproxy/net/tcp"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func dockerlistTojson(list map[int]docker.DockerItem) (string, error) {
	json, _ := json.Marshal(list)
	log.Printf("%+v\n", string(json))
	return string(json), nil
}

//request for docker container list
func GetDockerLs(w http.ResponseWriter, r *http.Request) {
	list, err := docker.DockerList()
	if err != nil {
		log.Fatal(err)
	}
	jsonS, _ := dockerlistTojson(list)

	//send data back to client
	fmt.Fprint(w, jsonS)

}

func checkPortRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	//keys, ok := r.URL.Query()["container"]
	key := r.FormValue("container")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Url Param 'container' is missing")
		return "", errors.New("Url Param 'container' is missing")
	}
	w.WriteHeader(http.StatusOK)
	return key, nil
}

func searchLnPort(cont string, gpRoutes map[string]docker.DockerItem) (string, error) {
	listen_port := ""
	for port, ditem := range gpRoutes {
		if ditem.Id == cont || ditem.Name == cont {
			listen_port = port
			return listen_port, nil
		}
	}
	return listen_port, errors.New("No suitable listen port")
}

//return port for cli in order to connect to the chosen docker postgres
func GetProxyPort(w http.ResponseWriter, r *http.Request) {
	container, err := checkPortRequest(w, r)
	if err != nil {
		return
	}

	//TODO err
	listen_port, _ := searchLnPort(container, fuftcp.GpRoutes)

	fmt.Fprint(w, listen_port)
}
