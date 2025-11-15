package maintenance

type LogMessage struct {
	Event          string `json:"event"`                  // Event log event
	Url            string `json:"url"`                    // Url
	Token          string `json:"token"`                  // Token used for auth
	RequestId      string `json:"requestId"`              // RequestId request ID
	CorrelationId  string `json:"correlationId"`          // CorrelationId correlation ID
	FromPlatform   string `json:"fromPlatform"`           // FromPlatform platform sender
	ToPlatform     string `json:"toPlatform"`             // FromPlatform platform receiver
	RequestBody    any    `json:"requestBody,omitempty"`  // RequestBody request body
	ResponseBody   any    `json:"responseBody,omitempty"` // ResponseBody response body
	Headers        any    `json:"headers,omitempty"`      // Headers
	ResponseStatus int    `json:"responseStatus"`         // ResponseStatus http status of the response
	OcpiStatus     int    `json:"ocpiStatus"`             // OcpiStatus OCPI status of the response
	Err            string `json:"err,omitempty"`          // Err error
	In             bool   `json:"in"`                     // In if true, incoming call, otherwise outgoing
	DurationMs     int64  `json:"durationMs"`             // DurationMs request duration ms
}
