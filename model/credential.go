package model

import (
	"time"
)

type OcpiCredentialRole struct {
	OcpiPartyId
	Role            string               `json:"role"`
	BusinessDetails *OcpiBusinessDetails `json:"business_details"`
}

type OcpiVersion struct {
	Version string `json:"version"`
	Url     string `json:"url"`
}

type OcpiVersionModuleEndpoint struct {
	Id   string `json:"identifier"`
	Role string `json:"role"`
	Url  string `json:"url"`
}

type OcpiVersionDetails struct {
	Version   string                      `json:"version"`
	Endpoints []OcpiVersionModuleEndpoint `json:"endpoints"`
}

type OcpiCredentials struct {
	Token string                `json:"token"`
	Url   string                `json:"url"`
	Roles []*OcpiCredentialRole `json:"roles"`
}

type OcpiClientInfo struct {
	OcpiPartyId
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

type OcpiResponseAny struct {
	OcpiResponse
	Data any `json:"data"`
}

type OcpiCredentialsResponse struct {
	OcpiResponse
	Data *OcpiCredentials `json:"data"`
}

type OcpiVersionsResponse struct {
	OcpiResponse
	Data []*OcpiVersion `json:"data"`
}

type OcpiVersionDetailsResponse struct {
	OcpiResponse
	Data *OcpiVersionDetails `json:"data"`
}

type OcpiClientInfoResponse struct {
	OcpiResponse
	Data []*OcpiClientInfo `json:"data"`
}
