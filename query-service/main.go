package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/f3rcho/cqrs/database"
	"github.com/f3rcho/cqrs/events"
	"github.com/f3rcho/cqrs/repository"
	"github.com/f3rcho/cqrs/search"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", listFeedsHandler).Methods(http.MethodGet)
	router.HandleFunc("/search", searchHandler).Methods(http.MethodGet)
	return
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB)
	repo, err := database.NewPostgressRepository(addr)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)

	es, err := search.NewElastic(fmt.Sprintf("http://%s", config.ElasticsearchAddress))
	if err != nil {
		log.Fatal(err)
	}
	search.SetSearchRepository(es)

	n, err := events.NewNats(fmt.Sprintf("nats://%s", config.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}
	err = n.OnCreatedFeed(onCreatedFeed)
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)

	defer events.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
