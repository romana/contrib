package main

import (
	"flag"
	"fmt"
	"log"
	"encoding/json"

	search "github.com/romana/contrib/rsearch"
)

func main() {
	var cfgFile = flag.String("c", "", "Kubernetes reverse search config file")
	var server = flag.Bool("s", false, "Start a server")
	var host = flag.String("h", "", "Protocol://host for client to connect to")
	var searchTag = flag.String("r", "", "Search resources by tag")
	flag.Parse()

	done := make(chan search.Done)

	config, err := search.NewConfig(*cfgFile)
	if err != nil {
		log.Fatalf("Can not read config file %s, %s\n", *cfgFile, err)
	}

	if *host != "" {
		config.Server.Host = *host
	}

	if *server {
		log.Println("Starting server")
		nsUrl := fmt.Sprintf("%s/%s", config.Api.Url, config.Api.NamespaceUrl)
		nsEvents, err := search.NsWatch(done, nsUrl, config)
		if err != nil {
			log.Fatal("Namespace watcher failed to start", err)
		}

		events := search.Conductor(nsEvents, done, config)
		req := search.Process(events, done, config)
		log.Println("All routines started")
		search.Serve(config, req)
	} else if len(*searchTag) > 0 {
		if config.Server.Debug {
			log.Println("Making request t the server")
		}
		r := search.SearchResource(config, search.SearchRequest{Tag: *searchTag})
		response, err := json.Marshal(r)
		if err != nil {
			log.Fatal("Failed to parse out server response, ", err)
		}
		fmt.Println(string(response))
	} else {
		log.Fatal("Either -s or -r must be given")
	}
}
