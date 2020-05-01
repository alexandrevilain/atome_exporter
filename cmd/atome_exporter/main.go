package main

import (
	"fmt"
	"log"

	"github.com/alexandrevilain/atome_exporter/internal/config"
	"github.com/alexandrevilain/atome_exporter/pkg/atome"
)

func main() {
	config, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	atome, err := atome.NewClient(config.Atome.Username, config.Atome.Password)
	if err != nil {
		log.Fatal(err)
	}
	err = atome.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	val, err := atome.RetriveDayConsumption()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
