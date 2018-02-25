package main

import (
	"log"
	"net/http"

	"github.com/to6ka/golang-test/service"
)

func main() {
	http.Handle("/", service.Service{}.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
