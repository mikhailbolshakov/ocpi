package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type cmdService struct {
	base
	storage      domain.CommandStorage
	tokenService domain.TokenService
}

func NewCmdService(tokenService domain.TokenService, storage domain.CommandStorage) domain.CommandService {
	return &cmdService{
		storage:      storage,
		tokenService: tokenService,
	}
}

var (
	cmdStatusMap = map[string]struct{}{
		domain.CmdStatusRequestAccepted:        {},
		domain.CmdStatusRequestRejected:        {},
		domain.CmdStatusRequestProcessedOk:     {},
		domain.CmdStatusRequestProcessedFailed: {},
		domain.CmdStatusRequestExpired:         {},
	}

	cmdTypeMap = map[string]struct{}{
		domain.CmdStartSession:      {},
		domain.CmdStopSession:       {},
		domain.CmdReserve:           {},
		domain.CmdCancelReservation: {},
		domain.CmdUnlockConnector:   {},
	}

	cmdProcStatusMap = map[string]struct{}{
		domain.CmdResultTypeAccepted:             {},
		domain.CmdResultTypeCancelledReservation: {},
		domain.CmdResultTypeEvseOccupied:         {},
		domain.CmdResultTypeEvseInoperative:      {},
		domain.CmdResultTypeFailed:               {},
		domain.CmdResultTypeNotSupported:         {},
		domain.CmdResultTypeRejected:             {},
		domain.CmdResultTypeTimeout:              {},
		domain.CmdResultTypeUnknownReservation:   {},
	}

	cmdResToStatus = map[string]string{
		domain.CmdResultTypeAccepted:             domain.CmdStatusRequestProcessedOk,
		domain.CmdResultTypeCancelledReservation: domain.CmdStatusRequestProcessedFailed,
		domain.CmdResultTypeEvseOccupied:         domain.CmdStatusRequestProcessedFailed,
		domain.CmdResultTypeEvseInoperative:      domain.CmdStatusRequestProcessedFailed,
		domain.CmdResultTypeFailed:               domain.CmdStatusRequestProcessedFailed,
		domain.CmdResultTypeNotSupported:         domain.CmdStatusRequestProcessedFailed,
		domain.CmdResultTypeRejected:             domain.CmdStatusRequestRejected,
		domain.CmdResultTypeTimeout:              domain.CmdStatusRequestProcessedFailed,
		domain.CmdResultTypeUnknownReservation:   domain.CmdStatusRequestProcessedFailed,
	}
)

func (s *cmdService) l() kit.CLogger {
	return ocpi.L().Cmp("cmd-svc")
}

func (s *cmdService) Create(ctx context.Context, cmd *domain.Command) (*domain.Command, error) {
	s.l().C(ctx).Mth("create").F(kit.KV{"cmdId": cmd.Id}).Dbg()

	if cmd.Id == "" {
		cmd.Id = kit.NewId()
	}

	// validate
	err := s.validateAndPopulateCreate(ctx, cmd)
	if err != nil {
		return nil, err
	}

	err = s.storage.CreateCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s *cmdService) Update(ctx context.Context, cmd *domain.Command) (*domain.Command, error) {
	l := s.l().C(ctx).Mth("update").F(kit.KV{"cmdId": cmd.Id}).Dbg()

	if cmd.Id == "" {
		return nil, errors.ErrCmdIdEmpty(ctx)
	}

	// get stored
	stored, err := s.storage.GetCommand(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrCmdNotFound(ctx)
	}

	// check last_updated
	if cmd.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateUpdate(ctx, cmd, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.UpdateCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s *cmdService) Get(ctx context.Context, id string) (*domain.Command, error) {
	s.l().C(ctx).Mth("get").Dbg()
	if id == "" {
		return nil, errors.ErrCmdIdEmpty(ctx)
	}
	return s.storage.GetCommand(ctx, id)
}

func (s *cmdService) GetByAuthRef(ctx context.Context, authRef string) (*domain.Command, error) {
	s.l().C(ctx).Mth("get-by-auth-ref").Dbg()
	if authRef == "" {
		return nil, errors.ErrCmdAuthRefEmpty(ctx)
	}
	return s.storage.GetCommandByAuthRef(ctx, authRef)
}

func (s *cmdService) DeleteCommandsByExt(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeleteCommandsByExt(ctx, extId)
}

func (s *cmdService) SearchCommands(ctx context.Context, cr *domain.CommandSearchCriteria) (*domain.CommandSearchResponse, error) {
	s.l().C(ctx).Mth("search").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchCommands(ctx, cr)
}

func (s *cmdService) validateAndPopulateCreate(ctx context.Context, cmd *domain.Command) error {
	return s.validate(ctx, cmd)
}

func (s *cmdService) validateAndPopulateUpdate(ctx context.Context, cmd, stored *domain.Command) error {
	if cmd.Details.Processing.Status != "" {
		cmd.Status = cmdResToStatus[cmd.Details.Processing.Status]
	}
	return s.validate(ctx, cmd)
}

func (s *cmdService) validate(ctx context.Context, cmd *domain.Command) error {

	err := s.validateOcpiItem(ctx, &cmd.OcpiItem)
	if err != nil {
		return err
	}
	if cmd.Id == "" {
		return errors.ErrCmdIdEmpty(ctx)
	}
	if err := s.validateId(ctx, cmd.Id, "id"); err != nil {
		return err
	}

	// status
	if cmd.Status == "" {
		return errors.ErrCmdEmptyAttr(ctx, "cmd", "status")
	}
	if _, ok := cmdStatusMap[cmd.Status]; !ok {
		return errors.ErrCmdInvalidAttr(ctx, "cmd", "status")
	}

	// command type
	if cmd.Cmd == "" {
		return errors.ErrCmdEmptyAttr(ctx, "cmd", "cmd_type")
	}
	if _, ok := cmdTypeMap[cmd.Cmd]; !ok {
		return errors.ErrCmdInvalidAttr(ctx, "cmd", "cmd_type")
	}

	// response url
	if cmd.Details.ResponseUrl == "" || !kit.IsUrlValid(string(cmd.Details.ResponseUrl)) {
		return errors.ErrCmdInvalidAttr(ctx, "cmd", "response_url")
	}

	// validate processing
	if cmd.Details.Processing.Status != "" {
		if _, ok := cmdProcStatusMap[cmd.Details.Processing.Status]; !ok {
			return errors.ErrCmdInvalidAttr(ctx, "cmd", "processing_status")
		}
	}

	// details
	switch cmd.Cmd {
	case domain.CmdStartSession:
		err = s.validateStartSession(ctx, cmd)
	case domain.CmdStopSession:
		err = s.validateStopSession(ctx, cmd)
	case domain.CmdReserve:
		err = s.validateReservation(ctx, cmd)
	case domain.CmdCancelReservation:
		err = s.validateCancelReservation(ctx, cmd)
	case domain.CmdUnlockConnector:
		err = s.validateUnlockConnector(ctx, cmd)
	}
	return err
}

func (s *cmdService) validateStartSession(ctx context.Context, cmd *domain.Command) error {
	if cmd.Details.StartSession == nil {
		return errors.ErrCmdEmptyAttr(ctx, "cmd", "start_session")
	}
	if cmd.Details.StartSession.LocationId == "" {
		return errors.ErrCmdEmptyAttr(ctx, "start_session", "location_id")
	}
	if cmd.Details.StartSession.EvseId == "" {
		return errors.ErrCmdEmptyAttr(ctx, "start_session", "evse_id")
	}
	if cmd.Details.StartSession.ConnectorId == "" {
		return errors.ErrCmdEmptyAttr(ctx, "start_session", "connector_id")
	}
	if err := s.validateId(ctx, cmd.Details.StartSession.EvseId, "evse_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.StartSession.LocationId, "location_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.StartSession.ConnectorId, "connector_id"); err != nil {
		return err
	}
	if cmd.Details.StartSession.Token == nil {
		return errors.ErrCmdEmptyAttr(ctx, "start_session", "token")
	}
	if cmd.Details.StartSession.KwhLimit != nil && *cmd.Details.StartSession.KwhLimit < 0 {
		return errors.ErrCmdInvalidAttr(ctx, "start_session", "kwh_limit")
	}
	return s.tokenService.ValidateToken(ctx, cmd.Details.StartSession.Token)
}

func (s *cmdService) validateStopSession(ctx context.Context, cmd *domain.Command) error {
	if cmd.Details.StopSession == nil {
		return errors.ErrCmdEmptyAttr(ctx, "cmd", "stop_session")
	}
	if cmd.Details.StopSession.SessionId == "" {
		return errors.ErrCmdEmptyAttr(ctx, "stop_session", "session_id")
	}
	if err := s.validateId(ctx, cmd.Details.StopSession.SessionId, "session_id"); err != nil {
		return err
	}
	return nil
}

func (s *cmdService) validateReservation(ctx context.Context, cmd *domain.Command) error {
	if err := s.validateId(ctx, cmd.Details.Reserve.ReservationId, "reservation_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.Reserve.EvseId, "evse_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.Reserve.LocationId, "location_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.Reserve.ConnectorId, "connector_id"); err != nil {
		return err
	}

	// check if reservation already exists
	sRs, err := s.storage.SearchCommands(ctx, &domain.CommandSearchCriteria{ReservationId: cmd.Details.Reserve.ReservationId})
	if err != nil {
		return err
	}

	// if exists another command with the same reservation_id
	// ocpi protocol allows another command with the same reservation_id, but we don't want to support this
	if len(sRs.Items) > 0 && kit.First(sRs.Items, func(c *domain.Command) bool { return c.Id != cmd.Id }) != nil {
		return errors.ErrCmdReservationIdAlreadyExists(ctx, cmd.Details.Reserve.ReservationId)
	}

	return nil
}

func (s *cmdService) validateCancelReservation(ctx context.Context, cmd *domain.Command) error {
	if err := s.validateId(ctx, cmd.Details.CancelReservation.ReservationId, "reservation_id"); err != nil {
		return err
	}
	return nil
}

func (s *cmdService) validateUnlockConnector(ctx context.Context, cmd *domain.Command) error {
	if err := s.validateId(ctx, cmd.Details.UnlockConnector.EvseId, "evse_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.UnlockConnector.LocationId, "location_id"); err != nil {
		return err
	}
	if err := s.validateId(ctx, cmd.Details.UnlockConnector.ConnectorId, "connector_id"); err != nil {
		return err
	}
	return nil
}
