package main

import (
	"net/http"

	"github.com/cosmouser/aad-ec/actions"
	"github.com/cosmouser/aad-ec/config"
)

func main() {
	http.HandleFunc("/", actions.IndexHandler)
	http.HandleFunc("/ece/getPlans", actions.APIHandler)
	http.ListenAndServe(config.Port, nil)
}
