//Copyright © 2018 Fuf
//work with http tasks
package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Send HTTP GET request and return answer
func SendRequest(req string) ([]byte, error) {
	resp, err := http.Get(req)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
