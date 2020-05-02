package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/alexandrevilain/atome_exporter/internal/config"
	"github.com/alexandrevilain/atome_exporter/internal/exporter"
	"github.com/alexandrevilain/atome_exporter/pkg/atome"
	"github.com/alexandrevilain/atome_exporter/pkg/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := logrus.New()

	config, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.New(logger, "atome_exporter.db", "atome")
	if err != nil {
		log.Fatal(err)
	}

	atome := atome.NewClient(logger, config.Atome.Username, config.Atome.Password, storage)
	exporter := exporter.New(logger, atome)

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())

	listen := fmt.Sprintf("%s:%d", config.ListenAddr, config.ListenPort)
	logger.Infof("Server running an listening: %s", listen)
	http.ListenAndServe(listen, nil)
}
