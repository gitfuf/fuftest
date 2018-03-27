//Copyright © 2018 Fuf fufproxy
package main

import (
	httpfuf "github.com/gitfuf/fuftest/fufproxy/net/http"
	tcpfuf "github.com/gitfuf/fuftest/fufproxy/net/tcp"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func main() {
	initConfig()

	go tcpfuf.StartGProxy()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", httpfuf.IndexPage)
	router.HandleFunc("/getdockerls", httpfuf.GetDockerLs)
	router.HandleFunc("/askproxyport", httpfuf.GetProxyPort)

	log.Fatal(http.ListenAndServe(":8888", router))

}

func initConfig() {

	file, err := os.OpenFile("fufproxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file :", err)
	}

	//mylog := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	log.SetOutput(file)

	viper.AddConfigPath(".")
	viper.SetConfigName("fufproxy")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

}
