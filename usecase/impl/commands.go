package impl

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"math"
	"time"
)

const (
	cmdTimeout       = time.Minute * 10
	cmdTimeoutErrMsg = "command timed out"
)

type commandUc struct {
	ucBase
	commandService       domain.CommandService
	localPlatformService domain.LocalPlatformService
	locService           domain.LocationService
	remoteCommandRep     usecase.RemoteCommandRepository
	partyService         domain.PartyService
	webhook              backend.WebhookCallService
	converter            usecase.CommandConverter
	tokenUc              usecase.TokenUc
	tokenService         domain.TokenService
	sessionService       domain.SessionService
}

func NewCommandUc(platformService domain.PlatformService, commandService domain.CommandService, remoteCommandRep usecase.RemoteCommandRepository,
	partyService domain.PartyService, locService domain.LocationService, webhook backend.WebhookCallService,
	localPlatform domain.LocalPlatformService, tokenUc usecase.TokenUc, tokenService domain.TokenService, sessionService domain.SessionService, tokenGen domain.TokenGenerator) usecase.CommandUc {
	return &commandUc{
		ucBase:               newBase(platformService, partyService, tokenGen),
		commandService:       commandService,
		locService:           locService,
		remoteCommandRep:     remoteCommandRep,
		partyService:         partyService,
		localPlatformService: localPlatform,
		webhook:              webhook,
		tokenUc:              tokenUc,
		tokenService:         tokenService,
		sessionService:       sessionService,
		converter:            NewCommandConverter(NewTokenConverter()),
	}
}

func (t *commandUc) l() kit.CLogger {
	return ocpi.L().Cmp("cmd-uc")
}

func (t *commandUc) OnRemoteStartSession(ctx context.Context, platformId string, rq *model.OcpiStartSession) (*model.OcpiCommandResponse, error) {
	t.l().C(ctx).Mth("on-start-sess-rem").F(kit.KV{"locId": rq.LocationId}).Dbg()

	rs := &model.OcpiCommandResponse{
		Result: domain.CmdResponseTypeRejected,
	}

	// get location
	con, err := t.locService.GetConnector(ctx, rq.LocationId, rq.EvseId, rq.ConnectorId)
	if err != nil {
		return nil, err
	}
	if con == nil {
		return nil, errors.ErrCmdConNotFound(ctx)
	}
	// check if requested starting session for a CP of the local platform
	if con.PlatformId != t.localPlatformService.GetPlatformId(ctx) {
		return nil, errors.ErrLocNotBelongLocalPlatform(ctx)
	}

	// pre-populate
	cmd := t.converter.StartSessionCommandModelToDomain(rq, platformId)
	cmd.ExtId = t.getFromPartyCtx(ctx)
	cmd.Deadline = kit.Now().Add(cmdTimeout)

	// get or create token
	cmd.Details.StartSession.Token, err = t.getOrCreateRemoteToken(ctx, platformId, rq.Token)
	if err != nil {
		return nil, err
	}

	cmd.Status = domain.CmdStatusRequestAccepted

	// create command request
	cmd, err = t.commandService.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// call webhook to local platform
	err = t.webhook.OnStartSession(ctx, t.converter.CommandDomainToBackend(cmd))
	if err != nil {
		return nil, err
	}

	rs.Result = domain.CmdResponseTypeAccepted
	rs.Timeout = int(math.Round(cmd.Deadline.Sub(kit.Now()).Seconds()))

	return rs, nil
}

func (t *commandUc) OnRemoteStopSession(ctx context.Context, platformId string, rq *model.OcpiStopSession) (*model.OcpiCommandResponse, error) {
	t.l().C(ctx).Mth("on-stop-sess-rem").F(kit.KV{"sessId": rq.SessionId}).Dbg()

	if rq.SessionId == "" {
		return nil, errors.ErrSessIdEmpty(ctx)
	}

	rs := &model.OcpiCommandResponse{
		Result: domain.CmdResponseTypeRejected,
	}

	// get session
	sess, err := t.sessionService.GetSession(ctx, rq.SessionId)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, errors.ErrSessNotFound(ctx)
	}

	// check if requested session for a CP of the local platform
	if sess.PlatformId != t.localPlatformService.GetPlatformId(ctx) {
		return nil, errors.ErrCmdSessNotBelongLocalPlatform(ctx)
	}

	// pre-populate
	cmd := t.converter.StopSessionCommandModelToDomain(rq, platformId)
	cmd.ExtId = t.getFromPartyCtx(ctx)
	cmd.Deadline = kit.Now().Add(cmdTimeout)
	cmd.Status = domain.CmdStatusRequestAccepted
	cmd.AuthRef = sess.Details.AuthRef

	// create command request
	cmd, err = t.commandService.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// call webhook to local platform
	err = t.webhook.OnStopSession(ctx, t.converter.CommandDomainToBackend(cmd))
	if err != nil {
		return nil, err
	}

	rs.Result = domain.CmdResponseTypeAccepted
	rs.Timeout = int(math.Round(cmd.Deadline.Sub(kit.Now()).Seconds()))

	return rs, nil
}

func (t *commandUc) OnRemoteReserve(ctx context.Context, platformId string, rq *model.OcpiReserveNow) (*model.OcpiCommandResponse, error) {
	t.l().C(ctx).Mth("on-res-rem").F(kit.KV{"resId": rq.ReservationId}).Dbg()

	if rq.ReservationId == "" {
		return nil, errors.ErrReservationIdEmpty(ctx)
	}

	rs := &model.OcpiCommandResponse{
		Result: domain.CmdResponseTypeRejected,
	}

	// get location
	con, err := t.locService.GetConnector(ctx, rq.LocationId, rq.EvseId, rq.ConnectorId)
	if err != nil {
		return nil, err
	}
	if con == nil {
		return nil, errors.ErrCmdConNotFound(ctx)
	}

	// check if requested reservation for a CP of the local platform
	if con.PlatformId != t.localPlatformService.GetPlatformId(ctx) {
		return nil, errors.ErrLocNotBelongLocalPlatform(ctx)
	}

	// pre-populate
	cmd := t.converter.ReserveNowCommandModelToDomain(rq, platformId)
	cmd.ExtId = t.getFromPartyCtx(ctx)
	cmd.Deadline = kit.Now().Add(cmdTimeout)

	// get or create token
	cmd.Details.Reserve.Token, err = t.getOrCreateRemoteToken(ctx, platformId, rq.Token)
	if err != nil {
		return nil, err
	}

	cmd.Status = domain.CmdStatusRequestAccepted

	// create command request
	cmd, err = t.commandService.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// call webhook to local platform
	err = t.webhook.OnReserveNow(ctx, t.converter.CommandDomainToBackend(cmd))
	if err != nil {
		return nil, err
	}

	rs.Result = domain.CmdResponseTypeAccepted
	rs.Timeout = int(math.Round(cmd.Deadline.Sub(kit.Now()).Seconds()))

	return rs, nil
}

func (t *commandUc) OnRemoteCancelReservation(ctx context.Context, platformId string, rq *model.OcpiCancelReservation) (*model.OcpiCommandResponse, error) {
	t.l().C(ctx).Mth("on-res-cancel-rem").F(kit.KV{"resId": rq.ReservationId}).Dbg()

	if rq.ReservationId == "" {
		return nil, errors.ErrReservationIdEmpty(ctx)
	}

	rs := &model.OcpiCommandResponse{
		Result: domain.CmdResponseTypeRejected,
	}

	// get session
	cmdRs, err := t.commandService.SearchCommands(ctx, &domain.CommandSearchCriteria{Cmd: domain.CmdReserve, ReservationId: rq.ReservationId})
	if err != nil {
		return nil, err
	}
	if len(cmdRs.Items) == 0 {
		return nil, errors.ErrCmdCancelReservationNotFound(ctx, rq.ReservationId)
	}
	resCmd := cmdRs.Items[0]

	// check if requested session for a CP of the local platform
	if resCmd.PlatformId != platformId {
		return nil, errors.ErrCmdCancelReservationInvalidPlatform(ctx)
	}

	// pre-populate
	cmd := t.converter.CancelReservationCommandModelToDomain(rq, platformId)
	cmd.ExtId = t.getFromPartyCtx(ctx)
	cmd.Deadline = kit.Now().Add(cmdTimeout)
	cmd.Status = domain.CmdStatusRequestAccepted

	// create command request
	cmd, err = t.commandService.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// call webhook to local platform
	err = t.webhook.OnCancelReservation(ctx, t.converter.CommandDomainToBackend(cmd))
	if err != nil {
		return nil, err
	}

	rs.Result = domain.CmdResponseTypeAccepted
	rs.Timeout = int(math.Round(cmd.Deadline.Sub(kit.Now()).Seconds()))

	return rs, nil
}

func (t *commandUc) OnRemoteUnlockConnector(ctx context.Context, platformId string, rq *model.OcpiUnlockConnector) (*model.OcpiCommandResponse, error) {
	return nil, errors.ErrCmdNotSupported(ctx)
}

func (t *commandUc) OnRemoteSetResponse(ctx context.Context, platformId, uid string, rq *model.OcpiCommandResult) error {
	t.l().C(ctx).Mth("on-set-rs-rem").F(kit.KV{"uid": uid}).Dbg()

	// get command request and check status
	cmd, err := t.commandService.Get(ctx, uid)
	if err != nil {
		return err
	}
	if cmd == nil {
		return errors.ErrCmdCommandNotFound(ctx, uid)
	}

	// check if response set by the same platform
	// the only exception is timeout set by the same platform
	if cmd.PlatformId == platformId && rq.Result != domain.CmdResultTypeTimeout {
		return errors.ErrCmdCommandInvalidPlatform(ctx, uid)
	}

	// check status
	if cmd.Status != domain.CmdStatusRequestAccepted {
		return errors.ErrCmdCommandBadStatus(ctx, uid)
	}

	cmd.Details.Processing.Status = rq.Result
	if len(rq.Message) > 0 {
		cmd.Details.Processing.ErrMsg = rq.Message[0].Text
	}

	// update command
	cmd.LastUpdated = kit.Now()
	cmd, err = t.commandService.Update(ctx, cmd)
	if err != nil {
		return err
	}

	// call webhook to local platform
	return t.webhook.OnCommandResponse(ctx, t.converter.CommandDomainToBackend(cmd))
}

func (t *commandUc) RemoteCommandsDeadlineCronHandler(ctx context.Context) {
	l := t.l().C(ctx).Mth("remote-cmd-deadline").Dbg()

	localPlatformId := t.localPlatformService.GetPlatformId(ctx)

	// retrieve all commands with expired deadline
	cmdRs, err := t.commandService.SearchCommands(ctx, &domain.CommandSearchCriteria{
		ExcPlatforms: []string{localPlatformId},
		Statuses:     []string{domain.CmdStatusRequestAccepted},
		DeadlineLE:   kit.NowPtr(),
		RetrieveAll:  true,
	})
	if err != nil {
		l.E(err).St().Err("retrieve commands")
		return
	}
	l.DbgF("found: %d", len(cmdRs.Items))

	for _, cmd := range cmdRs.Items {
		err = t.OnLocalCommandSetResponse(ctx, cmd.Id, domain.CmdResultTypeTimeout, cmdTimeoutErrMsg)
		if err != nil {
			t.l().C(ctx).Mth("remote-cmd-deadline").E(err).St().Err()
		}
	}

}

func (t *commandUc) OnLocalStartSession(ctx context.Context, rq *domain.Command) error {
	l := t.l().C(ctx).Mth("on-start-sess-loc").F(kit.KV{"id": rq.Id}).Dbg()

	// local platform
	localPlatform, err := t.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get connector
	con, err := t.locService.GetConnector(ctx, rq.Details.StartSession.LocationId, rq.Details.StartSession.EvseId, rq.Details.StartSession.ConnectorId)
	if err != nil {
		return err
	}
	if con == nil {
		return errors.ErrCmdConNotFound(ctx)
	}

	// platform
	platform, err := t.getConnectedPlatform(ctx, con.PlatformId)
	if err != nil {
		return err
	}

	// check if requested starting session for a CP of the remote platform
	if !platform.Remote {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// pre-populate
	if rq.Id == "" {
		rq.Id = kit.NewId()
	}
	rq.LastUpdated = kit.Now()

	// set deadline
	rq.Deadline = kit.Now().Add(cmdTimeout)

	// set response url
	rq.Details.ResponseUrl, err = t.getLocalResponseUrl(ctx, rq.Cmd, rq.Id)
	if err != nil {
		return err
	}

	// get or create auth token
	err = t.tokenUc.OnLocalTokenChanged(ctx, rq.Details.StartSession.Token)
	if err != nil {
		return err
	}
	rq.Status = domain.CmdStatusRequestAccepted

	// create command request
	cmd, err := t.commandService.Create(ctx, rq)
	if err != nil {
		return err
	}

	// set header to route message
	ctx = t.setFromPartyCtx(ctx, cmd.ExtId)
	ctx = t.setToPartyCtx(ctx, con.ExtId)

	ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCommands, model.OcpiReceiver)
	if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Commands) {
		// push session to a remote platform
		rq := buildOcpiRepositoryErrHandlerRequest(ep, t.tokenC(platform), localPlatform, platform, l)
		t.remoteCommandRep.PostCommandAsync(ctx, rq, cmd.Cmd, t.converter.StartSessionCommandDomainToModel(cmd))
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("commands not supported")
	}

	return nil
}

func (t *commandUc) OnLocalStopSession(ctx context.Context, rq *domain.Command) error {
	l := t.l().C(ctx).Mth("on-stop-sess-loc").F(kit.KV{"id": rq.Id}).Dbg()

	if rq.Details.StopSession.SessionId == "" {
		return errors.ErrSessIdEmpty(ctx)
	}

	// local platform
	localPlatform, err := t.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get session
	sess, err := t.sessionService.GetSession(ctx, rq.Details.StopSession.SessionId)
	if err != nil {
		return err
	}
	if sess == nil {
		return errors.ErrSessNotFound(ctx)
	}

	// platform
	platform, err := t.getConnectedPlatform(ctx, sess.PlatformId)
	if err != nil {
		return err
	}

	// check if requested stopping session for a CP of the remote platform
	if !platform.Remote {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// pre-populate
	if rq.Id == "" {
		rq.Id = kit.NewId()
	}
	rq.LastUpdated = kit.Now()
	rq.Status = domain.CmdStatusRequestAccepted
	rq.AuthRef = sess.Details.AuthRef
	// set deadline
	rq.Deadline = kit.Now().Add(cmdTimeout)

	// set response url
	rq.Details.ResponseUrl, err = t.getLocalResponseUrl(ctx, rq.Cmd, rq.Id)
	if err != nil {
		return err
	}

	// create command request
	cmd, err := t.commandService.Create(ctx, rq)
	if err != nil {
		return err
	}

	// set header to route message
	t.setFromPartyCtx(ctx, cmd.ExtId)
	t.setToPartyCtx(ctx, sess.ExtId)

	ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCommands, model.OcpiReceiver)
	if ep != "" && platform.Protocol.PushSupport.Commands {
		// push session to a remote platform
		rq := buildOcpiRepositoryErrHandlerRequest(ep, t.tokenC(platform), localPlatform, platform, l)
		t.remoteCommandRep.PostCommandAsync(ctx, rq, cmd.Cmd, t.converter.StopSessionCommandDomainToModel(cmd))
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("commands not supported")
	}

	return nil
}

func (t *commandUc) OnLocalReserve(ctx context.Context, rq *domain.Command) error {
	l := t.l().C(ctx).Mth("on-res-loc").F(kit.KV{"id": rq.Id}).Dbg()

	// local platform
	localPlatform, err := t.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get location
	loc, err := t.locService.GetLocation(ctx, rq.Details.Reserve.LocationId, false)
	if err != nil {
		return err
	}
	if loc == nil {
		return errors.ErrLocationNotFound(ctx)
	}

	// platform
	platform, err := t.getConnectedPlatform(ctx, loc.PlatformId)
	if err != nil {
		return err
	}

	// check if requested reservation for a CP of the remote platform
	if !platform.Remote {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// pre-populate
	if rq.Id == "" {
		rq.Id = kit.NewId()
	}
	rq.LastUpdated = kit.Now()

	// set deadline
	rq.Deadline = kit.Now().Add(cmdTimeout)

	// set response url
	rq.Details.ResponseUrl, err = t.getLocalResponseUrl(ctx, rq.Cmd, rq.Id)
	if err != nil {
		return err
	}

	// get or create auth token
	err = t.tokenUc.OnLocalTokenChanged(ctx, rq.Details.Reserve.Token)
	if err != nil {
		return err
	}
	rq.Status = domain.CmdStatusRequestAccepted

	// create command request
	cmd, err := t.commandService.Create(ctx, rq)
	if err != nil {
		return err
	}

	// set header to route message
	ctx = t.setFromPartyCtx(ctx, cmd.ExtId)
	ctx = t.setToPartyCtx(ctx, loc.ExtId)

	ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCommands, model.OcpiReceiver)
	if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Commands) {
		// push session to a remote platform
		rq := buildOcpiRepositoryErrHandlerRequest(ep, t.tokenC(platform), localPlatform, platform, l)
		t.remoteCommandRep.PostCommandAsync(ctx, rq, cmd.Cmd, t.converter.ReserveNowCommandDomainToModel(cmd))
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("commands not supported")
	}

	return nil
}

func (t *commandUc) OnLocalCancelReservation(ctx context.Context, rq *domain.Command) error {
	l := t.l().C(ctx).Mth("on-cancel-res-loc").F(kit.KV{"id": rq.Id}).Dbg()

	if rq.Details.CancelReservation.ReservationId == "" {
		return errors.ErrReservationIdEmpty(ctx)
	}

	// local platform
	localPlatform, err := t.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get session
	cmdResRs, err := t.commandService.SearchCommands(ctx, &domain.CommandSearchCriteria{
		ReservationId: rq.Details.CancelReservation.ReservationId,
		Cmd:           domain.CmdReserve,
	})
	if err != nil {
		return err
	}
	if len(cmdResRs.Items) == 0 {
		return errors.ErrCmdCancelReservationNotFound(ctx, rq.Details.CancelReservation.ReservationId)
	}
	cmdRes := cmdResRs.Items[0]

	// platform
	resCmdPlatform, err := t.getConnectedPlatform(ctx, cmdRes.PlatformId)
	if err != nil {
		return err
	}
	// check if source reservation of the remote platform
	if resCmdPlatform.Remote {
		return errors.ErrCmdCancelReservationInvalidPlatform(ctx)
	}

	// get location
	loc, err := t.locService.GetLocation(ctx, cmdRes.Details.Reserve.LocationId, false)
	if err != nil {
		return err
	}
	if loc == nil {
		return errors.ErrLocationNotFound(ctx)
	}

	// remote platform
	platform, err := t.getConnectedPlatform(ctx, loc.PlatformId)
	if err != nil {
		return err
	}

	// check if requested cancellation reservation for a CP of the remote platform
	if !platform.Remote {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// pre-populate
	if rq.Id == "" {
		rq.Id = kit.NewId()
	}
	rq.LastUpdated = kit.Now()
	rq.Status = domain.CmdStatusRequestAccepted
	// set deadline
	rq.Deadline = kit.Now().Add(cmdTimeout)

	// set response url
	rq.Details.ResponseUrl, err = t.getLocalResponseUrl(ctx, rq.Cmd, rq.Id)
	if err != nil {
		return err
	}

	// create command request
	cmd, err := t.commandService.Create(ctx, rq)
	if err != nil {
		return err
	}

	// set header to route message
	t.setFromPartyCtx(ctx, cmd.ExtId)
	t.setToPartyCtx(ctx, loc.ExtId)

	ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCommands, model.OcpiReceiver)
	if ep != "" && platform.Protocol.PushSupport.Commands {
		// push command to a remote platform
		rq := buildOcpiRepositoryErrHandlerRequest(ep, t.tokenC(platform), localPlatform, platform, l)
		t.remoteCommandRep.PostCommandAsync(ctx, rq, cmd.Cmd, t.converter.CancelReservationCommandDomainToModel(cmd))
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("commands not supported")
	}

	return nil
}

func (t *commandUc) OnLocalUnlockConnector(ctx context.Context, rq *domain.Command) error {
	return errors.ErrCmdNotSupported(ctx)
}

func (t *commandUc) OnLocalCommandSetResponse(ctx context.Context, uid, status, errMsg string) error {
	l := t.l().C(ctx).Mth("on-set-rs-loc").F(kit.KV{"uid": uid, "status": status}).Dbg()

	// local platform
	localPlatform, err := t.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get command request and check status
	cmd, err := t.commandService.Get(ctx, uid)
	if err != nil {
		return err
	}
	if cmd == nil {
		return errors.ErrCmdCommandNotFound(ctx, uid)
	}
	if cmd.Status != domain.CmdStatusRequestAccepted {
		return errors.ErrCmdCommandBadStatus(ctx, uid)
	}
	cmd.Details.Processing.Status = status
	cmd.Details.Processing.ErrMsg = errMsg
	cmd.LastUpdated = kit.Now()

	// update command
	cmd, err = t.commandService.Update(ctx, cmd)
	if err != nil {
		return err
	}

	// get target platform
	platform, err := t.getConnectedPlatform(ctx, cmd.PlatformId)
	if err != nil {
		return err
	}

	if cmd.Details.ResponseUrl != "" {

		// build response
		rs := &model.OcpiCommandResult{
			Result: cmd.Details.Processing.Status,
		}
		if errMsg != "" {
			rs.Message = append(rs.Message, &model.OcpiDisplayText{
				Text:     errMsg,
				Language: "en",
			})
		}

		// send response
		rq := buildOcpiRepositoryErrHandlerRequestG(cmd.Details.ResponseUrl, t.tokenC(platform), localPlatform, platform, rs, l)
		t.remoteCommandRep.PostCommandResponseAsync(ctx, rq)
	}

	return nil
}

func (t *commandUc) LocalCommandsDeadlineCronHandler(ctx context.Context) {
	l := t.l().C(ctx).Mth("local-cmd-deadline").Dbg()

	localPlatformId := t.localPlatformService.GetPlatformId(ctx)

	// retrieve all commands with expired deadline
	cmdRs, err := t.commandService.SearchCommands(ctx, &domain.CommandSearchCriteria{
		IncPlatforms: []string{localPlatformId},
		Statuses:     []string{domain.CmdStatusRequestAccepted},
		DeadlineLE:   kit.NowPtr(),
		RetrieveAll:  true,
	})
	if err != nil {
		l.E(err).St().Err("retrieve commands")
		return
	}
	l.DbgF("found: %d", len(cmdRs.Items))

	for _, cmd := range cmdRs.Items {
		err = t.OnRemoteSetResponse(ctx, cmd.PlatformId, cmd.Id, &model.OcpiCommandResult{
			Result: domain.CmdResultTypeTimeout,
			Message: []*model.OcpiDisplayText{
				{
					Language: cmdTimeoutErrMsg,
					Text:     "en",
				},
			},
		})
		if err != nil {
			t.l().C(ctx).Mth("local-cmd-deadline").E(err).St().Err()
		}
	}
}

func (t *commandUc) getOrCreateRemoteToken(ctx context.Context, platformId string, tkn *model.OcpiToken) (*domain.Token, error) {
	t.l().C(ctx).Mth("get-create-tkn-rem").F(kit.KV{"tknId": tkn.Id}).Dbg()

	// put token async
	err := t.tokenUc.OnRemoteTokenPut(ctx, platformId, tkn)
	if err != nil {
		return nil, err
	}

	// get token
	stored, err := t.tokenService.GetToken(ctx, tkn.Id)
	if err != nil {
		return nil, err
	}

	// check token is valid
	if stored != nil && stored.Details.Valid != nil && !*stored.Details.Valid {
		return nil, errors.ErrTknNotValid(ctx)
	}

	return stored, nil
}

func (t *commandUc) getLocalResponseUrl(ctx context.Context, cmd, cmdId string) (domain.Endpoint, error) {
	locPl, err := t.localPlatformService.Get(ctx)
	if err != nil {
		return "", err
	}
	return domain.Endpoint(fmt.Sprintf("%s/%s/%s", locPl.Endpoints[model.ModuleIdCommands][model.OcpiSender], cmd, cmdId)), nil
}
