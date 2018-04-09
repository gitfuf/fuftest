//Copyright © 2018 Fuf fufcli
//work with http tasks
package http

import (
	"io/ioutil"
	"log"
	"net/http"
)

//Send HTTP GET request and return answer
func SendRequest(req string) ([]byte, error) {
	resp, err := http.Get(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
