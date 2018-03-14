//Copyright © 2018 Fuf
//for work with tcp tasks
package tcp

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
)

var proxyPort string

func SetProxyPort(p string) {
	proxyPort = p
}

func StartServer(port string) error {
	proxyPort = port
	ln, err := net.Listen("tcp", ":"+viper.GetString("port"))
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
	endpoint := viper.GetString("proxy") + ":" + proxyPort //"192.168.99.100:7816"
	//fmt.Println("endpoint", endpoint)
	defer incoming_conn.Close()

	dest_conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Dial err", err)
	}
	//fmt.Println("connected to ", endpoint)
	defer dest_conn.Close()
	pipe(incoming_conn, dest_conn)

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
