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
	httpfuf "github.com/gitfuf/fuftest/fufcli/net/http"
	tcpfuf "github.com/gitfuf/fuftest/fufcli/net/tcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "get port to connect postgesql",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called: ", args)
		if len(args) > 0 { //todo check also inside dockeritem map
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
	body, err := httpfuf.SendRequest(request)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error2")
	}
	fmt.Printf("You can use in another terminal psql -h localhost - p %s etc \n", viper.GetString("port"))
	port := fmt.Sprintf("%s", body)

	tcpfuf.StartServer(port)

}
