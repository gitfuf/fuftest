//Copyright © 2018 Fuf fufcli
//for work with tcp tasks
package tcp

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
)

var (
	proxyPort string
	closeC    chan bool //chanel to notify close connection
)

func SetProxyPort(p string) {
	proxyPort = p
}

func StartServer(port string) error {
	proxyPort = port
	closeC = make(chan bool) //in order to stop server

	ln, err := net.Listen("tcp", ":"+viper.GetString("port")) //:7815
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)

		<-closeC
		log.Println("connection with postgresql closed -> stop server")
		conn.Close()
		return nil
	}

}

func handleConnection(incoming_conn net.Conn) {
	endpoint := viper.GetString("proxy") + ":" + proxyPort //"192.168.99.100:7816"
	log.Println("handleconnection proxy endpoint = ", endpoint)
	defer incoming_conn.Close()

	dest_conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer dest_conn.Close()
	pipe(incoming_conn, dest_conn)
	closeC <- true

}

//for proxy transfer
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
