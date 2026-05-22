package service

import (
	"fmt"

	pkg "oxio-phone-lookup/internal/pkg"

	"github.com/nyaruka/phonenumbers"
)

const ()

type Service interface {
	LookupPhoneNumber(params pkg.GetV1PhoneNumbersParams) (*pkg.PhoneNumberResponse, error)
}

type service struct{}

func New() Service {
	return &service{}
}

// LookupPhoneNumber takes a phone number and optional country code, parses and validates the phone number,
// and returns the normalized E.164 format along with the country code, area code, and local phone number.
func (s *service) LookupPhoneNumber(params pkg.GetV1PhoneNumbersParams) (*pkg.PhoneNumberResponse, error) {
	countryCode := ""
	if params.CountryCode != nil {
		countryCode = *params.CountryCode
	}
	phoneNumber := params.PhoneNumber

	// if countryCode isn't provided, it will try to parse it if it include a + prefix
	parsed, err := phonenumbers.Parse(phoneNumber, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse phone number '%s' and country code '%s' with error: %w", phoneNumber, countryCode, err)
	}

	if !phonenumbers.IsValidNumber(parsed) {
		return nil, fmt.Errorf("invalid phone number '%s' and country code '%s'", phoneNumber, countryCode)
	}

	regionCode := phonenumbers.GetRegionCodeForNumber(parsed)
	countryDialCode := parsed.GetCountryCode()
	nationalNum := fmt.Sprintf("%d", parsed.GetNationalNumber())

	const localNumerLength = 7
	areaCode := ""
	localNumber := nationalNum

	// area code is everything before the last 7 digits; only accurate for NANP (North American Numbering Plan) numbers
	// this will capture non-NANP numbers
	if len(nationalNum) > localNumerLength {
		areaCode = nationalNum[:len(nationalNum)-localNumerLength]
		localNumber = nationalNum[len(nationalNum)-localNumerLength:]
	}

	e164 := fmt.Sprintf("+%d%s", countryDialCode, nationalNum)

	return &pkg.PhoneNumberResponse{
		PhoneNumber:      e164,
		CountryCode:      regionCode,
		AreaCode:         areaCode,
		LocalPhoneNumber: localNumber,
	}, nil
}
