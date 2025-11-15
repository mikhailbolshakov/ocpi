package commands

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/backend/commands/sessions/start", c.StartSession).POST().ApiKey(),
		http.R("/backend/commands/sessions/stop", c.StopSession).POST().ApiKey(),
		http.R("/backend/commands/reservations", c.Reservation).POST().ApiKey(),
		http.R("/backend/commands/reservations", c.CancelReservation).DELETE().ApiKey(),
		http.R("/backend/commands/connectors/{connId}/unlock", c.UnLockConnector).DELETE().ApiKey(),
		http.R("/backend/commands/{cmdId}", c.GetCommand).GET().ApiKey(),
		http.R("/backend/commands/{cmdId}/response", c.PutCommandResponse).POST().ApiKey(),
		http.R("/backend/commands/search/query", c.SearchCommands).GET().ApiKey(),
	}
}
