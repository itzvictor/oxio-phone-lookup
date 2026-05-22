package main

import (
	"log"
	"net/http"

	"oxio-phone-lookup/web/controller"
	"oxio-phone-lookup/internal/service"

	"github.com/go-chi/chi/v5"
)

func main() {
	svc := service.New()
	ctrl := controller.New(svc)

	r := chi.NewRouter()
	ctrl.RegisterRoutes(r)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
