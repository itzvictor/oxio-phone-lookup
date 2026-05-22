package controller

import (
	"encoding/json"
	"net/http"

	"oxio-phone-lookup/internal/pkg"
	"oxio-phone-lookup/internal/service"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	svc service.Service
}

func New(service service.Service) *Controller {
	return &Controller{
		svc: service,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Get("/v1/phone-numbers", c.LookupPhoneNumber)
}

func (c *Controller) LookupPhoneNumber(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	params := pkg.GetV1PhoneNumbersParams{
		PhoneNumber: r.URL.Query().Get("phoneNumber"),
		CountryCode: &countryCode,
	}

	result, err := c.svc.LookupPhoneNumber(params)
	if err != nil {
		// Return the error response with 400 status code for invalid input
		// Along with the error massage
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(pkg.ErrorResponse{
			PhoneNumber: params.PhoneNumber,
			Error:       map[string]string{"message": err.Error()},
		})
		return
	}

	// Return the successful response and 500 error if encoding fails
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
