package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type DockerItem struct {
	Image, Id, Name, Endpoint string
	//State           bool
}

var dockerlist map[int]DockerItem

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

//request for docker container list
func getDockerLs(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetDockerLs request")
	formDockerList()
	theJson, _ := json.Marshal(dockerlist)
	//fmt.Printf("%+v\n", string(theJson))
	//send data back to client
	fmt.Fprint(w, string(theJson))

}
func formDockerList() {
	out, err := exec.Command("docker", "ps", "--filter", "status=running", "--format", "{{.Image}} {{.ID}} {{.Names}} {{.Ports}}").Output()
	if err != nil {
		fmt.Errorf("error running docker ps: %v", err)
		log.Fatal(err)

	}

	// have to fill dockerlist
	dockerlist = make(map[int]DockerItem)
	lines := strings.Split(string(out), "\n")

	for i, line := range lines {
		arr := strings.SplitN(line, " ", 4)
		if len(arr) < 4 {
			continue
		}
		dockerlist[i] = DockerItem{
			arr[0], arr[1], arr[2], arr[3],
		}
		//fmt.Println(dockerlist[i])
	}

}

//return port for cli, setup chosen postgres container
func getProxyPort(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["container"]
	if !ok || len(keys) < 1 {
		fmt.Println("Url Param 'container' is missing")
		return
	}
	chosenDBcontainer = keys[0]

	listen_port := viper.GetString("cliport")
	//fmt.Println("db container : ", chosenDBcontainer)
	fmt.Fprint(w, listen_port)
}

//----------- for web client --------------------------------
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//TODO a lot ;)
func setDockerTeminal(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)

	}
	defer ws.Close()
	containerId := "pg_fuf"
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		//run $docker exec pg_fuf ls -la
		arg := append(
			[]string{
				fmt.Sprintf("exec %s %s", containerId, message),
			})
		cmdname := "docker"
		cmd := exec.Command(cmdname, arg...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("docker exec res: %q\n", out.String())
		err = ws.WriteMessage(mt, out.Bytes())
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}
