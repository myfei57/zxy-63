package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"waterplant/internal/console"
	"waterplant/internal/store"
)

func main() {
	storePath := flag.String("store", "waterplant.json", "file path for the persistent store")
	addr := flag.String("addr", ":8080", "listen address for the control console")
	showVersion := flag.Bool("version", false, "print the control console version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(console.Version())
		return
	}

	s, err := store.Open(*storePath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer s.Close()

	if err := console.Seed(s); err != nil {
		log.Fatalf("seed defaults: %v", err)
	}

	server := console.NewServer(s)
	log.Printf("waterplant control console listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
