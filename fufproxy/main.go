//Copyright © 2018 Fuf
package main

import (
	"fmt"

	httpfuf "github.com/gitfuf/fuftest/fufproxy/net/http"
	tcpfuf "github.com/gitfuf/fuftest/fufproxy/net/tcp"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	initConfig()

	go tcpfuf.StartServer()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", httpfuf.IndexPage)
	//router.HandleFunc("/dockerterm", setDockerTeminal)
	router.HandleFunc("/getdockerls", httpfuf.GetDockerLs)
	router.HandleFunc("/askproxyport", httpfuf.GetProxyPort)

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
