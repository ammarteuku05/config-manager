package dto

type AgentRegisterRequest struct {
	Name string `json:"name"`
}

type AgentRegisterResponse struct {
	AgentID             string `json:"agent_id"`
	PollURL             string `json:"poll_url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
	Code                int    `json:"code"`
	RequestID           string `json:"request_id"`
}
