package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/ocpi/errors"
	"io"
	"net/http"
	"net/url"
)

type Sdk struct {
	logger  *kit.Logger
	baseUrl string
	apiKey  string
}

func New(url, apiKey string, logCfg *kit.LogConfig) *Sdk {
	sdk := &Sdk{
		baseUrl: url,
		apiKey:  apiKey,
	}
	sdk.logger = kit.InitLogger(logCfg)
	return sdk
}

func (s *Sdk) l() kit.CLogger {
	return s.logFn()().Cmp("ocpi-sdk")
}

func (s *Sdk) logFn() kit.CLoggerFunc {
	return func() kit.CLogger {
		return kit.L(s.logger).Srv("ocpi-sdk")
	}
}

func (s *Sdk) Close(ctx context.Context) {}

func (s *Sdk) do(ctx context.Context, url, verb string, payload []byte) ([]byte, error) {
	l := s.l().C(ctx).Mth("do").F(kit.KV{"url": url, "verb": verb, "pl": string(payload)}).Trc()
	client := &http.Client{}
	var rqReader io.Reader
	if payload != nil {
		rqReader = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(verb, url, rqReader)
	if err != nil {
		return nil, errors.ErrSdkRequest(ctx, err)
	}

	// setup separate connections for each call
	req.Close = true

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", s.apiKey)
	if rCtx, err := kit.MustRequest(ctx); err == nil {
		req.Header.Add("RequestId", rCtx.GetRequestId())
	} else {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.ErrSdkDoRequest(ctx, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, _ := io.ReadAll(resp.Body)
	l.C(ctx).F(kit.KV{"resp": string(data), "status": resp.StatusCode}).Dbg()

	// check app error
	httpErr := &kitHttp.Error{}
	_ = json.Unmarshal(data, &httpErr)
	if httpErr != nil && httpErr.Code != "" && httpErr.Message != "" {
		return nil, kit.NewAppErrBuilder(httpErr.Code, httpErr.Message).Err()
	}

	return data, nil
}

func (s *Sdk) POST(ctx context.Context, url string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, "POST", payload)
}

func (s *Sdk) PUT(ctx context.Context, url string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, "PUT", payload)
}

func (s *Sdk) DELETE(ctx context.Context, url string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, "DELETE", payload)
}

func (s *Sdk) GET(ctx context.Context, url string) ([]byte, error) {
	return s.do(ctx, url, "GET", nil)
}

func (s *Sdk) PATCH(ctx context.Context, url string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, "PATCH", payload)
}

func (s *Sdk) toUrlParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return "?" + values.Encode()
}
