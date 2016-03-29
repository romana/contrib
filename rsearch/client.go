// Copyright (c) 2016 Pani Networks
// All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package rsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SearchResource makes HTTP request to the server
// instance if this library.
func SearchResource(config Config, req SearchRequest) SearchResponse {
	// TODO need to make url configurable
	url := config.Server.Host + ":" + config.Server.Port
	data := []byte(`{ "tag" : "` + req.Tag + `"}`)
	if config.Server.Debug {
		log.Println("Making request with", string(data))
	}

	// Making request.
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("HTTP request failed", url, err)
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	sr := SearchResponse{}
	if config.Server.Debug {
		log.Println("Trying to decode", response.Body)
	}
	err = decoder.Decode(&sr)
	if err != nil {
		log.Println("Failed to decode", response.Body)
		panic(err)
	}

	if config.Server.Debug {
		fmt.Println("Decoded response form a server", sr)
	}
	return sr
}
