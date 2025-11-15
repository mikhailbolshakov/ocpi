package ocpi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = time.Minute
)

type restRequest struct {
	LogEvent    string
	Url         string
	Verb        string
	QueryParams map[string]string
	Header      map[string]string
	Token       string
	Body        any
	RespModel   any
}

type ocpiRestClient interface {
	Init(ctx context.Context, config *service.CfgOcpiRemote) error
	Close(ctx context.Context) error
	GetVersions(ctx context.Context, url, token, fromPlatform, toPlatform string) ([]*model.OcpiVersion, error)
	GetVersionDetails(ctx context.Context, url, token, fromPlatform, toPlatform string) (*model.OcpiVersionDetails, error)
	PostCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiCredentials) (*model.OcpiCredentials, error)
	PutCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiCredentials) (*model.OcpiCredentials, error)
	GetCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string) (*model.OcpiCredentials, error)
	DeleteCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string) error
	GetHubClientInfo(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiClientInfo, error)
	PutClientInfo(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiClientInfo) error
	PutLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiLocation) error
	PatchLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiLocation) error
	GetLocationPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiLocation, error)
	GetLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, locId string) (*model.OcpiLocation, error)
	PutEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiEvse, party *model.OcpiPartyId, locId string) error
	PatchEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiEvse, party *model.OcpiPartyId, locId string) error
	GetEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, locId, evseId string) (*model.OcpiEvse, error)
	PutCon(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiConnector, party *model.OcpiPartyId, locId, evseId string) error
	PatchCon(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiConnector, party *model.OcpiPartyId, locId, evseId string) error
	GetCon(ctx context.Context, url, token, fromPlatform, toPlatform, locId, evseId, conId string) (*model.OcpiConnector, error)
	PutTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trf *model.OcpiTariff) error
	PatchTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trf *model.OcpiTariff) error
	GetTariffPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiTariff, error)
	GetTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trfId string) (*model.OcpiTariff, error)
	PutToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tkn *model.OcpiToken) error
	PatchToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tkn *model.OcpiToken) error
	GetTokenPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiToken, error)
	GetToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tknId string) (*model.OcpiToken, error)
	PutSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sess *model.OcpiSession) error
	PatchSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sess *model.OcpiSession) error
	GetSessionPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiSession, error)
	GetSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sessId string) (*model.OcpiSession, error)
	PostCommand(ctx context.Context, url, token, cmdType, fromPlatform, toPlatform string, cmd any) error
	PostCommandResponse(ctx context.Context, url, token, fromPlatform, toPlatform string, rs *model.OcpiCommandResult) error
	PostCdr(ctx context.Context, url, token, fromPlatform, toPlatform string, cdr *model.OcpiCdr) error
	GetCdrsPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiCdr, error)
	GetCdr(ctx context.Context, url, token, fromPlatform, toPlatform string, cdrId string) (*model.OcpiCdr, error)
}

type clientImpl struct {
	cfg        *service.CfgOcpiRemote
	logService domain.OcpiLogService
	timeout    time.Duration
}

func newOcpiRestClient(logService domain.OcpiLogService) ocpiRestClient {
	return &clientImpl{
		logService: logService,
	}
}

func (s *clientImpl) l() kit.CLogger {
	return service.L().Cmp("ocpi-rest")
}

func (s *clientImpl) Init(ctx context.Context, config *service.CfgOcpiRemote) error {
	s.l().Mth("init").Dbg()
	s.cfg = config
	if s.cfg.Timeout != nil {
		s.timeout = time.Duration(*s.cfg.Timeout) * time.Second
	} else {
		s.timeout = defaultTimeout
	}
	return nil
}

func (s *clientImpl) Close(ctx context.Context) error {
	s.l().Mth("close").Dbg()
	return nil
}

func (s *clientImpl) GetVersions(ctx context.Context, url, token, fromPlatform, toPlatform string) ([]*model.OcpiVersion, error) {
	s.l().Mth("get-version").Dbg()
	rs := &model.OcpiVersionsResponse{}
	rq := s.prepareRq(ctx, url, token, domain.LogEventGetVersions).ResponseModel(&rs).B()
	err := s.makeRequest(ctx, rq, fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetVersionDetails(ctx context.Context, url, token, fromPlatform, toPlatform string) (*model.OcpiVersionDetails, error) {
	s.l().Mth("get-version-det").Dbg()
	rs := &model.OcpiVersionDetailsResponse{}
	rq := s.prepareRq(ctx, url, token, domain.LogEventGetVersionDetails).ResponseModel(&rs).B()
	err := s.makeRequest(ctx, rq, fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) PostCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiCredentials) (*model.OcpiCredentials, error) {
	s.l().Mth("post-cred").Dbg()
	rs := &model.OcpiCredentialsResponse{}
	rq := s.prepareRq(ctx, url, token, domain.LogEventPostCredentials).
		Verb(http.MethodPost).
		Body(pl).
		ResponseModel(&rs).
		B()
	err := s.makeRequest(ctx, rq, fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) PutCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiCredentials) (*model.OcpiCredentials, error) {
	s.l().Mth("put-cred").Dbg()
	rs := &model.OcpiCredentialsResponse{}
	rq := s.prepareRq(ctx, url, token, domain.LogEventPostCredentials).
		Verb(http.MethodPut).
		Body(pl).
		ResponseModel(&rs).
		B()
	err := s.makeRequest(ctx, rq, fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string) (*model.OcpiCredentials, error) {
	s.l().Mth("get-cred").Dbg()
	rs := &model.OcpiCredentialsResponse{}
	rq := s.prepareRq(ctx, url, token, domain.LogEventGetCredentials).
		Verb(http.MethodGet).
		ResponseModel(&rs).
		B()
	err := s.makeRequest(ctx, rq, fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) DeleteCredentials(ctx context.Context, url, token, fromPlatform, toPlatform string) error {
	s.l().Mth("del-cred").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventGetCredentials).
		Verb(http.MethodGet).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetHubClientInfo(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiClientInfo, error) {
	s.l().Mth("get-hub-client").Dbg()
	rs := &model.OcpiClientInfoResponse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetHubClients).
		Verb(http.MethodGet).
		Page(rq).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) PutClientInfo(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiClientInfo) error {
	s.l().Mth("put-client-info").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutClientInfo).
		Verb(http.MethodPut).
		Body(pl).
		UrlParams(pl.CountryCode, pl.PartyId).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PutLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiLocation) error {
	s.l().Mth("put-location").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutLocation).
		Verb(http.MethodPut).
		Body(pl).
		UrlParams(pl.CountryCode, pl.PartyId, pl.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PatchLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiLocation) error {
	s.l().Mth("patch-location").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPatchLocation).
		Verb(http.MethodPatch).
		Body(pl).
		UrlParams(pl.CountryCode, pl.PartyId, pl.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetLocationPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiLocation, error) {
	s.l().Mth("get-loc-page").Dbg()
	rs := &model.OcpiLocationsResponse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetLocations).
		Verb(http.MethodGet).
		Page(rq).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetLocation(ctx context.Context, url, token, fromPlatform, toPlatform string, locId string) (*model.OcpiLocation, error) {
	s.l().Mth("get-loc").Dbg()
	rs := &model.OcpiLocation{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetLocation).
		Verb(http.MethodGet).
		UrlParams(locId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) PutEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiEvse, party *model.OcpiPartyId, locId string) error {
	s.l().Mth("put-evse").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutEvse).
		Verb(http.MethodPut).
		Body(pl).
		UrlParams(party.CountryCode, party.PartyId, locId, pl.Uid).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PatchEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiEvse, party *model.OcpiPartyId, locId string) error {
	s.l().Mth("patch-evse").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPatchEvse).
		Verb(http.MethodPatch).
		Body(pl).
		UrlParams(party.CountryCode, party.PartyId, locId, pl.Uid).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetEvse(ctx context.Context, url, token, fromPlatform, toPlatform string, locId, evseId string) (*model.OcpiEvse, error) {
	s.l().Mth("get-evse").Dbg()
	rs := &model.OcpiEvse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetEvse).
		Verb(http.MethodGet).
		UrlParams(locId, evseId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) PutCon(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiConnector, party *model.OcpiPartyId, locId, evseId string) error {
	s.l().Mth("put-con").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutCon).
		Verb(http.MethodPut).
		Body(pl).
		UrlParams(party.CountryCode, party.PartyId, locId, evseId, pl.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PatchCon(ctx context.Context, url, token, fromPlatform, toPlatform string, pl *model.OcpiConnector, party *model.OcpiPartyId, locId, evseId string) error {
	s.l().Mth("patch-con").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPatchCon).
		Verb(http.MethodPatch).
		Body(pl).
		UrlParams(party.CountryCode, party.PartyId, locId, evseId, pl.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetCon(ctx context.Context, url, token, fromPlatform, toPlatform, locId, evseId, conId string) (*model.OcpiConnector, error) {
	s.l().Mth("get-con").Dbg()
	rs := &model.OcpiConnector{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetCon).
		Verb(http.MethodGet).
		UrlParams(locId, evseId, conId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) PutTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trf *model.OcpiTariff) error {
	s.l().Mth("put-trf").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutTariff).
		Verb(http.MethodPut).
		Body(trf).
		UrlParams(trf.CountryCode, trf.PartyId, trf.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PatchTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trf *model.OcpiTariff) error {
	s.l().Mth("patch-trf").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPatchTariff).
		Verb(http.MethodPatch).
		Body(trf).
		UrlParams(trf.CountryCode, trf.PartyId, trf.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetTariffPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiTariff, error) {
	s.l().Mth("get-trf-page").Dbg()
	rs := &model.OcpiTariffsResponse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetTariffs).
		Verb(http.MethodGet).
		Page(rq).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetTariff(ctx context.Context, url, token, fromPlatform, toPlatform string, trfId string) (*model.OcpiTariff, error) {
	s.l().Mth("get-trf").Dbg()
	rs := &model.OcpiTariff{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetTariff).
		Verb(http.MethodGet).
		UrlParams(trfId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) PutToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tkn *model.OcpiToken) error {
	s.l().Mth("put-tkn").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutToken).
		Verb(http.MethodPut).
		Body(tkn).
		UrlParams(tkn.CountryCode, tkn.PartyId, tkn.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PatchToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tkn *model.OcpiToken) error {
	s.l().Mth("patch-tkn").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPatchToken).
		Verb(http.MethodPatch).
		Body(tkn).
		UrlParams(tkn.CountryCode, tkn.PartyId, tkn.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetTokenPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiToken, error) {
	s.l().Mth("get-tkn-page").Dbg()
	rs := &model.OcpiTokensResponse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetTokens).
		Verb(http.MethodGet).
		Page(rq).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetToken(ctx context.Context, url, token, fromPlatform, toPlatform string, tknId string) (*model.OcpiToken, error) {
	s.l().Mth("get-tkn").Dbg()
	rs := &model.OcpiToken{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetToken).
		Verb(http.MethodGet).
		UrlParams(tknId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) PutSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sess *model.OcpiSession) error {
	s.l().Mth("put-sess").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPutSession).
		Verb(http.MethodPut).
		Body(sess).
		UrlParams(sess.CountryCode, sess.PartyId, sess.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PatchSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sess *model.OcpiSession) error {
	s.l().Mth("patch-sess").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPatchSession).
		Verb(http.MethodPatch).
		Body(sess).
		UrlParams(sess.CountryCode, sess.PartyId, sess.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetSessionPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiSession, error) {
	s.l().Mth("get-sess-page").Dbg()
	rs := &model.OcpiSessionsResponse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetSessions).
		Verb(http.MethodGet).
		Page(rq).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetSession(ctx context.Context, url, token, fromPlatform, toPlatform string, sessId string) (*model.OcpiSession, error) {
	s.l().Mth("get-sess").Dbg()
	rs := &model.OcpiSession{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetSession).
		Verb(http.MethodGet).
		UrlParams(sessId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) PostCommand(ctx context.Context, url, token, fromPlatform, toPlatform, cmdType string, cmd any) error {
	s.l().Mth("post-cmd").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPostCommand).
		Verb(http.MethodPost).
		Body(cmd).
		UrlParams(cmdType).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PostCommandResponse(ctx context.Context, url, token, fromPlatform, toPlatform string, rs *model.OcpiCommandResult) error {
	s.l().Mth("post-cmd-rs").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPostCommandResponse).
		Verb(http.MethodPost).
		Body(rs).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) PostCdr(ctx context.Context, url, token, fromPlatform, toPlatform string, cdr *model.OcpiCdr) error {
	s.l().Mth("post-cdr").Dbg()
	rq := s.prepareRq(ctx, url, token, domain.LogEventPostCdr).
		Verb(http.MethodPost).
		Body(cdr).
		UrlParams(cdr.CountryCode, cdr.PartyId, cdr.Id).
		B()
	return s.makeRequest(ctx, rq, fromPlatform, toPlatform)
}

func (s *clientImpl) GetCdrsPage(ctx context.Context, url, token, fromPlatform, toPlatform string, rq *model.OcpiGetPageRequest) ([]*model.OcpiCdr, error) {
	s.l().Mth("get-cdrs-page").Dbg()
	rs := &model.OcpiCdrsResponse{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetCdrs).
		Verb(http.MethodGet).
		Page(rq).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Data, nil
}

func (s *clientImpl) GetCdr(ctx context.Context, url, token, fromPlatform, toPlatform string, cdrId string) (*model.OcpiCdr, error) {
	s.l().Mth("get-cdr").Dbg()
	rs := &model.OcpiCdr{}
	b := s.prepareRq(ctx, url, token, domain.LogEventGetCdr).
		Verb(http.MethodGet).
		UrlParams(cdrId).
		ResponseModel(&rs)
	err := s.makeRequest(ctx, b.B(), fromPlatform, toPlatform)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs, nil
}

func (s *clientImpl) prepareRq(ctx context.Context, url, token, ev string) *rqBuilder {
	b := newRq(url, token, ev)

	// auth token
	b.Header(model.OcpiHeaderAuth, fmt.Sprintf("Token %s", token))

	if appCtx, ok := kit.Request(ctx); ok {
		// request
		requestId := appCtx.GetRequestId()
		if requestId == "" {
			requestId = kit.NewId()
		}
		b.Header(model.OcpiHeaderRequestId, requestId)

		if len(appCtx.GetKv()) > 0 {
			// correlation
			corrId := appCtx.Kv[model.OcpiCtxCorrelationId]
			if corrId == nil {
				corrId = kit.NewId()
			}
			b.Header(model.OcpiHeaderCorrelationId, corrId.(string))

			// from
			fromCc := appCtx.Kv[model.OcpiCtxFromCountryCode]
			if fromCc != nil {
				b.Header(model.OcpiHeaderFromCountryCode, fromCc.(string))
			}
			fromParty := appCtx.Kv[model.OcpiCtxFromParty]
			if fromParty != nil {
				b.Header(model.OcpiHeaderFromPartyId, fromParty.(string))
			}

			// to
			toCc := appCtx.Kv[model.OcpiCtxToCountryCode]
			if toCc != nil {
				b.Header(model.OcpiHeaderToCountryCode, toCc.(string))
			}
			toParty := appCtx.Kv[model.OcpiCtxToParty]
			if fromParty != nil {
				b.Header(model.OcpiHeaderToPartyId, toParty.(string))
			}
		}
	}

	return b
}

func (s *clientImpl) prepareLogMsg(rq *restRequest, fromPlatform, toPlatform string) *domain.LogMessage {
	return &domain.LogMessage{
		Event:         rq.LogEvent,
		Url:           rq.Url,
		Token:         rq.Token,
		RequestId:     rq.Header[model.OcpiHeaderRequestId],
		CorrelationId: rq.Header[model.OcpiHeaderCorrelationId],
		FromPlatform:  fromPlatform,
		ToPlatform:    toPlatform,
		Headers:       rq.Header,
	}
}

func (s *clientImpl) makeRequest(ctx context.Context, rq *restRequest, fromPlatform, toPlatform string) error {
	s.l().C(ctx).Mth("make").F(kit.KV{"url": rq.Url, "verb": rq.Verb}).Dbg()

	start := kit.Now()

	// logging
	log := s.prepareLogMsg(rq, fromPlatform, toPlatform)
	defer s.logService.Log(ctx, log)

	// setup timeout
	ctxExec, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	// payload
	var rqReader io.Reader
	if rq.Body != nil {
		bodyB, _ := json.Marshal(rq.Body)
		rqReader = bytes.NewReader(bodyB)
		log.RequestBody = rq.Body
	}

	// prepare request
	req, err := http.NewRequestWithContext(ctxExec, rq.Verb, rq.Url, rqReader)
	if err != nil {
		log.Err = err
		return errors.ErrOcpiRestSendRequest(ctx, err)
	}

	// setup separate connections for each call
	req.Close = true

	req.Header.Add("Content-Type", "application/json")

	// headers
	for k, v := range rq.Header {
		req.Header.Add(k, v)
	}

	// query params
	if len(rq.QueryParams) > 0 {
		q := req.URL.Query()
		for k, v := range rq.QueryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		log.Url = req.URL.String()
	}

	// make request
	resp, err := http.DefaultClient.Do(req)
	log.DurationMs = time.Since(start).Milliseconds()
	if err != nil {
		log.Err = err
		return errors.ErrOcpiRestSendRequest(ctx, err)
	}

	// parse body
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(resp.Body)
	log.ResponseBody = string(data)
	log.ResponseStatus = resp.StatusCode
	if err != nil {
		log.Err = err
		return errors.ErrOcpiRestReadBody(ctx, err)
	}

	// check response
	respObj := &model.OcpiResponse{}
	err = json.Unmarshal(data, &respObj)
	if err != nil {
		err = errors.ErrOcpiRestParseResponse(ctx, err, resp.Status)
		log.Err = err
		return err
	}
	if respObj == nil {
		err = errors.ErrOcpiRestEmptyResponse(ctx)
		log.Err = err
		return err
	}
	log.ResponseBody = respObj
	log.OcpiStatus = respObj.StatusCode

	// check http status
	if resp.StatusCode > 300 {
		err = errors.ErrOcpiRestStatus(ctx, resp.StatusCode)
		log.Err = err
		return err
	}

	// check ocpi status
	if respObj.StatusCode != model.OcpiStatusCodeOk {
		err = errors.ErrOcpiInvalidStatus(ctx, respObj.StatusCode, respObj.StatusMessage)
		log.Err = err
		return err
	}

	// parse requested model
	if rq.RespModel != nil {
		err = json.Unmarshal(data, &rq.RespModel)
		if err != nil || rq.RespModel == nil {
			err = errors.ErrOcpiRestParseModel(ctx)
			log.Err = err
			return err
		}
	}

	return nil
}
