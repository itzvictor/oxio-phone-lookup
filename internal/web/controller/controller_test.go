package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"oxio-phone-lookup/internal/pkg"
	"oxio-phone-lookup/internal/service/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type ControllerSuite struct {
	suite.Suite
	mockSvc *mocks.Service
	router  chi.Router
}

func (s *ControllerSuite) SetupTest() {
	s.mockSvc = mocks.NewService(s.T())
	ctrl := New(s.mockSvc)
	r := chi.NewRouter()
	ctrl.RegisterRoutes(r)
	s.router = r
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}

func (s *ControllerSuite) TestLookupPhoneNumber_HappyPath() {
	s.mockSvc.On("LookupPhoneNumber", pkg.GetV1PhoneNumbersParams{
		PhoneNumber: "+12125690123",
		CountryCode: new(""),
	}).Return(&pkg.PhoneNumberResponse{
		PhoneNumber:      "+12125690123",
		CountryCode:      "US",
		AreaCode:         "212",
		LocalPhoneNumber: "5690123",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/phone-numbers?phoneNumber=%2B12125690123", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	var resp pkg.PhoneNumberResponse
	s.Require().NoError(json.NewDecoder(w.Body).Decode(&resp))
	s.Equal("+12125690123", resp.PhoneNumber)
	s.Equal("US", resp.CountryCode)
	s.Equal("212", resp.AreaCode)
	s.Equal("5690123", resp.LocalPhoneNumber)
}

func (s *ControllerSuite) TestLookupPhoneNumber_ServiceError() {
	s.mockSvc.On("LookupPhoneNumber", pkg.GetV1PhoneNumbersParams{
		PhoneNumber: "invalid",
		CountryCode: new(""),
	}).Return(nil, errors.New("invalid phone number"))

	req := httptest.NewRequest(http.MethodGet, "/v1/phone-numbers?phoneNumber=invalid", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)

	var resp pkg.ErrorResponse
	s.Require().NoError(json.NewDecoder(w.Body).Decode(&resp))
	s.Equal("invalid", resp.PhoneNumber)
	s.Equal("invalid phone number", resp.Error["message"])
}
