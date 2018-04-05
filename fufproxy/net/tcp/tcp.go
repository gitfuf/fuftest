//Copyright © 2018 Fuf fufproxy
//for work with tcp tasks
package tcp

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	docker "github.com/gitfuf/fuftest/fufproxy/docker"
	gproxy "github.com/google/tcpproxy"
)

var (
	gp       gproxy.Proxy
	GpRoutes map[string]docker.DockerItem
)

func FillinGpRoutes() error {
	list, err := docker.DockerList()
	if err != nil {
		log.Println("No docker containers")
		return errors.New("No available routes")
	}

	var port string
	GpRoutes = make(map[string]docker.DockerItem)

	str := ""
	for i, pgsql := range list {
		if i < 10 {
			str = "0"
		}
		port = "79" + str + strconv.Itoa(i)
		gp.AddRoute(":"+port, gproxy.To(pgsql.ShortEndpoint()))
		log.Printf("Add route: %s <-> %s ", port, pgsql.ShortEndpoint())
		GpRoutes[port] = pgsql
	}
	return nil
}

func StartGProxy() {
	err := FillinGpRoutes()
	if err != nil {
		log.Fatalln("Google proxy wasn't run", err)
	}

	//setup signal for init close operation
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	//Go routine to catch signal and call close
	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		err := gp.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("successfull close google tcpproxy")
		os.Exit(1)
	}()

	log.Println("Startup google tcpproxy")
	log.Fatal(gp.Run())
}
