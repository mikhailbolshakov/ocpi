package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/errors"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = time.Second * 10
	apiKeyHeader   = "x-api-key"
	whEventHeader  = "x-event"
)

type webhookRestClient interface {
	Init(ctx context.Context, config *service.CfgWebHook) error
	Close(ctx context.Context) error
	CallAsync(ctx context.Context, url, apiKey, event string, payload any)
}

type clientImpl struct {
	cfg     *service.CfgWebHook
	timeout time.Duration
}

type request struct {
	Data any `json:"data"`
}

func newWhRestClient() webhookRestClient {
	return &clientImpl{}
}

func (s *clientImpl) l() kit.CLogger {
	return service.L().Cmp("wh-rest")
}

func (s *clientImpl) Init(ctx context.Context, config *service.CfgWebHook) error {
	s.l().Mth("init").Dbg()
	s.cfg = config
	if s.cfg.Timeout != nil {
		s.timeout = time.Duration(*s.cfg.Timeout) * time.Second
	} else {
		s.timeout = defaultTimeout
	}
	return nil
}

func (s *clientImpl) CallAsync(ctx context.Context, url, apiKey, event string, payload any) {
	l := s.l().Mth("call-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := s.makeRequest(ctx, url, apiKey, event, payload); err != nil {
			l.E(err).St().Err()
		}
	})
}

func (s *clientImpl) Close(ctx context.Context) error {
	s.l().Mth("close").Dbg()
	return nil
}

func (s *clientImpl) makeRequest(ctx context.Context, url, apiKey, event string, payload any) error {
	l := s.l().C(ctx).Mth("make").F(kit.KV{"url": url}).Dbg()

	// setup timeout
	ctxExec, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	// payload
	var rqReader io.Reader
	if payload != nil {
		bodyB, _ := json.Marshal(request{Data: payload})
		rqReader = bytes.NewReader(bodyB)
	}

	// prepare request
	req, err := http.NewRequestWithContext(ctxExec, http.MethodPost, url, rqReader)
	if err != nil {
		return errors.ErrWhRestSendRequest(ctx, err)
	}

	// api key
	req.Header.Add(apiKeyHeader, apiKey)
	// event
	req.Header.Add(whEventHeader, event)

	// make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.ErrWhRestSendRequest(ctx, err)
	}

	// parse body
	defer func() { _ = resp.Body.Close() }()

	// check response
	if resp.StatusCode > 300 {
		return errors.ErrWhRestStatus(ctx, resp.Status)
	}

	l.Dbg("ok")

	return nil
}
