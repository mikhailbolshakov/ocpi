package maintenance

import (
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (c *ctrlImpl) toLogMessage(m *domain.LogMessage, short *bool) *LogMessage {
	r := &LogMessage{
		Event:          m.Event,
		Url:            m.Url,
		RequestId:      m.RequestId,
		CorrelationId:  m.CorrelationId,
		FromPlatform:   m.FromPlatform,
		ToPlatform:     m.ToPlatform,
		ResponseStatus: m.ResponseStatus,
		OcpiStatus:     m.OcpiStatus,
		In:             m.In,
		DurationMs:     m.DurationMs,
	}
	if short == nil || !*short {
		r.RequestBody = m.RequestBody
		r.ResponseBody = m.ResponseBody
		r.Headers = m.Headers
		r.Token = m.Token[:5] + "****"
	}
	if m.Err != nil {
		r.Err = m.Err.Error()
	}
	return r
}

func (c *ctrlImpl) toLogMessages(m []*domain.LogMessage, short *bool) []*LogMessage {
	var r []*LogMessage
	for _, msg := range m {
		r = append(r, c.toLogMessage(msg, short))
	}
	return r
}
