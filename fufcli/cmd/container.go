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

	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

type DockerItem struct {
	Image, Id, Name, Endpoint string
	//State           bool
}

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "use ls to get docker container list",
	Long:  `use fufcli container ls to get docker container list`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("container called :", args)
		if len(args) > 0 && args[0] == "ls" {
			getDockerList()
		} else {
			fmt.Println("use fufcli container ls")
		}

	},
}

func init() {
	rootCmd.AddCommand(containerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//containerCmd.PersistentFlags().String("ls", "", "container list")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// containerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//
}
func getDockerList() {
	var request = fmt.Sprintf("http://%s/getdockerls", viper.GetString("web"))

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

	// Have a JSON map, lets turn it into a Go map
	var data map[int]DockerItem
	data = make(map[int]DockerItem)
	json.Unmarshal(body, &data)

	// Print what we got with keys
	fmt.Println("Have next docker containers:")
	for _, e := range data {
		fmt.Printf("image: %+v, id: %+v, name: %+v \n", e.Image, e.Id, e.Name)
	}

}
