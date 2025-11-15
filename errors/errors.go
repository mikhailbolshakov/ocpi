package errors

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/model"
	"net/http"
)

var (
	ErrPlatformStorageCreate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodePlatformStorageCreate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodePlatformStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodePlatformStorageDelete, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformStorageGetDb = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodePlatformStorageGetDb, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformEmpty, "platform empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformIdEmpty, "platform ID empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformNameEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformNameEmpty, "platform name empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformTokenAEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformTokenAEmpty, "platform token A empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformTokenEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformTokenEmpty, "platform token empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformVersionEpEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformVersionEpEmpty, "version endpoint empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformVersionEpNotValid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformVersionEpNotValid, "version endpoint invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformRoleInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformRoleInvalid, "platform role invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformRoleNotSupported = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformRoleNotSupported, "platform role not supported").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformNotFound = func(ctx context.Context, platformId string) error {
		return kit.NewAppErrBuilder(ErrCodePlatformNotFound, "platform not found").F(kit.KV{"platformId": platformId}).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformNotFoundByToken = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformNotFoundByToken, "platform not found by token").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrNoTokensSpecifiedForConnection = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeNoTokensSpecifiedForConnection, "no token specified").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrNoVersionsSpecifiedForConnection = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeNoVersionsSpecifiedForConnection, "no versions specified").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrCurrentVersionInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCurrentVersionInvalid, "current version invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyIdEmpty, "id empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyRoleInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyRoleInvalid, "role invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyRolesEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyRolesEmpty, "roles empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyStorageCreate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodePartyStorageCreate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodePartyStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformNotAvailable = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformNotAvailable, "platform not available").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformVersionsEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformVersionsEmpty, "platform versions empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrNoCompatibleVersionFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeNoCompatibleVersionFound, "no compatible version").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnsupportedVersionError}).HttpSt(http.StatusOK).Err()
	}
	ErrAuthFailed = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeAuthFailed, "no credential roles").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrAuthBackendFailed = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeAuthBackendFailed, "authentication failed").Business().C(ctx).HttpSt(http.StatusUnauthorized).Err()
	}
	ErrAppCtxPlatformIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeAppCtxPlatformIdEmpty, "platform id empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiRestSendRequest = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiRestSendRequest, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiRestReadBody = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiRestReadBody, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiRestParseResponse = func(ctx context.Context, err error, httpStatus string) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiRestParseResponse, "status: %s", httpStatus).Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiRestEmptyResponse = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiRestEmptyResponse, "empty response").C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiRestStatus = func(ctx context.Context, status int) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiRestStatus, "remote server error").C(ctx).F(kit.KV{"status": status}).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiInvalidStatus = func(ctx context.Context, status int, msg string) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiInvalidStatus, "remote server responded with ocpi error: %s", msg).C(ctx).F(kit.KV{"status": status}).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrOcpiRestParseModel = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeOcpiRestParseModel, "parsing model error").C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPlatformNotConnected = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePlatformNotConnected, "platform hasn't been connected yet").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyStorageGetDb = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodePartyStorageGetDb, "").Wrap(cause).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyStorageDelete = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodePartyStorageDelete, "").Wrap(cause).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocationNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocationNotFound, "location isn't found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownLocationError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocIdEmpty, "location id is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownLocationError}).HttpSt(http.StatusOK).Err()
	}
	ErrConNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeConNotFound, "connector not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownLocationError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocCannotMergeEvses = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocCannotMergeEvses, "evses cannot be merged").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownLocationError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseCannotMergeConnectors = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeEvseCannotMergeConnectors, "connectors cannot be merged").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownLocationError}).HttpSt(http.StatusOK).Err()
	}
	ErrExtIdInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeExtIdInvalid, "party_id or country_code is invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrMaxLenExceeded = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeMaxLenExceeded, "length of %s exceeded", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrCountryCodeInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCountryCodeInvalid, "country_code is invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyIdLen = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyIdLen, "party_id len is invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrLastUpdatedInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLastUpdatedInvalid, "last_updated is invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenClientError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocStorageTx = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeLocStorageTx, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeLocStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeLocStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeLocStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeEvseStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrConStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeConStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeEvseIdEmpty, "evse id empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrConIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeConIdEmpty, "connector id empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeEvseNotFound, "evse not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseStatusInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeEvseStatusInvalid, "evse status invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocStorageGetDb = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeLocStorageGetDb, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseStorageGetDb = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeEvseStorageGetDb, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrConStorageGetDb = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeConStorageGetDb, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseStorageTx = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeEvseStorageTx, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrConStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeConStorageDelete, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeEvseStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeEvseStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrConStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeConStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrConStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeConStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrWhRestSendRequest = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeWhRestSendRequest, "").Wrap(err).C(ctx).Err()
	}
	ErrWhRestStatus = func(ctx context.Context, status string) error {
		return kit.NewAppErrBuilder(ErrCodeWhRestStatus, "").F(kit.KV{"status": status}).C(ctx).Err()
	}
	ErrWhNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeWhNotFound, "webhook not found").Business().C(ctx).Err()
	}
	ErrWhIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeWhIdEmpty, "webhook id empty").Business().C(ctx).Err()
	}
	ErrWhApiKeyEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeWhApiKeyEmpty, "api key empty").Business().C(ctx).Err()
	}
	ErrWhEventsEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeWhEventsEmpty, "events empty").Business().C(ctx).Err()
	}
	ErrWhUrlInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeWhUrlInvalid, "url invalid").Business().C(ctx).Err()
	}
	ErrWhAlreadyExists = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeWhAlreadyExists, "webhook already exists").Business().C(ctx).Err()
	}
	ErrWhStorageCreate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeWhStorageCreate, "").Wrap(err).C(ctx).Err()
	}
	ErrWhStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeWhStorageUpdate, "").Wrap(err).C(ctx).Err()
	}
	ErrWhStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeWhStorageMerge, "").Wrap(err).C(ctx).Err()
	}
	ErrWhStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeWhStorageDelete, "").Wrap(err).C(ctx).Err()
	}
	ErrWhStorageGetDb = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeWhStorageGetDb, "").Wrap(err).C(ctx).Err()
	}
	ErrLogStorageCreate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeLogStorageCreate, "").Wrap(err).C(ctx).Err()
	}
	ErrDisplayTextInvalidAttr = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeDisplayTextInvalidAttr, "display text has invalid attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrImageInvalidAttr = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeImageInvalidAttr, "image has invalid attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrBusinessDetailsEmptyAttr = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeBusinessDetailsEmptyAttr, "business details has empty attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrBusinessDetailsInvalidAttr = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeBusinessDetailsInvalidAttr, "business details has invalid attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTimePeriodInvalid = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeTimePeriodInvalid, "date is invalid: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocDetailsHoursEmptyAttr = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeLocDetailsHoursEmptyAttr, "location hours has empty attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocDetailsHoursInvalidAttr = func(ctx context.Context, entity string, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeLocDetailsHoursInvalidAttr, "location hours has invalid attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocDetailsEmptyAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeLocDetailsEmptyAttr, "location has empty attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocDetailsInvalidAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeLocDetailsInvalidAttr, "location has invalid attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrEvseDetailsInvalidAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeEvseDetailsInvalidAttr, "evse has invalid attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrConDetailsInvalidAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeConDetailsInvalidAttr, "connector has invalid attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrConDetailsEmptyAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeConDetailsEmptyAttr, "connector has empty attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTrfIdEmpty, "tariff id is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTrfNotFound, "tariff not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeTrfInvalidAttr, "invalid tariff attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeTrfEmptyAttr, "empty tariff attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTrfStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTrfStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTrfStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTrfStorageDelete, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTknIdEmpty, "token id is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTknNotFound, "token not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownTokenError}).HttpSt(http.StatusNotFound).Err()
	}
	ErrTknEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeTknEmptyAttr, "empty token attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeTknInvalidAttr, "invalid token attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLangInvalid = func(ctx context.Context, entity string) error {
		return kit.NewAppErrBuilder(ErrCodeLangInvalid, "invalid language: %s", entity).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTknStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTknStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeTknStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeSessIdEmpty, "session id is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeSessNotFound, "session not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeSessInvalidAttr, "empty session attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeSessEmptyAttr, "invalid session attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdNotSupported = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCmdNotSupported, "command isn't supported").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdConNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCmdConNotFound, "connector not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusUnknownLocationError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocNotBelongLocalPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocNotBelongLocalPlatform, "location doesn't belong to local platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrLocNotBelongRemotePlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocNotBelongRemotePlatform, "location doesn't belong to remote platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyNotBelongRemotePlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyNotBelongRemotePlatform, "party doesn't belong to remote platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfNotBelongLocalPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTrfNotBelongLocalPlatform, "tariff doesn't belong to local platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrTrfNotBelongRemotePlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTrfNotBelongRemotePlatform, "tariff doesn't belong to remote platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyNotBelongLocalPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyNotBelongLocalPlatform, "party doesn't belong to local platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrPartyNotFoundByExt = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodePartyNotFoundByExt, "party isn't found by ext").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdCommandNotFound = func(ctx context.Context, uid string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdCommandNotFound, "command not found: %s", uid).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdCommandBadStatus = func(ctx context.Context, uid string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdCommandBadStatus, "command bad status: %s", uid).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdCommandInvalidPlatform = func(ctx context.Context, uid string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdCommandInvalidPlatform, "invalid platform: %s", uid).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCmdIdEmpty, "command id is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdAuthRefEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCmdAuthRefEmpty, "auth ref is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCmdNotFound, "command not found").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdEmptyAttr, "empty attr: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdInvalidAttr, "invalid attr: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageDelete, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCmdStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdStorageCreate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCmdStorageCreate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCmdStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCmdStorageDelete, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrTknNotValid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTknNotValid, "token isn't valid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCmdSessNotBelongLocalPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeSessNotBelongLocalPlatform, "session doesn't belong the local platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessAuthRefEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeSessAuthRefEmpty, "auth ref is empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessCmdInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeSessCmdInvalid, "command is invalid").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessCmdInvalidPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeSessCmdInvalidPlatform, "invalid platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrRouteNotValid = func(url string) error {
		return kit.NewAppErrBuilder(ErrCodeRouteNotValid, "route isn't valid: %s", url).Err()
	}
	ErrCdrSessionIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCdrSessionIdEmpty, "session id empty").C(ctx).Business().F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrSessInvalidPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCdrSessInvalidPlatform, "session invalid platform").C(ctx).Business().F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrInvalidPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCdrInvalidPlatform, "invalid platform").C(ctx).Business().F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrStorageGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCdrStorageGet, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrStorageMerge = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCdrStorageMerge, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrStorageUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCdrStorageUpdate, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrStorageDelete = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeCdrStorageDelete, "").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCdrIdEmpty, "cdr id empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrLocDetailsEmptyAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrLocDetailsEmptyAttr, "empty cdr location attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrLocDetailsInvalidAttr = func(ctx context.Context, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrLocDetailsInvalidAttr, "invalid cdr location attribute: %s", attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrEmptyAttr, "empty cdr attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrInvalidAttr, "invalid cdr attribute: %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrLocInvalidPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCdrLocInvalidPlatform, "cdr: invalid platform").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrNotFoundTariffForChargingPeriod = func(ctx context.Context, trfId string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrNotFoundTariffForChargingPeriod, "cdr: not found tariff for charging period").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrTokenInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrTokenInvalidAttr, "cdr token: invalid attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrTokenEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeCdrTokenEmptyAttr, "cdr token: empty attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrAuthInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeAuthInvalidAttr, "auth: invalid attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrAuthEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeAuthEmptyAttr, "auth: empty attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrChargingPeriodInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeChargingPeriodInvalidAttr, "charging period: invalid attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrChargingPeriodEmptyAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeChargingPeriodEmptyAttr, "charging period: empty attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrPriceInvalidAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodePriceInvalidAttr, "price: invalid attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrInvalidAmountAttr = func(ctx context.Context, entity, attr string) error {
		return kit.NewAppErrBuilder(ErrCodeInvalidAmountAttr, "amount: invalid attribute %s.%s", entity, attr).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrInvalidContext = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeInvalidContext, "invalid or empty context object").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrHeaderParamEmpty = func(ctx context.Context, a string) error {
		return kit.NewAppErrBuilder(ErrCodeInvalidContext, "obligatory header param not populated: %s", a).Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrCdrTokenEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCdrTokenEmpty, "cdr token empty").Business().C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageChargingPeriodsCreate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageChargingPeriodsCreate, "session storage: charging periods create").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageCreateTx = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageCreateTx, "session storage: charging periods create tx").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageChargingPeriodsUpdate = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageChargingPeriodsUpdate, "session storage: charging periods update").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSessStorageChargingPeriodsGet = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSessStorageChargingPeriodsGet, "session storage: charging periods get").Wrap(err).C(ctx).F(kit.KV{model.OcpiStatusField: model.OcpiStatusGenServerError}).HttpSt(http.StatusOK).Err()
	}
	ErrSdkRequest = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSdkRequest, "sdk: request").Wrap(err).C(ctx).Err()
	}
	ErrSdkDoRequest = func(ctx context.Context, err error) error {
		return kit.NewAppErrBuilder(ErrCodeSdkDoRequest, "sdk: make request").Wrap(err).C(ctx).Err()
	}
	ErrReceiverTokenInvalidBase64 = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeReceiverTokenInvalidBase64, "connect: token isn't a valid base64 string").Business().C(ctx).Err()
	}
	ErrReservationIdEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeReservationIdEmpty, "reservation: reservation_id is empty").F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Business().C(ctx).Err()
	}
	ErrCmdCancelReservationInvalidPlatform = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCmdCancelReservationInvalidPlatform, "cancel reservation: invalid platform").F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Business().C(ctx).Err()
	}
	ErrCmdCancelReservationNotFound = func(ctx context.Context, reservationId string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdCancelReservationNotFound, "cancel reservation: reservation isn't found by reservation_id (%s)", reservationId).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Business().C(ctx).Err()
	}
	ErrCmdReservationIdAlreadyExists = func(ctx context.Context, reservationId string) error {
		return kit.NewAppErrBuilder(ErrCodeCmdReservationIdAlreadyExists, "reservation: already exists reservation_id (%s)", reservationId).F(kit.KV{model.OcpiStatusField: model.OcpiStatusInvalidParamError}).HttpSt(http.StatusOK).Business().C(ctx).Err()
	}
)
