package model

type MessageDTO struct {
	ID        string `json:"id,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}
