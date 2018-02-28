// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "get port to connect postgesql",
	Long:  `get port to connect postgesql`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called: ", args)
		if len(args) > 0 { //todo check also inside dockeritem map
			//getDockerList()
			askProxyForConnection(args[0])
		} else {
			fmt.Println("use fufcli proxy <container name>")
		}

	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

}

//ask connection to chosen postgre container
func askProxyForConnection(container string) {
	var request = fmt.Sprintf("http://%s/askproxyport?container=%s", viper.GetString("web"), container)

	resp, err := http.Get(request)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error2")
	}
	fmt.Printf("You can use in another terminal psql -h localhost - p %s etc \n", viper.GetString("port"))
	port := fmt.Sprintf("%s", body)
	setupServer(port)

}

func setupServer(port_to_proxy string) string {
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
		go handleConnection(conn, port_to_proxy)
	}
}

func handleConnection(incoming_conn net.Conn, port string) {
	endpoint := viper.GetString("proxy") + ":" + port //"192.168.99.100:7816"
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
