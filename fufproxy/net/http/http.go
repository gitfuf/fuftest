//Copyright © 2018 Fuf
//work with http tasks
package http

import (
	"encoding/json"
	"fmt"
	docker "github.com/gitfuf/fuftest/fufproxy/docker"
	fuftcp "github.com/gitfuf/fuftest/fufproxy/net/tcp"
	"log"
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

//request for docker container list
func GetDockerLs(w http.ResponseWriter, r *http.Request) {
	list, err := docker.DockerList()
	if err != nil {
		log.Fatal(err)
	}
	theJson, _ := json.Marshal(list)
	log.Printf("%+v\n", string(theJson))
	//send data back to client
	fmt.Fprint(w, string(theJson))

}

//return port for cli in order to connect to the chosen docker postgres
func GetProxyPort(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["container"]
	if !ok || len(keys) < 1 {
		log.Println("Url Param 'container' is missing")
		return
	}
	//f docker.SetChosen(keys[0])

	var listen_port string //:= viper.GetString("cliport")
	for port, ditem := range fuftcp.GpRoutes {
		if ditem.Id == keys[0] || ditem.Name == keys[0] {
			listen_port = port
		}
	}

	fmt.Fprint(w, listen_port)
}
