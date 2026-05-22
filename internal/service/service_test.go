package service

import (
	"testing"

	"oxio-phone-lookup/internal/pkg"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	svc Service
}

func (s *ServiceSuite) SetupTest() {
	s.svc = New()
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestLookupPhoneNumber() {
	tests := []struct {
		name        string
		phoneNumber string
		countryCode *string
		wantPhone   string
		wantCountry string
		wantArea    string
		wantLocal   string
		wantErr     *string
	}{
		{
			name:        "E.164 with plus",
			phoneNumber: "+12125690123",
			wantPhone:   "+12125690123",
			wantCountry: "US",
			wantArea:    "212",
			wantLocal:   "5690123",
		},
		{
			name:        "E.164 with spaces",
			phoneNumber: "+52 631 3118150",
			wantPhone:   "+526313118150",
			wantCountry: "MX",
			wantArea:    "631",
			wantLocal:   "3118150",
		},
		{
			name:        "no plus with spaces",
			phoneNumber: "34 915 872200",
			countryCode: new("ES"),
			wantPhone:   "+34915872200",
			wantCountry: "ES",
			wantArea:    "91",
			wantLocal:   "5872200",
		},
		{
			name:        "local number with country code param",
			phoneNumber: "2125690123",
			countryCode: new("US"),
			wantPhone:   "+12125690123",
			wantCountry: "US",
			wantArea:    "212",
			wantLocal:   "5690123",
		},
		{
			name:        "local number missing country code",
			phoneNumber: "6313118150",
			wantErr:     new("failed to parse phone number '6313118150' and country code '' with error: invalid country code"),
		},
		{
			name:        "invalid characters",
			phoneNumber: "123-456-7890",
			wantErr:     new("failed to parse phone number '123-456-7890' and country code '' with error: invalid country code"),
		},
		{
			name:        "empty phone number",
			phoneNumber: "",
			wantErr:     new("failed to parse phone number '' and country code '' with error: the phone number supplied is not a number"),
		},
		{
			name:        "invalid country code",
			phoneNumber: "2125690123",
			countryCode: new("ESP"),
			wantErr:     new("failed to parse phone number '2125690123' and country code 'ESP' with error: invalid country code"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			params := pkg.GetV1PhoneNumbersParams{
				PhoneNumber: tt.phoneNumber,
				CountryCode: tt.countryCode,
			}
			result, err := s.svc.LookupPhoneNumber(params)
			if tt.wantErr != nil {
				s.EqualError(err, *tt.wantErr)
				s.Nil(result)
				return
			}
			s.Require().NoError(err)
			s.Equal(tt.wantPhone, result.PhoneNumber)
			s.Equal(tt.wantCountry, result.CountryCode)
			s.Equal(tt.wantArea, result.AreaCode)
			s.Equal(tt.wantLocal, result.LocalPhoneNumber)
		})
	}
}
