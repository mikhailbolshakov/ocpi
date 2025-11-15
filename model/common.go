package model

import "time"

type OcpiResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Timestamp     string `json:"timestamp,omitempty"`
}

type OcpiRequestHeader struct {
	RequestId       string `json:"requestId,omitempty"`
	CorrelationId   string `json:"correlationId,omitempty"`
	FromPartyId     string `json:"fromPartyId,omitempty"`
	FromCountryCode string `json:"fromCountryCode,omitempty"`
	ToPartyId       string `json:"toPartyId,omitempty"`
	ToCountryCode   string `json:"toCountryCode,omitempty"`
}

type OcpiGetPageRequest struct {
	OcpiRequestHeader
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
	Offset   *int       `json:"offset,omitempty"`
	Limit    *int       `json:"limit,omitempty"`
}

type OcpiDisplayText struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

type OcpiPartyId struct {
	PartyId     string `json:"party_id"`
	CountryCode string `json:"country_code"`
}

type OcpiImage struct {
	Url       string `json:"url"`                 // Url from where the image data can be fetched
	Thumbnail string `json:"thumbnail,omitempty"` // Thumbnail from where a thumbnail of the image can be fetched
	Category  string `json:"category"`            // Category what the image is used for
	Type      string `json:"type"`                // Type like: gif, jpeg, png, svg
	Width     int    `json:"width,omitempty"`     // Width of the full scale image
	Height    int    `json:"height,omitempty"`    // Height of the full scale image
}

type OcpiBusinessDetails struct {
	Name    string     `json:"name"`              // Name of the operator
	Logo    *OcpiImage `json:"log,omitempty"`     // Logo image link to the party’s logo
	Website string     `json:"website,omitempty"` // Website link to the party’s website
	Inn     string     `json:"inn,omitempty"`     // Inn (tax number) of the party (Yandex extension of the protocol. Isn't supported by OCPI)
}

type OcpiRegularHours struct {
	Weekday     int    `json:"weekday"`      // Weekday from Monday (1) till Sunday (7)
	PeriodBegin string `json:"period_begin"` // PeriodBegin in 24h format with leading zeros. Example: "18:15"
	PeriodEnd   string `json:"period_end"`   // PeriodEnd in 24h format with leading zeros. Example: "18:15"
}

type OcpiExceptionalPeriod struct {
	PeriodBegin time.Time `json:"period_begin"` // PeriodBegin begin of the exception. In UTC, time_zone field can be used to convert to local time
	PeriodEnd   time.Time `json:"period_end"`   // PeriodEnd end of the exception. In UTC, time_zone field can be used to convert to local time
}

type OcpiHours struct {
	TwentyFourSeven     bool                     `json:"twentyfourseven"`                // TwentyFourSeven true to represent 24 hours a day and 7 days a week
	RegularHours        []*OcpiRegularHours      `json:"regular_hours,omitempty"`        // RegularHours weekday-based, used if twentyfourseven=false
	ExceptionalOpenings []*OcpiExceptionalPeriod `json:"exceptional_openings,omitempty"` // ExceptionalOpenings for specified calendar dates, time-range based
	ExceptionalClosings []*OcpiExceptionalPeriod `json:"exceptional_closings,omitempty"` // ExceptionalClosings for specified calendar dates, time-range based
}
