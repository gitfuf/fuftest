// Copyright © 2018 Fuf fufcli
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

	docker "github.com/gitfuf/fuftest/fufcli/docker"
	httpfuf "github.com/gitfuf/fuftest/fufcli/net/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "use ls to get docker container list",
	Long:  `use fufcli container ls to get docker container list`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("container called :", args)
		if len(args) > 0 && args[0] == "ls" {
			askDockerList()
		} else {
			fmt.Println("use fufcli container ls")
		}

	},
}

func init() {
	rootCmd.AddCommand(containerCmd)
}

func askDockerList() {

	var request = fmt.Sprintf("http://%s/getdockerls", viper.GetString("web"))

	body, err := httpfuf.SendRequest(request)

	if err != nil {
		log.Fatal(err)
	}

	docker.SetDockerList(body)

	// Print what we got with keys
	fmt.Println("Have next docker containers:")
	for _, e := range docker.DockerList() {
		fmt.Printf("image: %+v, id: %+v, name: %+v \n", e.Image, e.Id, e.Name)
		log.Printf("image: %+v, id: %+v, name: %+v \n", e.Image, e.Id, e.Name)
	}

}
