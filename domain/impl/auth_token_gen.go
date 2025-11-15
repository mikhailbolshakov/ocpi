package impl

import (
	"context"
	"encoding/base64"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
)

type tokenGen struct{}

func NewTokenGenerator() domain.TokenGenerator {
	return &tokenGen{}
}

func (t *tokenGen) l() kit.CLogger {
	return ocpi.L().Cmp("token-gen")
}

func (t *tokenGen) Base64Encode(tkn domain.PlatformToken) domain.PlatformToken {
	return domain.PlatformToken(base64.StdEncoding.EncodeToString([]byte(tkn)))
}

func (t *tokenGen) TryBase64Decode(token domain.PlatformToken) (domain.PlatformToken, bool) {
	decoded, err := base64.StdEncoding.DecodeString(string(token))
	if err != nil || len(decoded) == 0 {
		// not base64 string
		return "", false
	}
	return domain.PlatformToken(decoded), true
}

func (t *tokenGen) Generate(ctx context.Context) (domain.PlatformToken, error) {
	token := kit.NewRandString()
	return domain.PlatformToken(token), nil
}
