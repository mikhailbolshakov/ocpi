package ocpi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"net/http"
	"time"
)

type mockClientImpl struct {
	logService domain.OcpiLogService
}

func newMockOcpiRestClient(logService domain.OcpiLogService) ocpiRestClient {
	return &mockClientImpl{
		logService: logService,
	}
}

func (s *mockClientImpl) l() kit.CLogger {
	return service.L().Cmp("mock-ocpi-rest")
}

func (s *mockClientImpl) Init(ctx context.Context, config *service.CfgOcpiRemote) error {
	s.l().Mth("init").Dbg()
	return nil
}

func (s *mockClientImpl) Close(ctx context.Context) error {
	s.l().Mth("close").Dbg()
	return nil
}

var (
	okResp = model.OcpiResponse{
		StatusCode:    model.OcpiStatusCodeOk,
		StatusMessage: model.OcpiStatusMessageSuccess,
		Timestamp:     kit.Now().Format(time.RFC3339),
	}
)

func (s *mockClientImpl) cred() *model.OcpiCredentials {
	return &model.OcpiCredentials{
		Token: kit.NewRandString(),
		Url:   "http://test.io/ocpi/versions",
		Roles: []*model.OcpiCredentialRole{
			{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				Role: model.OcpiRoleCPO,
				BusinessDetails: &model.OcpiBusinessDetails{
					Name: "US Test Limited",
					Logo: &model.OcpiImage{
						Url:       "https://us-test-limited.com/logo",
						Thumbnail: "https://us-test-limited.com/logo/thumbnail",
						Type:      "jpeg",
						Width:     512,
						Height:    512,
					},
					Website: "https://us-test-limited.com",
				},
			},
			{
				Role: model.OcpiRoleEMSP,
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     kit.NewRandString(),
					CountryCode: "US",
				},
			},
		},
	}
}

func (s *mockClientImpl) GetVersions(ctx context.Context, url, token, fromPlatform, toPlatform string) ([]*model.OcpiVersion, error) {
	s.l().Mth("get-version").Dbg()
	rs := &model.OcpiVersionsResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiVersion{
			{
				Version: "2.2.1",
				Url:     "http://test.io/ocpi/2.2.1/",
			},
			{
				Version: "2.2",
				Url:     "http://test.io/ocpi/2.2/",
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetVersions, url, token, fromPlatform, toPlatform, nil, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) GetVersionDetails(ctx context.Context, url, token, fromPlatform, toPlatform string) (*model.OcpiVersionDetails, error) {
	s.l().Mth("get-version-det").Dbg()
	rs := &model.OcpiVersionDetailsResponse{
		OcpiResponse: okResp,
		Data: &model.OcpiVersionDetails{
			Version: "2.2.1",
		},
	}
	for moduleId, d := range model.OcpiModules {
		if d.SenderOnly {
			rs.Data.Endpoints = append(rs.Data.Endpoints, model.OcpiVersionModuleEndpoint{
				Id:   moduleId,
				Role: model.OcpiSender,
				Url:  fmt.Sprintf("http://test.io/ocpi/2.2.1/%s", moduleId),
			})
		} else {
			rs.Data.Endpoints = append(rs.Data.Endpoints, model.OcpiVersionModuleEndpoint{
				Id:   moduleId,
				Role: model.OcpiSender,
				Url:  fmt.Sprintf("http://test.io/ocpi/sender/2.2.1/%s", moduleId),
			})
			rs.Data.Endpoints = append(rs.Data.Endpoints, model.OcpiVersionModuleEndpoint{
				Id:   moduleId,
				Role: model.OcpiReceiver,
				Url:  fmt.Sprintf("http://test.io/ocpi/receiver/%s", moduleId),
			})
		}
	}
	s.makeRequest(ctx, domain.LogEventGetVersions, url, token, fromPlatform, toPlatform, nil, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) PostCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiCredentials) (*model.OcpiCredentials, error) {
	s.l().Mth("post-cred").Dbg()
	rs := &model.OcpiCredentialsResponse{
		OcpiResponse: okResp,
		Data:         s.cred(),
	}
	s.makeRequest(ctx, domain.LogEventGetVersions, url, token, fromPlatform, toPlatform, pl, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) PutCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiCredentials) (*model.OcpiCredentials, error) {
	s.l().Mth("put-cred").Dbg()
	return s.PostCredentials(ctx, url, token, fromPlatform, toPlatform, pl)
}

func (s *mockClientImpl) GetCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string) (*model.OcpiCredentials, error) {
	s.l().Mth("get-cred").Dbg()
	rs := &model.OcpiCredentialsResponse{
		OcpiResponse: okResp,
		Data:         s.cred(),
	}
	s.makeRequest(ctx, domain.LogEventGetCredentials, url, token, fromPlatform, toPlatform, nil, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) DeleteCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string) error {
	s.makeRequest(ctx, domain.LogEventDelCredentials, url, token, fromPlatform, toPlatform, nil, nil)
	return nil
}

func (s *mockClientImpl) GetHubClientInfo(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiClientInfo, error) {
	s.l().Mth("get-client-info").Dbg()
	rs := &model.OcpiClientInfoResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiClientInfo{
			{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				Role:        model.OcpiRoleCPO,
				Status:      model.OcpiStatusConnected,
				LastUpdated: kit.Now(),
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetVersions, url, token, fromPlatform, toPlatform, rq, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) PutClientInfo(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiClientInfo) error {
	s.makeRequest(ctx, domain.LogEventPutClientInfo, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) PutLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiLocation) error {
	s.makeRequest(ctx, domain.LogEventPutLocation, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) PatchLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiLocation) error {
	s.makeRequest(ctx, domain.LogEventPatchLocation, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) GetLocationPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiLocation, error) {
	s.l().Mth("get-loc-page").Dbg()
	rs := &model.OcpiLocationsResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiLocation{
			{
				Id: kit.NewRandString(),
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				LastUpdated: kit.Now(),
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetLocations, url, token, fromPlatform, toPlatform, rq, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) GetLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, locId string) (*model.OcpiLocation, error) {
	s.l().Mth("get-loc").Dbg()
	rs := &model.OcpiLocation{
		Id: kit.NewRandString(),
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     "TEST",
			CountryCode: "US",
		},
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetLocation, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) PutEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiEvse, party *model.OcpiPartyId, locId string) error {
	s.makeRequest(ctx, domain.LogEventPutEvse, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) PatchEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiEvse, party *model.OcpiPartyId, locId string) error {
	s.makeRequest(ctx, domain.LogEventPatchEvse, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) GetEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, locId, evseId string) (*model.OcpiEvse, error) {
	s.l().Mth("get-evse").Dbg()
	rs := &model.OcpiEvse{
		Uid:         kit.NewRandString(),
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetLocation, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) PutCon(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiConnector, party *model.OcpiPartyId, locId, evseId string) error {
	s.makeRequest(ctx, domain.LogEventPutEvse, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) PatchCon(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiConnector, party *model.OcpiPartyId, locId, evseId string) error {
	s.makeRequest(ctx, domain.LogEventPatchEvse, url, token, fromPlatform, toPlatform, rq, nil)
	return nil
}

func (s *mockClientImpl) GetCon(ctx context.Context, url, token, fromPlatform, toPlatform, locId, evseId, conId string) (*model.OcpiConnector, error) {
	s.l().Mth("get-con").Dbg()
	rs := &model.OcpiConnector{
		Id:          kit.NewRandString(),
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetLocation, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) PutTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trf *model.OcpiTariff) error {
	s.makeRequest(ctx, domain.LogEventPutTariff, url, token, fromPlatform, toPlatform, trf, nil)
	return nil
}

func (s *mockClientImpl) PatchTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trf *model.OcpiTariff) error {
	s.makeRequest(ctx, domain.LogEventPatchTariff, url, token, fromPlatform, toPlatform, trf, nil)
	return nil
}

func (s *mockClientImpl) GetTariffPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiTariff, error) {
	s.l().Mth("get-trf-page").Dbg()
	rs := &model.OcpiTariffsResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiTariff{
			{
				Id: kit.NewRandString(),
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				LastUpdated: kit.Now(),
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetTariffs, url, token, fromPlatform, toPlatform, rq, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) GetTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trfId string) (*model.OcpiTariff, error) {
	s.l().Mth("get-trf").Dbg()
	rs := &model.OcpiTariff{
		Id: kit.NewRandString(),
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     "TEST",
			CountryCode: "US",
		},
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetTariff, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) PutToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tkn *model.OcpiToken) error {
	s.makeRequest(ctx, domain.LogEventPutToken, url, token, fromPlatform, toPlatform, tkn, nil)
	return nil
}

func (s *mockClientImpl) PatchToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tkn *model.OcpiToken) error {
	s.makeRequest(ctx, domain.LogEventPatchToken, url, token, fromPlatform, toPlatform, tkn, nil)
	return nil
}

func (s *mockClientImpl) GetTokenPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiToken, error) {
	s.l().Mth("get-tkn-page").Dbg()
	rs := &model.OcpiTokensResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiToken{
			{
				Id: kit.NewRandString(),
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				LastUpdated: kit.Now(),
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetTokens, url, token, fromPlatform, toPlatform, rq, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) GetToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tknId string) (*model.OcpiToken, error) {
	s.l().Mth("get-tkn").Dbg()
	rs := &model.OcpiToken{
		Id: kit.NewRandString(),
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     "TEST",
			CountryCode: "US",
		},
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetToken, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) PutSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sess *model.OcpiSession) error {
	s.makeRequest(ctx, domain.LogEventPutSession, url, token, fromPlatform, toPlatform, sess, nil)
	return nil
}

func (s *mockClientImpl) PatchSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sess *model.OcpiSession) error {
	s.makeRequest(ctx, domain.LogEventPatchSession, url, token, fromPlatform, toPlatform, sess, nil)
	return nil
}

func (s *mockClientImpl) GetSessionPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiSession, error) {
	s.l().Mth("get-sess-page").Dbg()
	rs := &model.OcpiSessionsResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiSession{
			{
				Id: kit.NewRandString(),
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				LastUpdated: kit.Now(),
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetSessions, url, token, fromPlatform, toPlatform, rq, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) GetSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sessId string) (*model.OcpiSession, error) {
	s.l().Mth("get-sess").Dbg()
	rs := &model.OcpiSession{
		Id: kit.NewRandString(),
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     "TEST",
			CountryCode: "US",
		},
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetSession, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) PostCommand(ctx context.Context, url, token, fromPlatform, toPlatform, cmdType string, cmd any) error {
	s.makeRequest(ctx, domain.LogEventPostCommand, url, token, fromPlatform, toPlatform, cmd, nil)
	return nil
}

func (s *mockClientImpl) PostCommandResponse(ctx context.Context, url, token, fromPlatform, toPlatform string, rs *model.OcpiCommandResult) error {
	s.makeRequest(ctx, domain.LogEventPostCommandResponse, url, token, fromPlatform, toPlatform, rs, nil)
	return nil
}

func (s *mockClientImpl) PostCdr(ctx context.Context, url, token, fromPlatform, toPlatform string, cdr *model.OcpiCdr) error {
	s.makeRequest(ctx, domain.LogEventPostCdr, url, token, fromPlatform, toPlatform, cdr, nil)
	return nil
}

func (s *mockClientImpl) GetCdrsPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiCdr, error) {
	s.l().Mth("get-cdrs-page").Dbg()
	rs := &model.OcpiCdrsResponse{
		OcpiResponse: okResp,
		Data: []*model.OcpiCdr{
			{
				Id:        kit.NewRandString(),
				SessionId: kit.NewId(),
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     "TEST",
					CountryCode: "US",
				},
				LastUpdated: kit.Now(),
			},
		},
	}
	s.makeRequest(ctx, domain.LogEventGetCdrs, url, token, fromPlatform, toPlatform, rq, rs)
	return rs.Data, nil
}

func (s *mockClientImpl) GetCdr(ctx context.Context, url, token, fromPlatform, toPlatform string, cdrId string) (*model.OcpiCdr, error) {
	s.l().Mth("get-cdr").Dbg()
	rs := &model.OcpiCdr{
		Id:        kit.NewRandString(),
		SessionId: kit.NewId(),
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     "TEST",
			CountryCode: "US",
		},
		LastUpdated: kit.Now(),
	}
	s.makeRequest(ctx, domain.LogEventGetCdr, url, token, fromPlatform, toPlatform, nil, rs)
	return rs, nil
}

func (s *mockClientImpl) makeRequest(ctx context.Context, logEv, url, token, fromPlatform, toPlatform string, rq, rs any) {
	s.l().C(ctx).Mth("make").Dbg()

	// logging
	log := &domain.LogMessage{
		Event:        logEv,
		Url:          url,
		Token:        token,
		FromPlatform: fromPlatform,
		ToPlatform:   toPlatform,
	}
	if rq != nil {
		rqJs, _ := json.Marshal(rq)
		log.RequestBody = string(rqJs)
	}
	defer s.logService.Log(ctx, log)

	rsJs, _ := json.Marshal(rs)
	log.ResponseBody = string(rsJs)
	log.ResponseStatus = http.StatusOK
}
