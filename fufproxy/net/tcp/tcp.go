//Copyright © 2018 Fuf
//for work with tcp tasks
package tcp

import (
	"fmt"
	docker "github.com/gitfuf/fuftest/fufproxy/docker"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
)

//start tcp server for cliport
func StartServer() {
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
	defer incoming_conn.Close()
	//fmt.Println("handle connection ", incoming_conn.RemoteAddr())

	endpoint := docker.EndPoint()
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
