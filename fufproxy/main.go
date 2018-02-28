package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

var chosenDBcontainer string

func main() {
	initConfig()

	go setupServer()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/dockerterm", setDockerTeminal)
	router.HandleFunc("/getdockerls", getDockerLs)
	router.HandleFunc("/askproxyport", getProxyPort)

	log.Fatal(http.ListenAndServe(":8888", router))

}

func initConfig() {
	fmt.Println("initConfig")

	viper.AddConfigPath(".")
	viper.SetConfigName("fufproxy")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

//start tcp server for cliport
func setupServer() {
	fmt.Println("Listen ", viper.GetString("cliport"))
	ln, err := net.Listen("tcp", ":"+viper.GetString("cliport")) // "7816"
	if err != nil {
		log.Fatal(err)
		fmt.Println("listen err", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			fmt.Println("listen err", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(incoming_conn net.Conn) {
	fmt.Println("PROXY Message Received")
	defer incoming_conn.Close()

	//get postgre endpoint from dockerlist
	if len(dockerlist) == 0 {
		formDockerList()
	}

	var endpoint string
	for _, docker := range dockerlist {
		if docker.Image == chosenDBcontainer || docker.Id == chosenDBcontainer || docker.Name == chosenDBcontainer {

			str := docker.Endpoint
			endpoint = str[0:strings.Index(str, "->")]
			fmt.Println("endpoint : ", endpoint)
		}
		fmt.Println("docker ", docker.Image, docker.Id, docker.Name)
	}

	//endpoint := "127.0.0.1:5434"
	fmt.Println("endpoint : ", endpoint)
	dest_conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Dial err", err)
	}
	fmt.Println("connected to ", endpoint)
	defer dest_conn.Close()
	err = pipe(incoming_conn, dest_conn)
	if err != nil {
		fmt.Println("pipe err", err)
	}
}

//for proxy forward
func pipe(a, b net.Conn) error {
	errors := make(chan error, 1)
	copy := func(write, read net.Conn) {
		_, err := io.Copy(write, read)
		errors <- err
	}
	go copy(a, b)
	go copy(b, a)
	return <-errors
}
