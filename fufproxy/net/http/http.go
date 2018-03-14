//Copyright © 2018 Fuf
//work with http tasks
package http

import (
	"encoding/json"
	"fmt"
	docker "github.com/gitfuf/fuftest/fufproxy/docker"
	"github.com/spf13/viper"
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
	//fmt.Printf("%+v\n", string(theJson))
	//send data back to client
	fmt.Fprint(w, string(theJson))

}

//return port for cli, setup chosen postgres container
func GetProxyPort(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["container"]
	if !ok || len(keys) < 1 {
		fmt.Println("Url Param 'container' is missing")
		return
	}
	docker.SetChosen(keys[0])

	listen_port := viper.GetString("cliport")
	fmt.Fprint(w, listen_port)
}
