package storage

import (
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *logStorageImpl) toDto(msg *domain.LogMessage) *log {
	det, _ := json.Marshal(&logDetails{
		Rq:     msg.RequestBody,
		Rs:     msg.ResponseBody,
		Header: msg.Headers,
	})
	r := &log{
		Event:          msg.Event,
		Url:            pg.StringToNull(msg.Url),
		Token:          pg.StringToNull(msg.Token),
		RequestId:      pg.StringToNull(msg.RequestId),
		CorrelationId:  pg.StringToNull(msg.CorrelationId),
		FromPlatformId: pg.StringToNull(msg.FromPlatform),
		ToPlatformId:   pg.StringToNull(msg.ToPlatform),
		Status:         msg.ResponseStatus,
		OcpiStatus:     msg.OcpiStatus,
		Details:        string(det),
		DurationMs:     msg.DurationMs,
		Incoming:       msg.In,
	}
	if msg.Err != nil {
		r.Err = pg.StringToNull(msg.Err.Error())
	}
	return r
}

func (s *logStorageImpl) toDomain(log *log) *domain.LogMessage {
	det := &logDetails{}
	if log.Details != "" {
		_ = json.Unmarshal([]byte(log.Details), &det)
	}
	r := &domain.LogMessage{
		Event:          log.Event,
		Url:            pg.NullToString(log.Url),
		Token:          pg.NullToString(log.Token),
		RequestId:      pg.NullToString(log.RequestId),
		CorrelationId:  pg.NullToString(log.CorrelationId),
		FromPlatform:   pg.NullToString(log.FromPlatformId),
		ToPlatform:     pg.NullToString(log.ToPlatformId),
		RequestBody:    det.Rq,
		ResponseBody:   det.Rs,
		Headers:        det.Header,
		ResponseStatus: log.Status,
		OcpiStatus:     log.OcpiStatus,
		In:             log.Incoming,
		DurationMs:     log.DurationMs,
	}
	if log.Err != nil {
		r.Err = fmt.Errorf(*log.Err)
	}
	return r
}
