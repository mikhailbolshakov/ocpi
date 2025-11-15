package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"sort"
	"strconv"
	"strings"
)

const (
	credWorkersNum = 16
	credPageSize   = 100
)

type credentialsUc struct {
	ucBase
	remotePlatformRep    usecase.RemotePlatformRepository
	platformService      domain.PlatformService
	localPlatformService domain.LocalPlatformService
	partyService         domain.PartyService
	converter            usecase.CredentialsConverter
	webhook              backend.WebhookCallService
	hubUc                usecase.HubUc
}

func NewCredentialsUc(platformService domain.PlatformService, localPlatformService domain.LocalPlatformService, tokenGen domain.TokenGenerator, remotePlatformRep usecase.RemotePlatformRepository, partyService domain.PartyService,
	webhook backend.WebhookCallService, hubUc usecase.HubUc) usecase.CredentialsUc {
	return &credentialsUc{
		ucBase:               newBase(platformService, partyService, tokenGen),
		platformService:      platformService,
		localPlatformService: localPlatformService,
		remotePlatformRep:    remotePlatformRep,
		partyService:         partyService,
		converter:            NewCredentialsConverter(),
		hubUc:                hubUc,
		webhook:              webhook,
	}
}

func (c *credentialsUc) l() kit.CLogger {
	return ocpi.L().Cmp("cred-uc")
}

func (c *credentialsUc) EstablishConnection(ctx context.Context, receiverPlatformId string) (*domain.Platform, error) {
	c.l().C(ctx).Mth("establish").F(kit.KV{"receiverPlatformId": receiverPlatformId}).Dbg()

	// check incoming receiver platform id
	if receiverPlatformId == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}

	//get receiver platform
	receiverPlatform, err := c.platformService.Get(ctx, receiverPlatformId)
	if err != nil {
		return nil, err
	}
	if receiverPlatform == nil {
		return nil, errors.ErrPlatformNotFound(ctx, receiverPlatformId)
	}

	// reconnect platform
	return c.reconnect(ctx, receiverPlatform)
}

func (c *credentialsUc) UpdateConnection(ctx context.Context, receiverPlatformId string) (*domain.Platform, error) {
	c.l().C(ctx).Mth("update").F(kit.KV{"receiverPlatformId": receiverPlatformId}).Dbg()

	// check incoming receiver platform id
	if receiverPlatformId == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}

	//get receiver platform
	receiverPlatform, err := c.platformService.Get(ctx, receiverPlatformId)
	if err != nil {
		return nil, err
	}
	if receiverPlatform == nil {
		return nil, errors.ErrPlatformNotFound(ctx, receiverPlatformId)
	}
	if receiverPlatform.TokenB == "" {
		return nil, errors.ErrPlatformNotConnected(ctx)
	}

	// reconnect platform
	return c.reconnect(ctx, receiverPlatform)
}

func (c *credentialsUc) AcceptConnection(ctx context.Context, senderPlatformId string, rq *model.OcpiCredentials) (*model.OcpiCredentials, error) {

	// check incoming sender token
	if senderPlatformId == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}

	// get local platform & parties (receiver)
	localPlatform, localParties, err := c.mustGetLocalPlatformWithParties(ctx)
	if err != nil {
		return nil, err
	}

	//get sender platform
	senderPlatform, err := c.platformService.Get(ctx, senderPlatformId)
	if err != nil {
		return nil, err
	}
	if senderPlatform == nil {
		return nil, errors.ErrPlatformNotFoundByToken(ctx)
	}
	if senderPlatform.Status == domain.ConnectionStatusSuspended {
		return nil, errors.ErrPlatformNotAvailable(ctx)
	}

	// if base64 used we have to populate auth header with encoded tokens
	// at the same time incoming payload brings not encoded tokens (OCPI 7.3.1)
	senderToken := domain.PlatformToken(rq.Token)
	if senderPlatform.TokenBase64 != nil && *senderPlatform.TokenBase64 {
		// encode sender token to use in Auth headers
		senderToken = c.tokenGen.Base64Encode(senderToken)
	}

	// request versions from sender
	senderVersions, err := c.remotePlatformRep.GetVersions(ctx, buildOcpiRepositoryRequest(domain.Endpoint(rq.Url), senderToken, localPlatform, senderPlatform))
	if err != nil {
		return nil, err
	}
	if len(senderVersions) == 0 {
		return nil, errors.ErrPlatformVersionsEmpty(ctx)
	}

	// find latest mutual version
	version := c.findProperVersion(senderVersions, localPlatform.VersionInfo.Available)
	if version == "" {
		return nil, errors.ErrNoCompatibleVersionFound(ctx)
	}

	// get sender endpoints
	senderEndpoints, err := c.remotePlatformRep.GetVersionDetails(ctx, buildOcpiRepositoryRequest(senderVersions[version], senderToken, localPlatform, senderPlatform))
	if err != nil {
		return nil, err
	}

	// generate receiver token
	localToken, err := c.tokenGen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	// merge credential roles from sender
	parties := c.converter.PartiesModelToDomain(senderPlatform.Id, senderPlatform.Status, rq.Roles...)
	err = c.partyService.MergeMany(ctx, parties...)
	if err != nil {
		return nil, err
	}

	// create connection per each role received
	senderPlatform.Endpoints = senderEndpoints
	senderPlatform.TokenB = localToken
	senderPlatform.TokenC = domain.PlatformToken(rq.Token)
	senderPlatform.VersionInfo.Available = senderVersions
	senderPlatform.VersionInfo.Current = version
	_, err = c.platformService.Merge(ctx, senderPlatform)
	if err != nil {
		return nil, err
	}

	// set platform status
	_, err = c.platformService.SetStatus(ctx, senderPlatform.Id, domain.ConnectionStatusConnected)
	if err != nil {
		return nil, err
	}

	rs := &model.OcpiCredentials{
		Token: string(localToken),
		Url:   string(localPlatform.VersionInfo.VersionEp),
		Roles: c.converter.PartiesDomainToModel(localParties...),
	}

	// call webhook
	err = c.webhook.OnPartiesChanged(ctx, c.converter.PartiesDomainToBackend(parties)...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (c *credentialsUc) OnRemotePartyPull(ctx context.Context) error {
	l := c.l().C(ctx).Mth("on-party-pull-rem").Dbg()

	// get local platform
	localPlatform, err := c.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get platforms to pull from
	platforms, err := c.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	ctx = kit.Copy(ctx)

	// execute workers
	ch := make(chan channelData, credWorkersNum)
	for i := 0; i < credWorkersNum; i++ {
		goroutine.New().WithLogger(l).Go(ctx, func() {
			for v := range ch {
				cred := v.data.(*model.OcpiCredentials)

				if cred == nil || len(cred.Roles) == 0 {
					continue
				}

				// merge credential roles from sender
				parties := c.converter.PartiesModelToDomain(v.platformId, domain.ConnectionStatusConnected, cred.Roles...)
				err = c.partyService.MergeMany(ctx, parties...)
				if err != nil {
					c.l().C(ctx).Mth("merge-parties").F(kit.KV{"platformId": v.platformId}).E(err).St().Err()
					continue
				}

				// call webhook
				err = c.webhook.OnPartiesChanged(ctx, c.converter.PartiesDomainToBackend(parties)...)
				if err != nil {
					c.l().C(ctx).Mth("webhook-call").F(kit.KV{"platformId": v.platformId}).E(err).St().Err()
					continue
				}

			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := c.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCredentials, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			credentials, err := c.remotePlatformRep.GetCredentials(ctx, buildOcpiRepositoryRequest(ep, c.tokenC(platform), localPlatform, platform))
			if err != nil {
				c.l().C(ctx).Mth("request-credentials").F(kit.KV{"platformId": platform.Id}).E(err).St().Err()
				return err
			}
			ch <- channelData{platformId: platform.Id, data: credentials}
			return nil
		})
	}

	// close channels when done
	goroutine.New().WithLogger(l).Go(ctx, func() {
		defer close(ch)
		_ = eg.Wait()
	})

	return nil
}

func (c *credentialsUc) OnRemoteGetCredentials(ctx context.Context, platformId string) (*model.OcpiCredentials, error) {
	c.l().C(ctx).Mth("on-get-rem").F(kit.KV{"platformId": platformId}).Dbg()

	// get connected platform
	platform, err := c.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return nil, err
	}

	// local platform
	localPlatform, err := c.localPlatformService.Get(ctx)
	if err != nil {
		return nil, err
	}

	// get local parties
	localParties, err := c.getLocalParties(ctx, localPlatform)
	if err != nil {
		return nil, err
	}

	// build response
	rs := &model.OcpiCredentials{
		Token: string(platform.TokenC),
		Url:   string(localPlatform.VersionInfo.VersionEp),
		Roles: c.converter.PartiesDomainToModel(localParties...),
	}

	return rs, nil
}

func (c *credentialsUc) OnRemoteDeleteCredentials(ctx context.Context, platformId string) error {
	c.l().C(ctx).Mth("on-delete-rem").F(kit.KV{"platformId": platformId}).Dbg()

	// get connected platform
	platform, err := c.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// set platform status
	_, err = c.platformService.SetStatus(ctx, platform.Id, domain.ConnectionStatusSuspended)
	return err
}

func (c *credentialsUc) OnLocalPartyChanged(ctx context.Context, party *domain.Party) error {
	l := c.l().C(ctx).Mth("on-party-changed-loc").F(kit.KV{"partyId": party.Id}).Dbg()

	// get local platform
	localPlatform, err := c.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get stored
	stored, err := c.partyService.Get(ctx, party.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != localPlatform.Id {
		return errors.ErrPartyNotBelongLocalPlatform(ctx)
	}

	party.PlatformId = localPlatform.Id

	// merge location to local platform
	party, err = c.partyService.Merge(ctx, party)
	if err != nil {
		return err
	}

	// no changes applied
	if party == nil {
		l.Warn("no changes applied")
		return nil
	}

	// if local platform is a hub
	if localPlatform.Role == domain.RoleHUB {
		err = c.hubUc.OnLocalClientInfoChanged(ctx, party)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *credentialsUc) reconnect(ctx context.Context, receiverPlatform *domain.Platform) (*domain.Platform, error) {

	// get local platform & parties (sender)
	localPlatform, localParties, err := c.mustGetLocalPlatformWithParties(ctx)
	if err != nil {
		return nil, err
	}

	// if base64 specified for the receiver platform, encode local token to pass to the receiver
	encodedTokenA := receiverPlatform.TokenA
	if receiverPlatform.TokenBase64 != nil && *receiverPlatform.TokenBase64 {
		encodedTokenA = c.tokenGen.Base64Encode(receiverPlatform.TokenA)
	}

	// request versions from receiver
	receiverVersions, err := c.remotePlatformRep.GetVersions(ctx, buildOcpiRepositoryRequest(receiverPlatform.VersionInfo.VersionEp, encodedTokenA, localPlatform, receiverPlatform))
	if err != nil {
		return nil, err
	}
	if len(receiverVersions) == 0 {
		return nil, errors.ErrPlatformVersionsEmpty(ctx)
	}

	// find latest mutual version
	version := c.findProperVersion(receiverVersions, localPlatform.VersionInfo.Available)
	if version == "" {
		return nil, errors.ErrNoCompatibleVersionFound(ctx)
	}

	// get receiver endpoints
	receiverEndpoints, err := c.remotePlatformRep.GetVersionDetails(ctx, buildOcpiRepositoryRequest(receiverVersions[version], encodedTokenA, localPlatform, receiverPlatform))
	if err != nil {
		return nil, err
	}

	// generate sender token
	localToken, err := c.tokenGen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	// update receiver platform
	receiverPlatform.TokenB = localToken
	receiverPlatform.VersionInfo.Available = receiverVersions
	receiverPlatform.VersionInfo.Current = version
	receiverPlatform.Endpoints = receiverEndpoints
	receiverPlatform, err = c.platformService.Merge(ctx, receiverPlatform)
	if err != nil {
		return nil, err
	}

	// send credentials request to the receiver platform
	credRq := &model.OcpiCredentials{
		Token: string(localToken),
		Url:   string(localPlatform.VersionInfo.VersionEp),
		Roles: c.converter.PartiesDomainToModel(localParties...),
	}

	// 7.2 credential role is symmetric, so it should have either RECEIVER or SENDER role
	ep := c.roleEndpoint(ctx, receiverEndpoints, model.ModuleIdCredentials, model.OcpiSender)
	if ep == "" {
		ep = c.roleEndpoint(ctx, receiverEndpoints, model.ModuleIdCredentials, model.OcpiReceiver)
		if ep == "" {
			return nil, errors.ErrPlatformRoleNotSupported(ctx)
		}
	}

	// post credentials
	var credRs *model.OcpiCredentials
	if receiverPlatform.TokenC == "" {
		credRs, err = c.remotePlatformRep.PostCredentials(ctx, buildOcpiRepositoryRequestG(ep, encodedTokenA, localPlatform, receiverPlatform, credRq))
	} else {
		credRs, err = c.remotePlatformRep.PutCredentials(ctx, buildOcpiRepositoryRequestG(ep, c.tokenC(receiverPlatform), localPlatform, receiverPlatform, credRq))
	}
	if err != nil {
		return nil, err
	}

	// merge credential roles from sender
	parties := c.converter.PartiesModelToDomain(receiverPlatform.Id, receiverPlatform.Status, credRs.Roles...)
	err = c.partyService.MergeMany(ctx, parties...)
	if err != nil {
		return nil, err
	}

	// update platform
	receiverPlatform.TokenC = domain.PlatformToken(credRs.Token)
	receiverPlatform, err = c.platformService.Merge(ctx, receiverPlatform)
	if err != nil {
		return nil, err
	}

	// set platform status
	receiverPlatform, err = c.platformService.SetStatus(ctx, receiverPlatform.Id, domain.ConnectionStatusConnected)
	if err != nil {
		return nil, err
	}

	// call webhook
	err = c.webhook.OnPartiesChanged(ctx, c.converter.PartiesDomainToBackend(parties)...)
	if err != nil {
		return nil, err
	}

	return receiverPlatform, nil
}

func (c *credentialsUc) verToInt(v string) int {
	a := strings.Replace(v, ".", "", -1)
	if len(a) == 2 {
		a = a + "0"
	} else if len(a) != 3 {
		return 0
	}
	vInt, _ := strconv.ParseInt(a, 10, 8)
	return int(vInt)
}

func (c *credentialsUc) findProperVersion(remoteVersions, localVersions domain.Versions) string {
	var remVerInt []int
	var locVer = map[int]string{}
	for v := range remoteVersions {
		remVerInt = append(remVerInt, c.verToInt(v))
	}
	for v := range localVersions {
		locVer[c.verToInt(v)] = v
	}
	sort.Sort(sort.Reverse(sort.IntSlice(remVerInt)))
	for _, r := range remVerInt {
		if h, ok := locVer[r]; ok {
			return h
		}
	}
	return ""
}

func (c *credentialsUc) roleEndpoint(ctx context.Context, endpoints domain.ModuleEndpoints, module, role string) domain.Endpoint {
	// check if push location is supported
	locModule, ok := endpoints[module]
	if !ok {
		// module not supported
		return ""
	}
	ep, ok := locModule[role]
	if !ok {
		// receiver interface isn't supported
		return ""
	}
	return ep
}

func (c *credentialsUc) mustGetLocalPlatformWithParties(ctx context.Context) (*domain.Platform, []*domain.Party, error) {
	// get local platform
	localPlatform, err := c.localPlatformService.Get(ctx)
	if err != nil {
		return nil, nil, err
	}

	// check local platform
	if localPlatform == nil {
		return nil, nil, errors.ErrPlatformNotFound(ctx, "")
	}
	if localPlatform.Status == domain.ConnectionStatusSuspended {
		return nil, nil, errors.ErrPlatformNotAvailable(ctx)
	}

	// get parties
	localParties, err := c.getLocalParties(ctx, localPlatform)
	if err != nil {
		return nil, nil, err
	}
	return localPlatform, localParties, nil
}

func (c *credentialsUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
	platforms, err := c.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		Remote:   kit.BoolPtr(true),                          // remote platforms
		ExcRoles: []string{domain.RoleHUB},                   // do not pull parties from HUBs
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}
