package tests

import (
	"encoding/json"
	"github.com/mikhailbolshakov/kit"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/sdk"
	"github.com/mikhailbolshakov/ocpi/transport/http"
	http2 "net/http"
	"time"
)

type emulatorSuite struct {
	kit.Suite
	cfg         *ocpi.Config
	localSdk    *sdk.Sdk
	emulatorSdk *sdk.Sdk
}

func (s *emulatorSuite) setCfg(cfg *ocpi.Config) {
	s.cfg = cfg
}

func (s *emulatorSuite) init() {
	s.Suite.Init(ocpi.LF())

	s.localSdk = sdk.New(s.cfg.Ocpi.Local.Url, s.cfg.Ocpi.Local.ApiKey, &kit.LogConfig{Level: kit.DebugLevel})
	s.emulatorSdk = sdk.New(s.cfg.Ocpi.Emulator.Url, s.cfg.Ocpi.Emulator.ApiKey, &kit.LogConfig{Level: kit.DebugLevel})

	s.connectEmulatorPlatform()
}

func (s *emulatorSuite) close() {
	s.localSdk.Close(s.Ctx)
	s.emulatorSdk.Close(s.Ctx)
}

func (s *emulatorSuite) connectEmulatorPlatform() {
	// register emulator platform in local platform
	_, err := s.localSdk.PostPlatform(s.Ctx, &backend.PlatformRequest{
		Id:           s.cfg.Ocpi.Emulator.Id,
		TokenA:       s.cfg.Ocpi.Emulator.TokenA,
		Name:         s.cfg.Ocpi.Emulator.Name,
		Role:         s.cfg.Ocpi.Emulator.Role,
		GetVersionEp: s.cfg.Ocpi.Emulator.VersionEp,
		Protocol: &backend.ProtocolDetails{
			PushSupport: backend.PushSupport{
				Credentials:   true,
				Cdrs:          true,
				Commands:      true,
				HubClientInfo: true,
				Locations:     true,
				Sessions:      true,
				Tariffs:       true,
				Tokens:        true,
			},
		},
	})
	s.NoError(err)

	// register local platform in emulator platform
	_, err = s.emulatorSdk.PostPlatform(s.Ctx, &backend.PlatformRequest{
		Id:           s.cfg.Ocpi.Local.Platform.Id,
		TokenA:       s.cfg.Ocpi.Emulator.TokenA,
		Name:         s.cfg.Ocpi.Local.Platform.Name,
		Role:         s.cfg.Ocpi.Local.Platform.Role,
		GetVersionEp: s.cfg.Ocpi.Local.Url,
		Protocol: &backend.ProtocolDetails{
			PushSupport: backend.PushSupport{
				Credentials:   true,
				Cdrs:          true,
				Commands:      true,
				HubClientInfo: true,
				Locations:     true,
				Sessions:      true,
				Tariffs:       true,
				Tokens:        true,
			},
		},
	})
	s.NoError(err)

	// connect from local to emulator
	emulatorPlatform, err := s.localSdk.ConnectPlatform(s.Ctx, s.cfg.Ocpi.Emulator.Id)
	s.NoError(err)
	s.NotEmpty(emulatorPlatform)
	s.Equal(backend.ConnectionStatusConnected, emulatorPlatform.Status)
}

func (s *emulatorSuite) openWebhookServer(port, url string, whCallback func(ev string, body map[string]interface{})) *kitHttp.Server {
	// create HTTP server
	server := kitHttp.NewHttpServer(&kitHttp.Config{
		Port:            port,
		ReadTimeoutSec:  10,
		WriteTimeoutSec: 10,
	}, ocpi.LF())
	routeBuilder := http.NewRouteBuilder(server, nil)
	routes := []*http.Route{
		http.R(url, func(w http2.ResponseWriter, r *http2.Request) {
			controller := &kitHttp.BaseController{}
			body := make(map[string]interface{})
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&body); err != nil {
				controller.RespondError(w, err)
				return
			}
			whCallback(r.Header.Get("x-event"), body)
			controller.RespondOK(w, nil)
		}).POST(),
	}
	routeBuilder.SetRoutes(routes)
	s.NoError(routeBuilder.Build())
	server.Listen()
	time.Sleep(time.Second * 2)
	return server
}
