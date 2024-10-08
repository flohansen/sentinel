package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/flohansen/sentinel/internal/cli"
	"github.com/flohansen/sentinel/internal/proxy"
)

var (
	configFile = flag.String("config", ".sentinel.json", "The path to the config file")
)

func main() {
	flag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			info, ok := debug.ReadBuildInfo()
			if !ok {
				log.Fatal("could not read version")
			}

			fmt.Printf("%s\n", info.Main.Version)
		case "init":
			if err := cli.Init(); err != nil {
				log.Fatalf("init error: %s", err)
			}
		default:
			log.Fatalf("unknown argument '%s'", os.Args[1])
		}

		return
	}

	cfg := proxy.NewConfig(proxy.FromFile(*configFile))
	proxy := proxy.NewProxy(proxy.WithClient(&http.Client{}), proxy.WithConfig(cfg))
	log.Fatalf("error: %s", proxy.Start())
}
