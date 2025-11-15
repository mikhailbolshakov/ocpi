package commands

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/commands/{command}/{uid}", c.SenderSetCommandResponse).POST().Auth(http.TokenB).OcpiLogging(),
		// receiver
		http.R("/ocpi/2.2.1/receiver/commands/{command}", c.ReceiverExecCommand).POST().Auth(http.TokenB).OcpiLogging(),
	}
}
