package dto

type ConfigRequest struct {
	Config map[string]interface{} `json:"config"`
}

type ConfigResponse struct {
	Config    map[string]interface{} `json:"config"`
	Version   string                 `json:"version"`
	Code      int                    `json:"code"`
	RequestID string                 `json:"request_id"`
}
