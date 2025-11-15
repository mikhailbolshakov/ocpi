package sdk

import (
	"context"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
)

func (s *Sdk) DeletePartyByExt(ctx context.Context, partyId, countryCode string) error {
	l := service.L().C(ctx).Mth("del-party-ext").Dbg()
	_, err := s.DELETE(ctx, fmt.Sprintf("%s/maintenance/parties/%s/%s", s.baseUrl, partyId, countryCode), nil)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}
