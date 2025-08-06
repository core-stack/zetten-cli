package main

import (
	"log"
	"os"

	"github.com/core-stack/zetten-cli/internal/service/mirror"
	"github.com/kardianos/service"
)

func main() {
	svcConfig := &service.Config{
		Name:        "zetten-service",
		DisplayName: "Zetten service",
		Description: "Sync folders and packages.",
	}

	prg := &mirror.ZettenService{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatalf("Error control: %v", err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
