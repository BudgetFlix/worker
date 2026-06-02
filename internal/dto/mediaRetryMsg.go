package dto

type MediaRetryMsg struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ErrorMsg string `json:"errorMsg,omitempty"`
}